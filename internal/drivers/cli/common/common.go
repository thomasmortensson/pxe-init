package common

import (
	"fmt"
	"os"
	"strconv"
)

func ErrorPrint(msg string, err error) {
	fmt.Fprintln(os.Stderr, msg, err)
}

func ErrorExit(msg string, err error) {
	ErrorPrint(msg, err)
	os.Exit(-1)
}

func GetEnvOrDefaultString(envVar, defaultStr string) string {
	result := os.Getenv(envVar)
	if result == "" {
		return defaultStr
	}
	return result
}

func GetEnvOrDefaultInt(envVar string, defaultInt int) int {
	result := os.Getenv(envVar)
	if result == "" {
		return defaultInt
	}
	integer, err := strconv.ParseInt(result, 10, 64) // nolint:gomnd // Fairly intuitive that these are a base and bitwidth respectively
	if err != nil {
		return defaultInt
	}
	return int(integer)
}
