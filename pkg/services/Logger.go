package services

import (
	"log"
	"os"
)

type Logger struct {
}

func (this *Logger) LogText(message string, filename string) {
	if filename == "" {
		filename = "log.log"
	}
	fs, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		log.Panic(err.Error())
	}
	defer func() {
		if err := fs.Close(); err != nil {
			log.Panic(err)
		}
	}()

	line := message + "\n"
	bytesN, err := fs.Write([]byte(line))

	if err != nil || bytesN == 0 {
		logger := log.Default()
		logger.Println(err.Error())
	}

}
