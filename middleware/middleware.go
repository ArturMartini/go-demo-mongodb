package middleware

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go-demo-mongodb/canonical"
	"go-demo-mongodb/service"
	"io/ioutil"
	"net/http"
	"strconv"
)

var svc service.Service

func init() {
	if svc == nil {
		svc = service.NewService()
	}
}

func Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/players", addPlayer).Methods(http.MethodPost)
	router.HandleFunc("/players/{id}", getPlayer).Methods(http.MethodGet)
	router.HandleFunc("/players", getAllPlayers).Methods(http.MethodGet)
	router.HandleFunc("/players/{id}", updatePlater).Methods(http.MethodPut)
	router.HandleFunc("/players/{id}", deletePlayer).Methods(http.MethodDelete)
	router.HandleFunc("/players/{id}/images", uploadImage).Methods(http.MethodPost)
	router.HandleFunc("/players/{id}/images", downloadImage).Methods(http.MethodGet)
	return http.ListenAndServe(":8080", router)
}

func downloadImage(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	bin, err := svc.Download(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(bin)
	w.WriteHeader(http.StatusOK)
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if header.Size > (1024 * 50) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"The file size is greater than 500kb\"}"))
	}

	bin, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := mux.Vars(r)["id"]

	err = svc.Upload(id, bin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func addPlayer(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	player := canonical.Player{}
	err = json.Unmarshal(bytes, &player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = svc.Add(&player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unexpected error ocurred"))
		return
	}

	body, err := json.Marshal(player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func updatePlater(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	player := canonical.Player{}
	err = json.Unmarshal(bytes, &player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	player.Id = id

	err = svc.Update(&player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unexpected error ocurred"))
		return
	}

	body, err := json.Marshal(player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	player, err := svc.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unexpected error ocurred"))
		return
	}

	if player.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := json.Marshal(player)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func getAllPlayers(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	players, err := svc.GetAll(offset, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unexpected error ocurred"))
		return
	}

	body, err := json.Marshal(players)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func deletePlayer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := svc.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unexpected error ocurred"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
