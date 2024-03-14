package doghole

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"
)

func PrepareLog() error {
	var err error
	err = PrepareLogDir()
	if err != nil {
		return err
	}
	err = PrepareLogFile()
	if err != nil {
		return err
	}
	return nil
}

func PrepareLogDir() error {
	if _, err := os.Stat("./log"); os.IsNotExist(err) {
		err := os.Mkdir("./log", 0664)
		if err != nil {
			return errors.New("failed to create log dir: " + err.Error())
		}
	}
	return nil
}

func PrepareLogFile() error {
	logFileName := "./log/" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".log"
	logFile, err := os.OpenFile(logFileName, syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND|syscall.O_SYNC, 0664)
	if err != nil {
		return errors.New("failed to create log file " + logFileName + ": " + err.Error())
	}
	logWriter := io.MultiWriter(os.Stdout, logFile)
	log.Default().SetOutput(logWriter)
	return nil
}
