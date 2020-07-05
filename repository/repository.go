package repository

import (
	"context"
	"go-demo-mongodb/canonical"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Repository interface {
	Add(*canonical.Player) error
	Update(*canonical.Player) error
	Get(id string) (canonical.Player, error)
	GetAll(offset int, limit int) ([]canonical.Player, error)
	Delete(id string) error
	Upload(id string, bin []byte) error
	Download(id string) ([]byte, error)
}

type repository struct {
	client *mongo.Client
}

type HexId struct {
	ID primitive.ObjectID `bson:"_id"`
}

type Image struct {
	Data []byte `bson:"data"`
}

const database = "test"

var repo Repository

func NewRepository() Repository {
	if repo == nil {
		repo = &repository{
			client: connect(),
		}
	}
	return repo
}

func connect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		log.Panic(err)
	}
	//defer disconnect(err, client, ctx)
	return client
}

func disconnect(err error, client *mongo.Client, ctx context.Context) {
	if err = client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (r repository) Add(player *canonical.Player) error {
	result, err := r.client.Database(database).Collection("players").
		InsertOne(context.Background(), player)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		player.Id = oid.Hex()
	}

	return nil
}

func (r repository) Update(player *canonical.Player) error {
	objId, _ := primitive.ObjectIDFromHex(player.Id)
	_, err := r.client.Database(database).Collection("players").
		UpdateOne(context.Background(), bson.M{"_id": objId}, bson.D{{Key: "$set", Value: player}})
	return err
}

func (r repository) Get(id string) (canonical.Player, error) {
	player := canonical.Player{}
	objID, _ := primitive.ObjectIDFromHex(id)
	err := r.client.Database(database).Collection("players").
		FindOne(context.Background(), bson.M{"_id": objID}).Decode(&player)

	player.Id = objID.Hex()

	if err == mongo.ErrNoDocuments {
		return canonical.Player{}, nil
	}

	return player, err
}

func (r repository) GetAll(offset int, limit int) ([]canonical.Player, error) {
	players := []canonical.Player{}
	ctx := context.Background()
	findOptions := options.Find()
	findOptions.SetLimit(5)
	findOptions.SetSkip(int64(offset))

	cur, err := r.client.Database(database).Collection("players").
		Find(ctx, bson.D{}, findOptions)

	if err != nil {
		return players, err
	}

	defer cur.Close(ctx)
	for cur.Next(context.TODO()) {
		var p canonical.Player
		hex := HexId{}

		cur.Decode(&hex)
		err := cur.Decode(&p)
		if err != nil {
			log.Fatal(err)
		}

		p.Id = hex.ID.Hex()
		players = append(players, p)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return players, nil

}

func (r repository) Delete(id string) error {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, err := r.client.Database(database).Collection("players").
		DeleteOne(context.Background(), bson.M{"_id": objId})
	return err
}

func (r repository) Upload(id string, bin []byte) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	img := Image{
		Data: bin,
	}

	player := canonical.Player{}

	r.client.Database(database).Collection("players").
		FindOne(context.Background(), bson.M{"_id": objId}).Decode(&player)

	if player.ImgId == "" {
		result, err := r.client.Database(database).Collection("images").
			InsertOne(context.Background(), &img)
		if err != nil {
			return err
		}

		imgId := ""
		if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
			imgId = oid.Hex()
		}

		_, err = r.client.Database(database).Collection("players").
			UpdateOne(context.Background(), bson.M{"_id": objId}, bson.D{{Key: "$set", Value: bson.M{"imgid": imgId}}})
		if err != nil {
			return err
		}

	} else {

		imgObjId, err := primitive.ObjectIDFromHex(player.ImgId)
		if err != nil {
			return err
		}
		_, err = r.client.Database(database).Collection("images").
			UpdateOne(context.Background(), bson.M{"_id": imgObjId}, bson.D{{Key: "$set", Value: bson.M{"data": bin}}})
	}

	return err
}

func (r repository) Download(id string) ([]byte, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	player := canonical.Player{}

	err = r.client.Database(database).Collection("players").
		FindOne(context.Background(), bson.M{"_id": objId}).Decode(&player)
	if err != nil {
		return nil, err
	}

	imgId, err := primitive.ObjectIDFromHex(player.ImgId)
	if err != nil {
		return nil, err
	}

	img := Image{}
	err = r.client.Database(database).Collection("images").
		FindOne(context.Background(), bson.M{"_id": imgId}).Decode(&img)

	return img.Data, err
}
