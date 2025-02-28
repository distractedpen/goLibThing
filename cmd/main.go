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
    "strconv"

	"gamelib.cloud/game/models"
	"gamelib.cloud/game/service"
)

func handleGetGames(w http.ResponseWriter, r *http.Request, s *service.Service) {
	resp := make(map[string]any)
	gamesResults, err := s.GetGamesService(r.Context())
	if nil != err {
		resp["error"] = "Error getting games list"
		resp["status"] = "error"
		log.Printf("%s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
    }

	// convert result into json
	w.Header().Set("Content-Type", "application/json")
	if nil != err {
		resp["error"] = "JSON ERROR"
		resp["status"] = "error"
		log.Printf("JSON ERROR\n")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp["data"] = gamesResults
		resp["status"] = "success"
		log.Printf("Get Success\n")
	}
    log.Printf("GET Success\n")
	jsonRep, _ := json.Marshal(resp)
	w.Write(jsonRep)
}

func handleGetGameById(w http.ResponseWriter, r *http.Request, s *service.Service) {
    resp := make(map[string]any)

    idRaw := r.PathValue("id")
    
    id, err := strconv.ParseInt(idRaw, 10, 64)
    if nil != err {
		log.Printf("Error Parsing id\n")
		resp["status"] = "error"
		resp["error"] = "Error Parsing Id"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
        return
    }

    game, err := s.GetGameByIdService(r.Context(), id)
    if nil != err {
		log.Printf("Server Error\n")
		resp["status"] = "error"
		resp["error"] = "Server Error"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
        return
    }

    log.Printf("GET Success\n")
    resp["data"] = game
    jsonResp, _ := json.Marshal(resp)
    w.Write(jsonResp)

}

func handlePostGame(w http.ResponseWriter, r *http.Request, s *service.Service) {
	defer r.Body.Close()
	resp := make(map[string]any)

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

	log.Printf("POST Success\n")
	resp["status"] = "success"
	resp["data"] = newGame
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func handleUpdateGame(w http.ResponseWriter, r *http.Request, s *service.Service) {
	defer r.Body.Close()
	resp := make(map[string]any)

    idRaw := r.PathValue("id")
    
    id, err := strconv.ParseInt(idRaw, 10, 64)
    if nil != err {
		log.Printf("Error Parsing Id\n")
		resp["status"] = "error"
		resp["error"] = "Error Parsing Id"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
        return
    }

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

	var updateGameData models.UpdateGameData 
	err = json.Unmarshal(bodyBytes, &updateGameData)
	if nil != err {
		log.Printf("JSON ERROR\n")
		resp["status"] = "error"
		resp["error"] = "JSON Error"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResp)
        return
	}

	newGame, err := s.UpdateGameService(r.Context(), id, updateGameData)
	if nil != err {
		log.Printf("Error writing to database\n")
		resp["status"] = "error"
		resp["error"] = "Error writing to database"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResp)
        return
	}

	log.Printf("UPDATE Success\n")
	resp["status"] = "success"
	resp["data"] = newGame
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func handleDeleteGame(w http.ResponseWriter, r *http.Request, s *service.Service) {
    resp := make(map[string]any)
    idRaw := r.PathValue("id")
    id, err := strconv.ParseInt(idRaw, 10, 64)
    if nil != err {
		log.Printf("Error Parsing Id\n")
		resp["status"] = "error"
		resp["error"] = "Error Parsing Id"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
        return
    }

    game, err := s.DeleteGameService(r.Context(), id)
    if nil != err {
		log.Printf("Server Error\n")
		resp["status"] = "error"
		resp["error"] = "Server Error"
		jsonResp, _ := json.Marshal(resp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)
        return
    }

    
	log.Printf("DELETE Success\n")
	resp["status"] = "success"
	resp["data"] = game 
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
	s := service.Service{Db: conn}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /games", func(w http.ResponseWriter, r *http.Request) {
		handleGetGames(w, r, &s)
	})
	mux.HandleFunc("GET /games/{id}", func(w http.ResponseWriter, r *http.Request) {
		handleGetGameById(w, r, &s)
	})
	mux.HandleFunc("POST /games", func(w http.ResponseWriter, r *http.Request) {
		handlePostGame(w, r, &s)
	})
	mux.HandleFunc("PUT /games/{id}", func(w http.ResponseWriter, r *http.Request) {
		handleUpdateGame(w, r, &s)
	})
	mux.HandleFunc("DELETE /games/{id}", func(w http.ResponseWriter, r *http.Request) {
		handleDeleteGame(w, r, &s)
	})

	log.Printf("Listening on :8080\n")
	http.ListenAndServe(":8080", mux)
}
