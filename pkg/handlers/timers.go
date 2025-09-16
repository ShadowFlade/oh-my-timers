package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	DB "shadowflade/timers/pkg/db"
	"shadowflade/timers/pkg/global"
	"shadowflade/timers/pkg/interfaces"
	"shadowflade/timers/pkg/services"
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
	if err != nil || userIdCookie.Value == "0" {
		templates.ExecuteTemplate(w, "index", interfaces.TimerTemplate{
			Items:                   make([]interfaces.Timer, 0),
			IsMoreThan10:            false,
			UserID:                  0,
			ShowNewUserAlertTrigger: false,
		})
	} else {
		userIdVal = userIdCookie.Value
		// fmt.Printf("%s cookie user id ", userIdVal)
	}

	userId, err = strconv.Atoi(userIdVal)

	userTimers := timersDb.GetAllUsersTimers(userId)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	templates.ExecuteTemplate(w, "index", interfaces.TimerTemplate{
		Items:                   userTimers,
		IsMoreThan10:            false,
		UserID:                  userId,
		ShowNewUserAlertTrigger: false,
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
	userTimers := make([]interfaces.Timer, 1, 1)
	userTimers = append(userTimers, newTimer)

	newTimerID, err := db.CreateTimer(newTimer)

	if err != nil {
		panic(err.Error())
	}

	templates.ExecuteTemplate(w, "timer", newTimer)
	if err != nil {
		fmt.Print(err.Error(), userID, newTimerID)
	}
}

func (this *TimerHandler) StartTimer(w http.ResponseWriter, r *http.Request) {
	db := DB.Db{}
	body, _ := io.ReadAll(r.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	timerId := response["timer_id"]

	if timerId == 0 || timerId == nil {
		return
	}

	timerId, _ = strconv.Atoi(timerId.(string))

	affectedRows, err := db.StartTimer(timerId.(int))
	if err != nil {
		log.Panic(err.Error())
	}
	resp, _ := json.Marshal(map[string]interface{}{
		"isSuccess": true,
		"data": map[string]interface{}{
			"affectedRows": affectedRows,
		},
	})
	w.Write(resp)
}

func (this *TimerHandler) PauseTimer(w http.ResponseWriter, r *http.Request) {
	db := DB.Db{}
	body, _ := io.ReadAll(r.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	timerId := response["timer_id"]
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

func (this *TimerHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var response map[string]string
	json.Unmarshal(body, &response)
	password := response["password"]

	if password == "" {
		resp, _ := json.Marshal(map[string]interface{}{
			"error":     "Password is empty. Could not create user",
			"isSuccess": false,
		})
		w.Write(resp)
	}

	userService := services.User{}
	hashedPassword, _ := userService.HashPassword(password)

	userDb := DB.User{}
	user := interfaces.NewUser("USER", hashedPassword)
	newUserID := userDb.CreateUser(user)
	cookie := &http.Cookie{
		Name:     "user_id",
		Value:    strconv.Itoa(int(newUserID)),
		Path:     "/",
		MaxAge:   3600 * 24 * 7,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	detectCookie := &http.Cookie{
		Name:     "user_id_detected",
		Value:    "1", // Simple flag
		Path:     "/",
		MaxAge:   3600 * 24 * 7,
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, detectCookie)

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	newUserResp, _ := json.Marshal(map[string]interface{}{
		"newUserId": newUserID,
		"isSuccess": true,
	})
	w.Write(newUserResp)
}
