package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	DB "shadowflade/timers/pkg/db"
	"shadowflade/timers/pkg/global"
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

	userId, err = strconv.Atoi(userIdVal)

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

	newTimer := interfaces.NewTimer(int32(userID), "your mom", "red")

	newTimerID, err := db.CreateTimer(newTimer)

	if err != nil {
		panic(err.Error())
	}

	err = templates.ExecuteTemplate(w, "timer", map[string]interface{}{
		"userID": userID,
		"id":     newTimerID,
	})
	if err != nil {
		fmt.Print(err.Error(), userID, newTimerID)
	}
}

func (this *TimerHandler) PauseTimer(w http.ResponseWriter, r *http.Request) {
	db := DB.Db{}

	body, _ := io.ReadAll(r.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	timerId := response["timer_id"]
	fmt.Println(timerId, " TImer id")
	if timerId == 0 || timerId == nil {
		return
	}
	timerId, _ = strconv.Atoi(timerId.(string))

	newDuration, _ := db.PauseTimer(timerId.(int))
	w.Write([]byte(string(newDuration)))

}

func (this *TimerHandler) UpdateTimerTitle(w http.ResponseWriter, r *http.Request) {

}

func (this *TimerHandler) DeleteTimer(w http.ResponseWriter, r *http.Request) {

}

func (this *TimerHandler) UpdateTimer(w http.ResponseWriter, r *http.Request) {

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
