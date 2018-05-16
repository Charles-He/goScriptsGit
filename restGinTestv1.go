package main

import (
	"net/http"
	"strconv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/buger/jsonparser"
	"gopkg.in/olahol/melody.v1"
	"encoding/json"
)

var db *gorm.DB
var m = melody.New()
var matchID string

func init() {
	//open a db connection
	var err error
	db, err = gorm.Open("mysql", "root:wddb_1q2w3e4r@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	//Migrate the schema
	db.AutoMigrate(&pgninfo{},&users{},&pgnmove{})
}

func main() {

	router := gin.Default()
	router.Use(cors.Default())
	// m := melody.New()

	v1 := router.Group("/api/v1/todos")
	{
		// API routers
		v1.POST("/addpgninfo", createPgnInfo)
		v1.POST("/addpgnmove", createPgnMove)
		v1.POST("/adduser", createUser)
		v1.GET("/getuser/:id", fetchSingleUser)
		v1.GET("/getmatches/:userid", fetchMatches)
		v1.GET("/getmoves/:matchid", fetchMoves)
		
		// chess live show page
		v1.GET("/matches", func(ctx *gin.Context) {
			http.ServeFile(ctx.Writer, ctx.Request, "/web/index.html")
		})
	
		v1.GET("/channel/:name/ws", func(ctx *gin.Context) {
			matchID = ctx.Param("name")
			m.HandleRequest(ctx.Writer, ctx.Request)
			
		})
		
		// v1.PUT("/:id", updateTodo)
		// v1.DELETE("/:id", deleteTodo)
	}
	
	m.HandleConnect(sendMove)
	
	router.Run()

}

type (
	users struct {
		Username string `gorm:"type:varchar(30); primary_key" json:"username"`
		Userid int `gorm:"type:int; not null" json:"userid" json:"userid"`
		Password string `gorm:"type:varchar(10); not null" json:"password"`
		Mailbox string `gorm:"type:varchar(30); not null" json:"mailbox"`
	}
	
	// todoModel describes a todoModel type
	pgninfo struct {
		// gorm.Model
		Userid int `gorm:"type:varchar(30); not null" json:"userid"`
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

// createTodo add a new todo
func createPgnInfo(c *gin.Context) {
	buf := make([]byte, 1024)  
    n, _ := c.Request.Body.Read(buf)  
    fmt.Println(string(buf[0:n]))  
	
    infouserid,_,_,_ := jsonparser.Get(buf[0:n], "userid")
	infouseridint,_ := strconv.Atoi(string(infouserid))
	infomatchid,_,_,_ := jsonparser.Get(buf[0:n], "matchid")
	infomatchidint,_ := strconv.Atoi(string(infomatchid))
	event,_,_,_ := jsonparser.Get(buf[0:n], "event")
	site,_,_,_ := jsonparser.Get(buf[0:n], "site")
	date,_,_,_ := jsonparser.Get(buf[0:n], "date")
	round,_,_,_ := jsonparser.Get(buf[0:n], "round")
	roundint,_ := strconv.Atoi(string(round))
	white,_,_,_ := jsonparser.Get(buf[0:n], "white")
	black,_,_,_ := jsonparser.Get(buf[0:n], "black")
	result,_,_,_ := jsonparser.Get(buf[0:n], "result")
	whitetype,_,_,_ := jsonparser.Get(buf[0:n], "whitetype")
	blacktype,_,_,_ := jsonparser.Get(buf[0:n], "blacktype")
	timecontrol,_,_,_ := jsonparser.Get(buf[0:n], "timecontrol")
	timecontrolint,_ := strconv.Atoi(string(timecontrol))
	rotation,_,_,_ := jsonparser.Get(buf[0:n], "rotation")
	rotationint,_ := strconv.Atoi(string(rotation))
	
	// save to database
	pgnInfo := pgninfo{
		Userid: infouseridint, Matchid: infomatchidint, Event: string(event), Site: string(site),
		Date: string(date), Round: roundint, White: string(white), Black: string(black),
		Result: string(result), WhiteType: string(whitetype), BlackType: string(blacktype), TimeControl: timecontrolint,
		Rotation: rotationint,
	}
	db.Save(&pgnInfo)
	c.JSON(http.StatusOK, string(event)) 
}
	
func createPgnMove(c *gin.Context) {	
	buf := make([]byte, 1024)  
    n, _ := c.Request.Body.Read(buf)  
    fmt.Println(string(buf[0:n])) 
	
	mvoematchid,_,_,_ := jsonparser.Get(buf[0:n], "matchid")
	mvoematchidint,_ := strconv.Atoi(string(mvoematchid))
	step,_,_,_ := jsonparser.Get(buf[0:n], "step")
	stepint,_ := strconv.Atoi(string(step))
	color,_,_,_ := jsonparser.Get(buf[0:n], "color")
	san,_,_,_ := jsonparser.Get(buf[0:n], "san")
		
	pgnMove := pgnmove{
		Matchid: mvoematchidint, Step: stepint, Color: string(color), San: string(san),
	}
	db.Save(&pgnMove)
		
	c.JSON(http.StatusOK, string(color))
}

func createUser(c *gin.Context) {	
	buf := make([]byte, 1024)  
    n, _ := c.Request.Body.Read(buf)  
    fmt.Println(string(buf[0:n])) 
	
	username,_,_,_ := jsonparser.Get(buf[0:n], "username")
	userid,_,_,_ := jsonparser.Get(buf[0:n], "userid")
	useridint,_ := strconv.Atoi(string(userid))
	password,_,_,_ := jsonparser.Get(buf[0:n], "password")
	mailbox,_,_,_ := jsonparser.Get(buf[0:n], "mailbox")
	
	Users := users{
		Username: string(username), Userid: useridint, Password: string(password), Mailbox: string(mailbox),
	}
	db.Save(&Users)
	
    c.JSON(http.StatusOK, string(username))  
}

// fetchSingleUser fetch a single user
func fetchSingleUser(c *gin.Context) {
	var user users
	userID := c.Param("id")
	db.First(&user, "userid = ?", userID)
	
	if user.Userid == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		return
	}
	
	_user := users{Username: user.Username, Userid: user.Userid, Password: user.Password, Mailbox: user.Mailbox}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _user})	
}

func fetchMatches(c *gin.Context) {
	var matchinfo []pgninfo
	userID := c.Param("userid")
	db.Where("userid = ?", userID).Find(&matchinfo)
	
	if matchinfo[0].Userid == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		return
	}
/*
	var _matchinfo []pgninfo
	for i, item := range matchinfo {
		_matchinfo := pgninfo{
			Userid: item[i].Userid, Matchid: item[i].Matchid, Event: item[i].Event, Site: item[i].Site,
			Date: item[i].Date, Round: item[i].Round, White: item[i].White, Black: matchinfo[i].Black,
			Result: matchinfo[i].Result, WhiteType: matchinfo[i].WhiteType, BlackType: matchinfo[i].BlackType, TimeControl: matchinfo[i].TimeControl,
			Rotation: matchinfo[i].Rotation,
		}
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": matchinfo})
	}
*/
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": matchinfo})
}

func fetchMoves(c *gin.Context) {
	var moveinfo []pgnmove
	matchID := c.Param("matchid")
	db.Where("matchid = ?", matchID).Find(&moveinfo)
	
	if moveinfo[0].Matchid == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No user found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": moveinfo})
}

func sendMove(s *melody.Session) {
	var moveinfo []pgnmove
	db.Where("matchid = ?", matchID).Find(&moveinfo)
	
	if moveinfo[0].Matchid == 0 {
		fmt.Println("no match matchid")
	}
	
	b, _ := json.Marshal(moveinfo)
	
	// m.Broadcast(msg)
	sArray := []*melody.Session{s}
	m.BroadcastMultiple(b, sArray)
}
