package handlers

import (
	"encoding/json"
	"net/http"
	"shadowflade/timers/pkg/db"
	"shadowflade/timers/pkg/interfaces"
)

type TimerHandler struct {
}

func (this *TimerHandler) Create(w http.ResponseWriter, r *http.Request) {
	db := db.Db{}
	var body []byte
	r.Body.Read(body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)

	timerStart, isTimerStartOk := response["start"].(string)
	timerEnd, isTimerEndOk := response["end"].(string)
	userID, isUserIDOk := response["userId"].(int)

	handlerThoseFuckingErrors(isTimerStartOk, isTimerEndOk, isUserIDOk)

	timer := interfaces.Timer{
		Start:  timerStart,
		End:    timerEnd,
		UserID: int32(userID),
	}

	newId, err := db.CreateTimer(timer)
	if err != nil {
		panic(err.Error())
	}
	result := map[string]interface{}{
		"timerId": newId,
	}
	responseJson, err := json.Marshal(result)
	if err != nil {
		panic(err.Error())
	}

	w.Write(responseJson)
}

func handlerThoseFuckingErrors(isTimerStartOk bool, isTimerEndOk bool, isUserIDOk bool) {
	if !isTimerStartOk {
		panic("No timer start date provided")
	}
	if !isTimerEndOk {
		panic("No timer end date provieded")
	}
	if !isUserIDOk {
		panic("No use id provided")
	}
}
