.PHONY: migrate-add, gotopostgres

migrate-add: ## Create new migration file, usage: migrate-add [name=<migration_name>]
	goose -dir migrations create $(name) sql

gotopostgres:
	docker exec -it project-postgres-1 psql -U postgres -d blog