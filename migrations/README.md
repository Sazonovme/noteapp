[DOCS FOR MIGRATIONS]

1. Создать новый файл миграции
migrate create -ext sql -dir migrations -seq -digits 4 <MIGRATION-NAME>

2. Использовать миграцию
migrate -path <PATH> -database <DATABASE-PATH> up

[EXAMPLE]
migrate -path migrations -database "postgres://localhost/notesdb?sslmode=disable" up