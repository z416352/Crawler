package utils

import (
	"Crawler/pkg/logger"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)
const format = "2006-01-02 15:04:05"

func Send_message_to_bot(msg string) bool {
	data := make(map[string]string)
	data["Message"] = msg
	data["Chat_ID"] = "998618031"
	b, _ := json.Marshal(data)

	resp, err := http.Post("http://localhost:5000/Send_Message",
		"application/json",
		bytes.NewBuffer(b),
	)

	if err != nil {
		logger.UtilsLog.Error(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logger.UtilsLog.Debugf("Send_message_to_bot : %s", string(body))

	return true
}

// Find the latest candlestick time.
func NewestKlineTime(mins int) int64 {
	var hours, days int

	if mins >= 60 {
		hours = mins / 60
		mins = mins % 60
	}
	if hours >= 24 {
		days = hours / 24
		hours = hours % 24
	}
	logger.UtilsLog.Tracef("Time Frame: %d days, %d hrs, %d mins", days, hours, mins)

	currTime := time.Now()

	if mins != 0 {
		currTime = currTime.Add(-time.Minute * time.Duration(currTime.Minute()%mins))
	} else if mins == 0 && hours != 0 {
		currTime = currTime.Add(-time.Minute * time.Duration(currTime.Minute()))
	}

	if hours != 0 {
		currTime = currTime.Add(-time.Hour * time.Duration(currTime.Hour()%hours))
	} else if hours == 0 && days != 0 {
		currTime = currTime.Add(-time.Hour * time.Duration(currTime.Hour()))
	}

	if days != 0 {
		currTime = currTime.Add(-time.Minute * time.Duration(currTime.Minute()))
		currTime = currTime.Add(-time.Hour * time.Duration(currTime.Hour()))

		currTime = currTime.AddDate(0, 0, -currTime.Day()%days) // AddDate(years, months, days)
	}

	currTime = currTime.Add(-time.Second * time.Duration(currTime.Second()))

	unix := time.Date(
		currTime.Year(),
		currTime.Month(),
		currTime.Day(),
		currTime.Hour(),
		currTime.Minute(),
		currTime.Second(),
		0,
		time.UTC,
	).UnixMilli()

	logger.UtilsLog.Debugf("The newest K lines time: %v", time.UnixMilli(unix).Format(format))

	return unix
}


// Shift 'n' intervals timeframe
func ShiftTimeframe(mins int, start_unix int64, n int) int64 {
	date := time.UnixMilli(start_unix)
	logger.UtilsLog.Tracef("---------------------  Shift Timeframe ---------------------")
	logger.UtilsLog.Tracef("Time frame: %d mins, n: %d --> Shift %d mins", mins, n, mins*n)
	logger.UtilsLog.Tracef("Before Shift : %v", date.Format(format))

	// shift time
	date = date.Add(-time.Minute * time.Duration(mins*n))

	logger.UtilsLog.Tracef("After Shift : %v", date.Format(format))

	return date.UnixMilli()
}
