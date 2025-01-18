package utils

import (
	log "github.com/sirupsen/logrus"
)

func HandleError(err error, causeFatal bool) {
	if err != nil {
		if causeFatal {
			log.Fatal(err.Error())
		} else {
			log.Error(err.Error())
		}
	}
}
