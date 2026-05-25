.PHONY: all
all: bot

.PHONY: check
check:
	@$(if $(wildcard .env),,echo Error: .env file not found. Please create it first. && exit 1)

.PHONY: bot
bot: check
	go run .\cmd\bot