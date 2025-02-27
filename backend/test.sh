#!/bin/bash
GAME1='{"name": "Final Fantasy XIV", "developer": "Creative Business Unit 3"}'
GAME2='{"name": "Crash Bandicoot", "developer": "Naughty Dog"}'
GAME3='{"name": "Nier Automata", "developer": "Platinum Games"}'
UPDATEGAME='{"name": "Crash Bandicoot 2", "developer": "Naughty Dog"}'
curl -X POST --data "$GAME1" "localhost:8080/games"
curl -X POST --data "$GAME2" "localhost:8080/games"
curl -X POST --data "$GAME3" "localhost:8080/games"
curl -X GET "localhost:8080/games"
curl -X PUT --data "$UPDATEGAME" "localhost:8080/games/1"
