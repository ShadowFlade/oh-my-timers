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
* ~~deleting timer~~
* ~~update title~~
* ~~разобраться почему при перезагрузке страницы показывает таймер запускающийся заново каждый раз~~
* ~~сделать чтобы при удалении таймеры не сразу меняли положение, а было как в firefox - спустя 1.5 секунды~~ (такой херней заниматься нет желания - не сделано)
* add report button in dev mode
* добавить выбор цвета?
* добавить в миграции поддержку sqlite
