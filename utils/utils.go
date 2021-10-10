package utils

import "log"

func ErrorHandle(err error, msg string) {
	if err != nil {
		log.Fatalf("%s %v", msg, err)
	}
}

func IsAvailable(arr []string, item string) int {
	for index, value := range arr {
		if value == item {
			return index
		}
	}

	return -1

}
