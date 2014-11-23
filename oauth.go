package main

import (
	//	"bytes"
	//	"encoding/json"
	//	"fmt"
	"github.com/golang/oauth2"
	"github.com/golang/oauth2/google"
	//	"io/ioutil"
	"log"
	"net/http"
)

//const clientID = "429009888040-hntt6h7j1pphiavktik6cc81cu33kknq.apps.googleusercontent.com"
const clientID = "429009888040-hntt6h7j1pphiavktik6cc81cu33kknq@developer.gserviceaccount.com"
const key_api = "AIzaSyA7OF9QLo35MVkzUw2Szizpe6HiGquOxtk"

func InitAuth() (http.Client, string) {
	f, err := oauth2.New(
		google.ServiceAccountJSONKey("./nuposter-19a3aab776f0.json"),
		oauth2.Scope(
			"https://www.googleapis.com/auth/calendar",
			//			"https://www.googleapis.com/auth/blogger",
		),
	)

	if err != nil {
		log.Fatal(err)
	}
	// Initiate an http.Client. The following GET request will be
	// authorized and authenticated on the behalf of
	// your service account.
	client := http.Client{Transport: f.NewTransport()}
	return client, key_api
}

/*
func main() {
	client := InitAuth()

	event := []byte(`
	{
		"end":
		{
			"dateTime": "2011-06-03T10:25:00.000-07:00"
		},
		"start":
		{
			"dateTime": "2011-06-03T10:00:00.000-07:00"
		}
	}
	`)

	//	resq, err := client.Post("https://www.googleapis.com/calendar/v3/calendars/3aigle8maruau3tt9p5dnsme68@group.calendar.google.com/events", "application/json", bytes.NewBuffer(event))
	resq, err := client.Post("https://www.googleapis.com/calendar/v3/calendars/d2e8nb7tkmp21gbfl656vqh4j4@group.calendar.google.com/events?key=AIzaSyA7OF9QLo35MVkzUw2Szizpe6HiGquOxtk", "application/json", bytes.NewBuffer(event))

	defer resq.Body.Close()
	fmt.Println("response Status:", resq.Status)
	fmt.Println("response Headers:", resq.Header)
	body, _ := ioutil.ReadAll(resq.Body)
	fmt.Println("response Body:", string(body))

	resq, err = client.Get("https://www.googleapis.com/calendar/v3/calendars/r8mbv5n092d6e9ous5macum008@group.calendar.google.com/events?key=AIzaSyA7OF9QLo35MVkzUw2Szizpe6HiGquOxtk")

	defer resq.Body.Close()
	fmt.Println("response Status:", resq.Status)
	fmt.Println("response Headers:", resq.Header)
	body, _ = ioutil.ReadAll(resq.Body)
	fmt.Println("response Body:", string(body))

	//fmt.Println("Response: ", body)
	//fmt.Println("Response: ", body)

	fmt.Println("Error: ", err)

	r, err := http.Post("https://www.googleapis.com/calendar/v3/calendars/3aigle8maruau3tt9p5dnsme68@group.calendar.google.com/events?fields=klaskdasd&key=AIzaSyD8Rm0BfwMk5dnzQ3oPESOgS8FRVXKKY7Q", "application/json", bytes.NewBuffer(event))
	fmt.Println("Response: ", r)
	fmt.Println("Error: ", err)

}
*/
