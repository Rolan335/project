.PHONY: migrate-add

migrate-add: ## Create new migration file, usage: migrate-add [name=<migration_name>]
	goose -dir migrations create $(name) sql