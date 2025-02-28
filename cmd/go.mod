module gamelib.cloud/main

go 1.24.0

require (
	gamelib.cloud/game/models v0.0.0-00010101000000-000000000000
	gamelib.cloud/game/service v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v5 v5.7.2
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

replace gamelib.cloud/game/service => ../game/service

replace gamelib.cloud/game/models => ../game/models
