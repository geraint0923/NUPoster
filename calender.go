package main

import (
	"fmt"
	//	"log"
	//	"net/http"
	//	"os"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"
	//	calendar "code.google.com/p/google-api-go-client/calendar/v3"
)

type timeS struct {
	dataTime string
}

type insertEvent struct {
	summary  string
	location string
	start    timeS
	end      timeS
}

const calendarID = "d2e8nb7tkmp21gbfl656vqh4j4@group.calendar.google.com"

func CreateEvent(title string, location string, startime int64, endtime int64) bool {
	client, key_api := InitAuth()

	sTime := timeS{time_Int2Str(startime)}

	eTime := timeS{time_Int2Str(endtime)}

	jEvent := insertEvent{title, location, sTime, eTime}

	event, err := json.Marshal(jEvent)
	if err != nil {
		panic(err)
		return false
	}

	resp, err := client.Post("https://www.googleapis.com/calendar/v3/calendars/"+calendarID+"/events?key="+key_api, "application/json", bytes.NewBuffer(event))
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		return false
	}

	return true
}

func subscribeEvent() {
}

func receiveEvent() {
}

func time_Int2Str(nsec int64) string {
	timestamp := time.Unix(0, nsec)

	timestampStr := timestamp.Format(time.RFC3339)

	str := strings.Split(timestampStr, "Z")
	fmt.Printf("*********%q\n", str)

	newStr := str[0] + ".000-" + str[1]

	return newStr
}
