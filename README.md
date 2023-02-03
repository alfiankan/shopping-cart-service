# Golang Shopping Cart Service API

## How to test
This application uses unit testing to carry out tests, for easier integration testing this application uses dockertest. Once you've run your tests, dockertest spins up all the dependencies/infrastructure needed to run your tests and cleanup when done, so tests are tested on real infrastructure not mocks.

To run integration test use this command :

```bash
make test
```

If you just want to run http test endpoint use this command :

## How to run
- first copy .env.example to .env
- then run makefile script `make docker` then `make migrate` then `make run`

### Makefile script
this project conatains Makefile :

|command|description|
|---|---|
| make migrate | migrate up database to writedb |
| make migrate-down | migrate down database from writedb |
| make seed | seed/insert sample data to system |
| make swagger | rebuild swagger api docs |
| make docker | compose up minimum dependencies (all single node) |
| make run | run api |



## API Docs
Open swagger ui at `http://localhost:3000/swagger/index.html`

