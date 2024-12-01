package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

var Log = zerolog.New(zerolog.ConsoleWriter{})

// SetupLogger :
func SetupLogger(logger zerolog.Logger) error {
	var level, filePath string

	level = "debug"
	fmt.Println(level)

	// filePath = "F:\\Novanew\\novajobs-ultra-aura-development\\trainers\\logger\\logger.log"
	filePath = "F:\\novaultraaura\\novajobs-ultra-aura-shyamal2\\trainers\\logger\\logger.log"

	if filePath == "" {
		return errors.New("logger file path not found")
	}

	basePath := filepath.Dir(filePath)
	created, errStr := checkPathExists(basePath, 1)
	if !created {
		return errors.New(errStr)
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Log = zerolog.New(zerolog.ConsoleWriter{Out: file, NoColor: false, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	return nil
}

// checkPathExists :
func checkPathExists(dirDump string, createDirDump int) (bool, string) {
	var err error
	if _, err = os.Stat(dirDump); err != nil {
		if os.IsNotExist(err) {
			if createDirDump == 1 {
				err = os.MkdirAll(dirDump, 0700)
				if err != nil {
					errmsg := fmt.Sprintf("Err: Path to backup \"%s\". %s", dirDump, err.Error())
					return false, errmsg
				} else {
					return true, ""
				}
			} else {
				errmsg := fmt.Sprintf("Err: Path to backup \"%s\" doesn't exists", dirDump)
				return false, errmsg
			}
		} else {
			errmsg := fmt.Sprintf("Err: Path to backup \"%s\" doesn't exists", dirDump)
			return false, errmsg
		}
	} else {
		return true, ""
	}
}
