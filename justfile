migrate:
   go build -o goose-custom ./migrations
   goose-custom mysql up
