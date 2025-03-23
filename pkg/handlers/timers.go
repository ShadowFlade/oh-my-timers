package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shadowflade/timers/global"
	DB "shadowflade/timers/pkg/db"
	"shadowflade/timers/pkg/interfaces"
	"shadowflade/timers/pkg/views"
	"strconv"
)

type TimerHandler struct {
}

func (this *TimerHandler) RenderUserTimers(w http.ResponseWriter, r *http.Request) {

	timersDb := DB.Timer{}
	views := views.Views{}

	templates := views.GetTemplates()
	var userId int
	var userIdVal string

	userIdCookie, err := r.Cookie(global.COOKIE_USER_ID_NAME)
	if err != nil {
		userIdVal = "0"
	} else {
		userIdVal = userIdCookie.Value
	}

	fmt.Print(userIdCookie.Value, " COOKIE USER ID!!")

	userId, err = strconv.Atoi(userIdVal)
	fmt.Print(userId)

	if userId == 0 {
		userId = int(this.createUser("USER", w))
		http.SetCookie(w, &http.Cookie{Name: global.COOKIE_USER_ID_NAME, Value: strconv.Itoa(userId)})
		templates.ExecuteTemplate(
			w,
			"index",
			interfaces.TimerTemplate{
				Items:        make([]interfaces.Timer, 0),
				IsMoreThan10: false,
				UserID:       userId,
			},
		)
	}
	userTimers := timersDb.GetAllUsersTimers(userId)

	templates.ExecuteTemplate(w, "index", interfaces.TimerTemplate{
		Items:        userTimers,
		IsMoreThan10: false,
		UserID:       userId,
	})

}

func (this *TimerHandler) CreateTimer(w http.ResponseWriter, r *http.Request) {
	db := DB.Db{}
	views := views.Views{}
	templates := views.GetTemplates()
	cookie, err := r.Cookie(global.COOKIE_USER_ID_NAME)
	if err != nil {
		panic(err.Error())
	}
	cookieVal := cookie.Value
	userID, err := strconv.Atoi(cookieVal)
	if err != nil {
		panic(err.Error())
	}
	newTimer := interfaces.Timer{
		UserID: int32(userID),
	}
	newTimerID, err := db.CreateTimer(newTimer)

	if err != nil {
		panic(err.Error())
	}

	templates.ExecuteTemplate(w, "timer", map[string]interface{}{
		"userId": userID,
		"id":     newTimerID,
	})
}

func (this *TimerHandler) UpdateTimer(w http.ResponseWriter, r *http.Request) {
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
	cookie := http.Cookie{Name: "userID", Value: strconv.Itoa(int(newUserID))}
	http.SetCookie(w, &cookie)
	return newUserID
}
