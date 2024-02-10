package worker

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/s29papi/wag3r-bot/worker/client"
	"github.com/s29papi/wag3r-bot/worker/types"
)

// check if a new cast has been added
func (w *Worker) process(d []byte) {
	var userMentions types.UserMentions
	if err := json.Unmarshal(d, &userMentions); err != nil {
		log.Println(err)
		// handle a return here
	}
	for _, notifs := range userMentions.Notifications {
		t := timestamp2secs(notifs.Cast.Timestamp)
		if t == 0 || t <= *w.lastProcReqTime {
			continue
		}
		if notifs.Type != "mention" {
			continue
		}

		// parse text
		var info = &types.Game{}
		errNo := process(notifs.Cast.Text, info)
		switch errNo {
		case ERR_MAX_NUMBER_LINES_NO:
		}

		// handle wrong input

		// other parts of today
		// make a request to vercel

		// vercel store in db

		// returns url

		w.lastProcReqTime = &t
	}

	w.pauseFn <- struct{}{}
}

// process takes a cast's text and returns either nil if successful and error
// if it fails.
func process(text string, info *types.Game) int64 {
	reg := regexp.MustCompile(`^\s*\*\s*(.+?):\s*(.*)$`)
	mentionReg := regexp.MustCompile(`@\w+`)
	textLines := strings.Split(text, "\n")

	if len(textLines) >= MAXNUMBERLINES {
		return int64(ERR_MAX_NUMBER_LINES_NO)
	}

	var challengeStarted bool

	for _, textLine := range textLines {
		textLine = strings.TrimSpace(textLine)

		if textLine == "" {
			continue
		}

		if !challengeStarted {
			if textLine == "Open /stadium Challenge:" {
				challengeStarted = true
			}
			continue
		}

		textLine = mentionReg.ReplaceAllString(textLine, "")

		if strings.HasPrefix(textLine, "*") {
			match := reg.FindStringSubmatch(textLine)
			if len(match) == 3 {
				key := strings.Title(match[1])
				value := match[2]

				switch key {
				case "Game":
					info.Name = value
				case "Match Set Up":
					info.Setup = value
				case "Match Date":
					info.Date = value
				case "Wager Amount":
					wagerAmountStr := strings.TrimSpace(value)
					if !strings.HasPrefix(wagerAmountStr, "$") {
						return ERR_MISSING_CURRENCY_SYMBOL_NO
					}

					wagerAmountParts := strings.Fields(value)
					if len(wagerAmountParts) != 3 {
						return ERR_INVALID_LENGTH_WAGERAMOUNT_NO
					}

					amountStr := wagerAmountParts[1]
					token := wagerAmountParts[2]
					amount, err := strconv.ParseFloat(amountStr, 64)
					if err != nil {
						return ERR_INVALID_AMOUNT_NO
					}
					info.Amount = amount
					info.Token = token
				default:
					return ERR_UNEXPECTED_FIELD_NO
				}
			}
		}
		if info.Name == "" || info.Setup == "" || info.Date == "" || info.Amount == 0 || info.Token == "" {
			return ERR_MISSING_REQ_FIELD_NO
		}

	}
	return 0 // sucess
}

func checkNewCast(service *client.HTTPService) {
	req := client.ChannelCastRequest()
	go service.SendRequest(http.MethodGet, req)
}

func checkNewMentions(service *client.HTTPService) {
	req := client.UserMentionsRequest()
	go service.SendRequest(http.MethodGet, req)
}

func timestamp2secs(t string) int64 {
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		log.Println("Error parsing timestamp:", err)
		return 0
	}
	return parsedTime.Unix()
}