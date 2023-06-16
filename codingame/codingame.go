package codingame

import (
	"net/http"
	"fmt"
	"strings"
	"io"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type CodinGameResponse struct {
	PublicHandle string `json:"publicHandle"`
}

func LaunchGame(db *sqlx.DB) gin.HandlerFunc{
	return func(c *gin.Context) {

		var requestBody RequestBody
		if err := c.Bind(&requestBody); err != nil {
			log.Error().Err(err).Msg("Could not read request Body")
			c.JSON(400, "")
			// Handle error
		}

		user := CodingameInformation{}
		err := db.Get(&user, "SELECT * FROM codingame_informations WHERE mattermost_user_id = ? LIMIT 1", requestBody.UserId)
		if err != nil{
			c.JSON(200, gin.H{
				"response_type": "in_channel",
				"text": "You're not connected, you should connect using /codingame_login",
			})
			return
		}

		url := "https://www.codingame.com/services/ClashOfCode/createPrivateClash"
		method := "POST"

		payload := strings.NewReader(fmt.Sprintf(`[%v,["Javascript"],["FASTEST","SHORTEST","REVERSE"]]`, user.CodingameUserId))

		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			log.Error().Err(err).Msg("Could not create request")
			c.JSON(200, gin.H{
				"response_type": "in_channel",
				"text": "Error while creating the request, our bad.",
			})
			return
		}
		
		req.Header.Add("content-type", "application/json;charset=UTF-8")
		req.Header.Add("cookie", fmt.Sprintf("cgSession=%v;rememberMe=%v;", user.CgSession, user.RememberMe))
	
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("Could not execute the request")
			c.JSON(200, gin.H{
				"response_type": "in_channel",
				"text": "Error while creating the request, our bad.",
			})
		}
		defer res.Body.Close()
		
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error().Err(err).Msg("Could not read response Body")
			return
		}

		fmt.Println(string(body))
	
		var response CodinGameResponse
		if err := json.Unmarshal(body, &response); err != nil {
			log.Error().Err(err).Msg("Could not read response Body json")
		}
	
		c.JSON(200, gin.H{
			"response_type": "in_channel",
			"text": fmt.Sprintf("https://www.codingame.com/clashofcode/clash/%v %v", response.PublicHandle, requestBody.Text),
		})
		
	}
}