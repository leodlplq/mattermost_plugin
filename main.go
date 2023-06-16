package main

import (
	"github.com/leodlplq/codingame/codingame"	
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func main(){
	db, err := sqlx.Connect("mysql", "leo:leo@(localhost)/mattermost_plugin")
	if err != nil {
			log.Panic().Err(err).Msg("Could not connect to DB")
	}

	r := gin.Default()
	r.GET("/ping", codingame.Ping())
	r.POST("/login", codingame.Login(db))
	r.POST("/launchGame", codingame.LaunchGame(db))
	
	r.Run() // listen and serve on 0.0.0.0:8080
}