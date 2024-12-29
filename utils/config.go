package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func getConfig(name string) (string, bool) {
	return os.LookupEnv(fmt.Sprintf("LUDIVAULT_%s", strings.ToUpper(name)))
}

func GetOptionalConfig(name string, fallback string) string {
	value, ok := getConfig(name)
	if !ok {
		return fallback
	}
	return value
}

func GetRequiredConfig(name string) string {
	value, ok := getConfig(name)
	if !ok {
		log.Panicf("Config %s not found", name)
	}
	return value
}
