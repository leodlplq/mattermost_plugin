package codingame

import (
	// "database/sql"
	"strings"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
  "github.com/rs/zerolog/log"
)

type RequestBody struct{
	Text string `form:"text"`
	UserId string `form:"user_id"`
}

type CodingameInformation struct{
	Id int `db:"id"`
	MattermostUserId string `db:"mattermost_user_id"`
	CodingameUserId string `db:"codingame_user_id"`
	CgSession string `db:"cgSession"`
	RememberMe string `db:"rememberMe"`
}

func Login(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		var requestBody RequestBody
		if err := c.Bind(&requestBody); err != nil {
			log.Error().Err(err).Msg("Could not read request Body")
			c.JSON(400, "")
			// Handle error
		}

		//0 : codingame_user_id, 1: cgSession, 2: rememberMe
		userData := strings.Split(requestBody.Text, " ")
		if len(userData) != 3 {
			c.JSON(200, gin.H{
				"response_type": "in_channel",
				"text": fmt.Sprintf("Params are not good, %v, remember the parameters are Codingame user ID, cgSession, rememberMe, if you can't remember how to get those information, use /codingame_help", requestBody.Text),
			})
			return
		}

		user := CodingameInformation{}
		err := db.Get(&user, "SELECT * FROM codingame_informations WHERE mattermost_user_id = ? LIMIT 1", requestBody.UserId)
		if err == nil {
			_, err = db.Exec("UPDATE codingame_informations SET codingame_user_id = ?, cgSession = ?, rememberMe = ? WHERE mattermost_user_id = ?", userData[0], userData[1], userData[2], requestBody.UserId)
			c.JSON(200, gin.H{
				"response_type": "in_channel",
				"text": "You're data have been modified, you can launch a game with /codingame command !",
			})
			return
		}

		_, err = db.Exec("INSERT INTO codingame_informations (mattermost_user_id, codingame_user_id, cgSession, rememberMe) VALUES (?,?,?,?)", requestBody.UserId, userData[0], userData[1], userData[2] )

		if err != nil {
			log.Error().Err(err).Msg("Could not create information")
		}

		c.JSON(200, gin.H{
			"response_type": "in_channel",
			"text": "You're data have been created, you can launch a game with /codingame command !",
		})
	}
}