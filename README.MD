# Golang rest application integration with mongodb

This project is an example how create rest application and execute CRUD operations with mongoDB and defult driver

For simulate a scenario, we have a register players with schema bellow

```json
{
    "name": "Ronaldinho",
    "age": 40,
    "position": "st",
    "foot": "right",
    "genre": "male",
    "rating": 4.9,
    "country": "BRA",
    "url": "https://your-url.com"
}
```
Using MongoDB in docker-compose

```yaml
# Use root/example as user/password credentials
version: '3.1'

services:

  mongo:
    image: mongo
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: example
    ports:
      - 27017:27017

  mongo-express:
    image: mongo-express
    restart: always
    ports:
      - 8081:8081
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: root
      ME_CONFIG_MONGODB_ADMINPASSWORD: example

```

Supports operations:
* Add player
* Get player
* Get all players
* Get all player with pagination (offset, limit)
* Delete player
* Upload image player
* Download image player

Endpoints:
* POST   /players
* GET    /players/{id}
* GET    /players?offset=0&limit=100
* DELETE /players/{id}
* POST   /players/{id}/images
* GET    /players/{id}/images

### Requirements
* Go
* Docker
* Docker Compose
* Terminal

For run this project:
```shell
    git clone https://github.com/ArturMartini/go-demo-mongodb.git
    cd go-demo-mongodb
    go run main.go
    sh data.sh
```

Now you can explore the endpoints

