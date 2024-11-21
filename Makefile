export DATABASE_URL="mysql://root:secret@tcp(localhost:3306)/fs_forum_api?query"

migrate-create:
	@ migrate create -ext sql -dir scripts/migrations -seq $(name)

migrate-up:
	@ migrate -database ${DATABASE_URL} -path scripts/migrations up

migrate-down:
	@ migrate -database ${DATABASE_URL} -path scripts/migrations down