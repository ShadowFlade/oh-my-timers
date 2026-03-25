set dotenv-load

# build migration binary and runs migration
migrate direction:
   go build -o goose-custom ./migrations
   ./goose-custom mysql {{direction}}

#[parallel]
dev:
    air
    xdg-open "localhost:$APP_PORT"

#url := "localhost:8069"
#open:
#    xdg-open {{url}}
