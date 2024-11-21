package helper

import (
	"errors"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func ErrMsg(err error, msg string) error {
	if err != nil {
		logger.Error(msg, err)
		return errors.New(msg + err.Error())
	}
	return nil
}
func SuccessMsg(msg string){
	logger.Info(msg)
}