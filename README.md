The idea of EFOS checker is to provide a way for journalist and investigators a way to check a list of entities that might be present on the 69b list from SAT.

- This will keep an up-to date database of definitive EFOS.

## How to use

### Requirements

- Docker
- Golang (>1.13)

### Dev

1. Run `cp .env.example .env`
1. Run `make composeup` in your terminal
1. Run `make migrateup` once database is ready
1. Run `make downloadefos`
1. Add list of company names to `nombres.csv`
1. Run `make searchefos`
1. See results on screen
1. Run `make composestop`
