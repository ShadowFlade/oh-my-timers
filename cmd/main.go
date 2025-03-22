package cmd

import (
	"log"
	"net/http"
	"shadowflade/timers/global"
)

func main() {

	http.Handle("createTimer",
	log.Fatal(http.ListenAndServe(":" + global.PORT , nil))
}
