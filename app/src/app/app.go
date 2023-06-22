package app

import (
	"itmo-profile/internal/gin"
	"log"
)

func Run() {

	r := gin.Create()
	err := r.Run(":5242")
	if err != nil {
		return
	}
	// GIN END
	log.Println(r)
	//log.Println(conf)

}
