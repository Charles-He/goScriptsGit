package main

import (
    "fmt"
    "math/rand"
    "log"
    "net/http"
    "strings"
    "bytes"
    "encoding/json"
	"os"
	"bufio"
	"regexp"
	"strconv"
	"io/ioutil"
	"time"
)

/*
type (
	users struct {
		Username string `gorm:"type:varchar(30); primary_key" json:"username"`
		Userid int `gorm:"type:int; not null" json:"userid"`
		Password string `gorm:"type:varchar(10); not null" json:"password"`
		Mailbox string `gorm:"type:varchar(30); not null" json:"mailbox"`
	}
	
	// todoModel describes a todoModel type
	pgninfo struct {
		// gorm.Model
		Userid int `gorm:"type:varchar(30); not null json:"userid"`
		Matchid	int `gorm:"type:int; not null; primary_key" json:"matchid"`
		Event string    `gorm:"type:varchar(50); not null" json:"event"`
		Site string `gorm:"type:varchar(50); not null" json:"site"`
		Date string `gorm:"type:date" json:"date"`
		Round int `gorm:"type:int" json:"round"`
		White string `gorm:"type:varchar(50); not null" json:"white"`
		Black string `gorm:"type:varchar(50); not null" json:"black"`
		Result string `gorm:"type:varchar(10); not null" json:"result"`
		WhiteType string `gorm:"type:varchar(30); not null" json:"whitetype"`
		BlackType string `gorm:"type:varchar(30); not null" json:"blacktype"`
		TimeControl int `gorm:"type:int" json:"timecontrol"`
		Rotation int `gorm:"type:int" json:"rotation"`
	}

	// transformedTodo represents a formatted todo
	pgnmove struct {
		Matchid int `gorm:"type:int; not null; primary_key" json:"matchid"`
		Step int `gorm:"type:int; not null; primary_key" json:"step"` 
		Color string `gorm:"type:varchar(10); not null; primary_key" json:"color"`
		San string `gorm:"type:varchar(10); not null" json:"san"`
	}
)
*/

func main() {
	userdata := make(map[string]string)
	infodata := make(map[string]string)
	movedata := make(map[string]string)
	
	userUrl := "http://192.168.1.2:8080/api/v1/todos/adduser"
	moveUrl := "http://192.168.1.2:8080/api/v1/todos/addpgnmove"
	infoUrl := "http://192.168.1.2:8080/api/v1/todos/addpgninfo"

	// assign user values
	userdata["username"] = "Charles"
	userdata["userid"] = "1"
	userdata["password"] = "1234a"
	userdata["mailbox"] = "Charles@hotmail.com"
	
	// post userdata
	senddata (userdata, userUrl)

	// initial some values
	marker := 0                            // use to determain when to post infodata
	mid := rand.Intn(10000)                // will be changed for every match
	infodata["userid"] = "1"                 // fot this test there is only one user
	
	
	file, err := os.Open("testpgn.pgn")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
	
	scanner := bufio.NewScanner(file)
	
	var stepbuffer bytes.Buffer
	
    for scanner.Scan() {
		s := scanner.Text()
		// fmt.Println(s)
		
		if strings.Contains(s, "[") {
			re, _ := regexp.Compile(`(\w+) "(.*)"`) 
			result_slice := re.FindAllStringSubmatch(s, -1)
			if result_slice[0][1] != "" {
				// fmt.Println(result_slice[0][1], result_slice[0][2])
		
				infodata[strings.ToLower(result_slice[0][1])] = result_slice[0][2]
			
				// fmt.Printf("%v", result_slice)
				// break
			}
		} else if strings.Contains(s, ".") {
			// post infodata first
			if marker == 0 {
				infodata["matchid"] = strconv.Itoa(mid)
				infodata["rotation"] = "1"
				senddata (infodata, infoUrl)
				marker = 1
			}
			
			stepbuffer.WriteString(s)
			stepbuffer.WriteString(" ")
			
			re, _ := regexp.Compile(`\d-`)
			match := re.MatchString(s)
			if match {
				stepstring := stepbuffer.String()
				stepbuffer.Reset()
				marker = 0
				
				re, _ := regexp.Compile(`(\d+)\.(\S+) ((\S+)?)`) 
				result_slice := re.FindAllStringSubmatch(stepstring, -1)
				
				// assign valuse to movedata
				// mid := rand.Intn(10000)
				for _, value := range result_slice {
					movedata["matchid"] = strconv.Itoa(mid)
					movedata["step"] = value[1]
					movedata["color"] = "white"
					movedata["san"] = value[2]
					senddata (movedata, moveUrl)
					
					if value[3] != "" {
						movedata["matchid"] = strconv.Itoa(mid)
						movedata["step"] = value[1]
						movedata["color"] = "black"
						movedata["san"] = value[3]
						senddata (movedata, moveUrl)
					}
					
					// fmt.Printf("%v", result_slice)
					// break
				}
				
				mid = rand.Intn(10000)
				
			}

		}
		
    }
}

func senddata(a map[string]string, url string) {
	bs, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
    fmt.Println(string(bs))
    

	// var jsonStr = []byte{bs}
	// url := "http://192.168.1.2:8080/api/v1/todos/"
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))

	time.Sleep(1 * time.Second)
}