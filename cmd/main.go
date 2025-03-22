package cmd

import (
	"log"
	"net/http"
	"shadowflade/timers/global"
	"shadowflade/timers/pkg/handlers"
)

func main() {

	http.Handle("createTimer",handlers.TimerHandler)
	log.Fatal(http.ListenAndServe(":" + global.PORT , nil))
}
func ServerHandler() {

}

