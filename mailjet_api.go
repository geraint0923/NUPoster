package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type MailJet struct {
	username string
	password string
}

type Mail struct {
	Locale      string
	Sender      string
	SenderEmail string
	Subject     string
	ListId      int
}

func InitMailJet(username string, password string) *MailJet {
	mail := new(MailJet)
	mail.username = username
	mail.password = password
	return mail
}

func parseJSON(data string) map[string]interface{} {
	jsonData := []byte(data)
	//	fmt.Println(data)
	u := map[string]interface{}{}
	err := json.Unmarshal(jsonData, &u)
	if err != nil {
		fmt.Println("heihei")
		panic(err)
	}
	return u
}

func (m *MailJet) BuildGroup(name string) int {
	var jsonStr = []byte("{\"name\":" + name + "}")
	req, _ := http.NewRequest("POST", "https://api.mailjet.com/v3/REST/contactslist", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//	fmt.Println(string(body))
	jsonObj := parseJSON(string(body))
	if jsonObj["Data"] == nil {
		return -1
	}
	retVal := jsonObj["Data"].([]interface{})[0].(map[string]interface{})["ID"].(float64)
	return int(retVal)
}

func (m *MailJet) DeleteGroup(groupId int) {
	req, _ := http.NewRequest("DELETE", "https://api.mailjet.com/v3/REST/contactslist/"+ strconv.Itoa(groupId), nil )
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}


func (m *MailJet) addContact(addr string) int {
	var jsonStr = []byte("{\"Email\":" + addr + "}")
	req, _ := http.NewRequest("POST", "https://api.mailjet.com/v3/REST/contact", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	jsonObj := parseJSON(string(body))
	if jsonObj["Data"] == nil {
		return -1
	}
	retVal := jsonObj["Data"].([]interface{})[0].(map[string]interface{})["ID"].(float64)
	return int(retVal)
}

func (m *MailJet) getContactID(addr string) int {
	req, _ := http.NewRequest("GET", "https://api.mailjet.com/v3/REST/contact/"+addr, nil)
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	jsonObj := parseJSON(string(body))
	if jsonObj["Data"] == nil {
		return -1
	}
	retVal := jsonObj["Data"].([]interface{})[0].(map[string]interface{})["ID"].(float64)
	return int(retVal)
}

func (m *MailJet) AddToGroup(listId int, addr string) {
	res := m.addContact(addr)
	fmt.Println("new created contact id = " + strconv.Itoa(res))
        fmt.Println(addr)
	if res < 0 {
		return
	}
	
	var jsonStr = []byte("{\"ContactID\":" + strconv.Itoa(res) + ",\"ListID\":" + strconv.Itoa(listId) + ",\"IsActive\": true}")
	req, _ := http.NewRequest("POST", "https://api.mailjet.com/v3/REST/listrecipient", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	/*
		jsonObj := parseJSON(string(body))
		if jsonObj["Data"] == nil {
			return
		}
	*/
}

func (m *MailJet) CreateNews(listID int) int{
        var jsonStr = []byte("{\"Locale\": \"en_us\", \"Sender\":\"NUPoster\", \"SenderEmail\": \"yanghu@u.northwestern.edu\",\"Subject\": \"Welcome to NU Poster\", \"ContactsListID\": " + strconv.Itoa(listID) + "}")
	req, _ := http.NewRequest("POST", "https://api.mailjet.com/v3/REST/newsletter", bytes.NewBuffer(jsonStr))
       	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

      	body, _ := ioutil.ReadAll(resp.Body)
//	fmt.Println(string(body))
	jsonObj := parseJSON(string(body))
	if jsonObj["Data"] == nil {
		return -1
	}
	retVal := jsonObj["Data"].([]interface{})[0].(map[string]interface{})["ID"].(float64)
	return int(retVal)
}



func (m *MailJet) SendToGroup(newsId int) {
        var jsonStr = []byte("{\"RefID\":" + strconv.Itoa(newsId) + ",\"JobType\":\"Send newsletter\", \"Status\": \"upload\"  }")
	fmt.Println("sending newsId " + strconv.Itoa(newsId))
	req, _ := http.NewRequest("POST", "https://api.mailjet.com/v3/REST/batchjob", bytes.NewBuffer(jsonStr))
       	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	//body, _ := ioutil.ReadAll(resp.Body)
//	fmt.Println(string(body))
	

}

func (m *MailJet) AddHtml(listID int, htmlContent string){
        var jsonStr = []byte(htmlContent)	
	req, _ := http.NewRequest("PUT", "https://api.mailjet.com/v3/DATA/Newsletter/"+ strconv.Itoa(listID)+"/HTML/text:html/LAST", bytes.NewBuffer(jsonStr))
       	req.Header.Set("Content-Type", "text/html")
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

}



func (m *MailJet) SendToUser(addr string, htmlContent string) {
        groupID := m.BuildGroup("temp")
        m.AddToGroup(groupID, addr)
	newsID := m.CreateNews(groupID)
       	fmt.Println("NewsID: " + strconv.Itoa(newsID))
        m.AddHtml(newsID, htmlContent)
        m.SendToGroup(newsID)
        m.DeleteGroup(groupID)
        fmt.Println("=====================")
}

/*
func main() {
	mail := InitMailJet("3bfc81a965fbc505da2950df6e2a5a4b", "6a26367d687abc1c4d543483d1870365")
//	mail.AddToGroup(1, "geraint0923@gmail.com")
       // res := mail.getContactID("yanghu2019@u.northwestern.edu")
//	fmt.Println("Result: " + strconv.Itoa(res))
        
       // newsID := mail.CreateNews(27)
       //	fmt.Println("NewsID: " + strconv.Itoa(newsID))
      //  mail.AddHtml(newsID, "<p>hello world</p> <img src=\"http://www.northwestern.edu/newscenter/images/toplevel2012/2014/11/tf1029.jpg\">")
       // mail.SendToGroup(newsID)
        
       // html_req, _ := http.NewRequest("GET", "http://104.236.41.69/contact.html", nil)
//	client := &http.Client{}
//	resp, _ := client.Do(html_req)
//	defer resp.Body.Close()
  //     	body, _ := ioutil.ReadAll(resp.Body)
    //    htmlContent := string(body)
      //  fmt.Println(htmlContent)
        mail.SendToUser("yanghu2019@gmail.com", "<p>hello! world! </p>") 
        //mail.SendToUser("yanghu2019@u.northwestern.edu", "<p>hello world</p> <img src=\"http://www.northwestern.edu/newscenter/images/toplevel2012/2014/11/tf1029.jpg\">")

}
*/
