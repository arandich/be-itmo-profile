package tools

import (
	"log"
	"strconv"
)

func ConvertStrToInt(str string) int {
	number, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
		return 0
	}
	return number
}
