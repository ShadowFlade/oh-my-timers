# DISCLAIMER

This code style is disapproved by golang community and myself (i just love `this`)

### Start project
`air`
build migration binary, e.g:
`cd migrations`
`go build -o goose-custom *.go`
then run it
`./goose-custom mysql up`

## TODO

### v1.0
* ~~deleting timer~~
* ~~update title~~
* ~~разобраться почему при перезагрузке страницы показывает таймер запускающийся заново каждый раз~~
* ~~сделать чтобы при удалении таймеры не сразу меняли положение, а было как в firefox - спустя 1.5 секунды~~ (такой херней заниматься нет желания - не сделано)
* add report button in dev mode - ? (the idea is to accept report messages and send user id and all his data so i can replicate bug - mb its too much for a fucking timer app)
* ~~добавить выбор цвета?~~ 
* добавить в миграции поддержку sqlite - ?


### fixes
* when starting (clicking start) timer with high latency (vpn) it does not start the timer untill the request finishes - NEED TO FIX
* optimize assets - minimize css and js on build (on push/pull) - fix pls

### >v1.0
* add commit timer time to categories - they will sum up and reset timer (and start timer if it was active at the moment) - branch categories - actually, not sure if i wanna do it
* b24 integration - add side menu for entering b24hook url which will be saved into cookies -add button "load tasks" which will fetch your active tasks (assigned to you? - so also should save to cookies your id or smth) and make timers out of it