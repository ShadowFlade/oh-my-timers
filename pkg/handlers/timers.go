package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"shadowflade/timers/global"
	DB "shadowflade/timers/pkg/db"
	"shadowflade/timers/pkg/interfaces"
	"strconv"
)

type TimerHandler struct {
}

func (this *TimerHandler) RenderUserTimers(w http.ResponseWriter, r *http.Request) {

	timersDb := DB.Timer{}
	var userId int
	var userIdVal string

	userIdCookie, err := r.Cookie(global.COOKIE_USER_ID_NAME)
	if err != nil {
		userIdVal = "0"
	} else {
		userIdVal = userIdCookie.Value
	}

	userId, err = strconv.Atoi(userIdVal)
	fmt.Print(userId)

	templates, err := template.ParseGlob("views/*.html")
	newTemplates := template.Must(templates, err)

	if userId == 0 {
		userId = int(this.createUser("USER", w))
		templates.ExecuteTemplate(w, "index", make([]interfaces.Timer, 0))
	}
	userTimers := timersDb.GetAllUsersTimers(userId)
	if err != nil {
		panic(err.Error())
	}
	log.Print(userTimers, newTemplates)
	templates.ExecuteTemplate(w, "index", userTimers)

}

func (this *TimerHandler) Create(w http.ResponseWriter, r *http.Request) {
	db := DB.Db{}
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

func (this *TimerHandler) createUser(name string, w http.ResponseWriter) int64 {
	userDb := DB.User{}
	newUserID := userDb.CreateUser(name)
	cookie := http.Cookie{Name: "userID", Value: string(newUserID)}
	http.SetCookie(w, &cookie)
	return newUserID
}
