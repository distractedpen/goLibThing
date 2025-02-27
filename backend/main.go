package main

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"

	"gamelib.cloud/models"
	"gamelib.cloud/services"
)

func handleGetGames(w http.ResponseWriter, r *http.Request, s *services.Service) {
	resp := make(map[string]string)
	gamesResults, err := s.GetGamesService(r.Context())
	if nil != err {
		resp["error"] = "Error getting games list"
		resp["status"] = "error"
		log.Printf("%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
    }

	// convert result into json
	w.Header().Set("Content-Type", "application/json")
	gamesJson, err := json.Marshal(gamesResults)
	if nil != err {
		resp["error"] = "JSON ERROR"
		resp["status"] = "error"
		log.Printf("JSON ERROR\n")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp["data"] = string(gamesJson)
		resp["status"] = "success"
		log.Printf("Get Success\n")
		w.WriteHeader(http.StatusOK)
	}
	jsonRep, _ := json.Marshal(resp)
	w.Write(jsonRep)
}

func handlePostGame(w http.ResponseWriter, r *http.Request, s *services.Service) {
	defer r.Body.Close()
	resp := make(map[string]string)

	bodyBytes, err := io.ReadAll(r.Body)
	if nil != err {
		log.Printf("Error Reading Response Body\n")
		resp["status"] = "error"
		resp["error"] = "Error Reading Response Body"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
        return
	}

	if len(bodyBytes) == 0 {
		log.Printf("Game Data is required\n")
		resp["status"] = "error"
		resp["error"] = "Game Data is required"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
        return
	}

	var newGameData models.NewGameData
	err = json.Unmarshal(bodyBytes, &newGameData)
	if nil != err {
		log.Printf("JSON ERROR\n")
		resp["status"] = "error"
		resp["error"] = "JSON Error"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResp)
        return
	}

	newGame, err := s.AddGameService(r.Context(), newGameData)
	if nil != err {
		log.Printf("Error writing to database\n")
		resp["status"] = "error"
		resp["error"] = "Error writing to database"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResp)
        return
	}

	newGameJson, _ := json.Marshal(newGame)
	log.Printf("POST Success")
	resp["status"] = "success"
	resp["data"] = string(newGameJson)
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func main() {

	// loading environment
	if err := godotenv.Load(); nil != err {
		log.Printf("No .env file found.\n")
	}
	log.Printf("Loaded env.\n")

	// connect to database
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URI"))
	if nil != err {
		panic(err)
	}
	defer conn.Close(context.Background())

	log.Printf("Connected to Database.\n")
	s := services.Service{Db: conn}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /games", func(w http.ResponseWriter, r *http.Request) {
		handleGetGames(w, r, &s)
	})
	// mux.HandleFunc("GET /games/{id}", handleGetGameById)
	mux.HandleFunc("POST /games", func(w http.ResponseWriter, r *http.Request) {
		handlePostGame(w, r, &s)
	})
	// mux.HandleFunc("PUT /games/{id}", handlePutGame)
	// mux.HandleFunc("DELETE /games/{id}", handleDeleteGame)

	log.Printf("Listening on :8080\n")
	http.ListenAndServe(":8080", mux)
}
