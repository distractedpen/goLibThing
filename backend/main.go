package main

// TODO: Convert from Fprintf to json.Encoder

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
)

// Dto
type NewGameData struct {
	Name      string `json:"name"`
	Developer string `json:"developer"`
}
type UpdateGameData = NewGameData

// Database Object
type Game struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Developer string `json:"developer"`
}

func (g *NewGameData) CreateNewGame() Game {
	return Game{
		Id:        currIndex,
		Name:      g.Name,
		Developer: g.Developer,
	}
}

func (g *UpdateGameData) MapToGame(id int) Game {
    return Game{
        Id: id,
        Name: g.Name,
        Developer: g.Developer,
    }
}

var GAMES []Game
var currIndex int

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, gamer!\n")
}

func handleGetGames(w http.ResponseWriter, r *http.Request) {

	gamesJson, err := json.Marshal(GAMES)
	if nil != err {
		fmt.Fprintf(w, "JSON Error\n")
		return
	}
	fmt.Fprintf(w, "%s\n", string(gamesJson))
}

func handleGetGameById(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if nil != err {
		fmt.Fprintf(w, "Id not valid.\n")
		return
	}

	index := slices.IndexFunc(GAMES, func(g Game) bool {
		return g.Id == id
	})

	if index == -1 {
		fmt.Fprintf(w, "No game found with id %s\n", id)
		return
	}

	game := GAMES[index]
	gameJson, err := json.Marshal(game)
	if nil != err {
		panic(err)
	}
	fmt.Fprintf(w, "%s\n", string(gameJson))
}

func handlePostGame(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if nil != err {
		fmt.Fprintf(w, "Error reading body\n")
		return
	}

	if len(bodyBytes) == 0 {
		fmt.Fprintf(w, "Error: Must provide game data.\n")
		return
	}

	fmt.Printf("%s", bodyBytes)

	var newGameData NewGameData
	err = json.Unmarshal(bodyBytes, &newGameData)
	if nil != err {
		fmt.Fprintf(w, "Error: JSON Unmarshal Error.")
		panic(err)
	}

	newGame := newGameData.CreateNewGame()

	GAMES = append(GAMES, newGame)
	currIndex++

	fmt.Fprintf(w, "Added a new game with details: %d, %s, %s\n",
		newGame.Id, newGame.Name, newGame.Developer)
}

// PathValue: Id
// Body: { name string | nil, developer string | nil }
func handlePutGame(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if nil != err {
		fmt.Fprintf(w, "Id not valid.\n")
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if nil != err {
		fmt.Fprintf(w, "Error: reading body.\n")
		return
	}

	if len(bodyBytes) == 0 {
		fmt.Fprintf(w, "Error: Must provide game data.\n")
		return
	}

    fmt.Printf("%s", bodyBytes)

	var updatedGameData UpdateGameData
	err = json.Unmarshal(bodyBytes, &updatedGameData)
	if nil != err {
		// fmt.Fprintf(w, "Error: JSON Unmarshal Error.")
        panic(err)
		return
	}

    updatedGame := updatedGameData.MapToGame(id)

	index := slices.IndexFunc(GAMES, func(g Game) bool {
		return g.Id == id
	})

    GAMES = slices.Replace(GAMES, index, index+1, updatedGame)

	fmt.Fprintf(w, "Updated game with id %s. Name: %s and Developer: %s.\n",
		id, updatedGame.Name, updatedGame.Developer)
}

// PathValue: id
func handleDeleteGame(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if nil != err {
		fmt.Fprintf(w, "Id not valid.\n")
		return
	}

	GAMES = slices.DeleteFunc(GAMES, func(g Game) bool {
		return g.Id == id
	})

	fmt.Fprintf(w, "Delete an existing game with id %s.\n", id)
}

func main() {
	GAMES = make([]Game, 0)
	currIndex = 0
	mux := http.NewServeMux()
	mux.HandleFunc("GET /games", handleGetGames)
	mux.HandleFunc("GET /games/{id}", handleGetGameById)
	mux.HandleFunc("POST /games", handlePostGame)
	mux.HandleFunc("PUT /games/{id}", handlePutGame)
	mux.HandleFunc("DELETE /games/{id}", handleDeleteGame)
	mux.HandleFunc("/{$}", handleRoot)
	fmt.Printf("Listening on :8080\n")
	http.ListenAndServe(":8080", mux)
}
