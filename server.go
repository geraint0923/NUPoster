package main

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	//	"strconv"
)

const EVENT_NUM = 2

type ViewRenderModel struct {
	Username   string
	PosterList []Poster
}

type PosterItem struct {
	Title       string                `form:"title"`
	Start       string                `form:"start"`
	End         string                `form:"end"`
	Location    string                `form:"location"`
	Talk        string                `form:"talk"`
	Show        string                `form:"show"`
	Athletics   string                `form:"althletics"`
	Recruit     string                `form:"recruit"`
	LostFound   string                `form:"lost_found"`
	FreeFood    string                `form:"free_food"`
	Info        string                `form:"info"`
	ImageUpload *multipart.FileHeader `form:"image"`
	Image       string
}

type SubscribeModel struct {
	Email     string `form:"email"`
	Talk      string `form:"talk"`
	Show      string `form:"show"`
	Athletics string `form:"althletics"`
	Recruit   string `form:"recruit"`
	LostFound string `form:"lost_found"`
	FreeFood  string `form:"free_food"`
}

// talks shows athletics recruits lost&found free_food
var tagMailList map[string]int
var dbmap *gorp.DbMap
var poster_dbmap *gorp.DbMap

var welcomeEmail string

func initDb() (*gorp.DbMap, *gorp.DbMap) {
	// Delete our SQLite database if it already exists so we have a clean start
	dbName := "nuposter.db"
	/*
		_, err := os.Open(dbName)
		if err == nil {
			os.Remove(dbName)
		}
	*/

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatalln("Fail to create database", err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	poster_dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(PosterUserModel{}, "users").SetKeys(true, "Id")
	poster_dbmap.AddTableWithName(Poster{}, "posters").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln("Could not build tables", err)
	}
	err = poster_dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln("Could not build tables", err)
	}

	/*
		user := PosterUserModel{1, "testuser", "password", false}
			err = dbmap.Insert(&user)
			if err != nil {
				log.Fatalln("Could not insert test user", err)
			}
	*/
	insertUser(dbmap, "testuser", "password")
	return dbmap, poster_dbmap
}

func insertUser(dbmap *gorp.DbMap, username string, passwd string) {
	user := PosterUserModel{1, username, passwd, false}
	err := dbmap.Insert(&user)
	if err != nil {
		log.Fatalln("failed to signup new user", err)
	}
}

var mailJet *MailJet

func procSubscribe(sub SubscribeModel) {
	if sub.Talk != "" {
		mailJet.AddToGroup(tagMailList["talk"], sub.Email)
	}
	if sub.Show != "" {
		mailJet.AddToGroup(tagMailList["show"], sub.Email)
	}
	if sub.Athletics != "" {
		mailJet.AddToGroup(tagMailList["athletics"], sub.Email)
	}
	if sub.Recruit != "" {
		mailJet.AddToGroup(tagMailList["recruit"], sub.Email)
	}
	if sub.LostFound != "" {
		mailJet.AddToGroup(tagMailList["lost_found"], sub.Email)
	}
	if sub.FreeFood != "" {
		mailJet.AddToGroup(tagMailList["free_food"], sub.Email)
	}
	fmt.Println("Email =>" + sub.Email)
}

func procSignup(username string) {
	msg := fmt.Sprintf(welcomeEmail, username)
	mailJet.SendToUser(username, msg)
	fmt.Println("=88=8=8=8=8=8= => finish procSignup   " + username + "         " + msg)
}

func main() {
	store := sessions.NewCookieStore([]byte("secret123"))
	dbmap, poster_dbmap = initDb()

	mailJet = InitMailJet("3bfc81a965fbc505da2950df6e2a5a4b", "6a26367d687abc1c4d543483d1870365")

	tagMailList = make(map[string]int)

	tagMailList["talk"] = 22
	tagMailList["sport"] = 23
	tagMailList["show"] = 24
	tagMailList["recruit"] = 25
	tagMailList["lost_found"] = 26
	tagMailList["free_food"] = 27

	welcomeEmail = "<html> <div align=\"middle\"><img src=\"http://yyj.cloudapp.com/welcome.jpg\"/><br/>Hi %s,<br/>Welcome to NUPoster!<br/> You can begin to post your own poster on NUPoster! Enjoy yourself!</div></html>"

	m := martini.Classic()
	m.Use(render.Renderer())

	// Default our store to use Session cookies, so we don't leave logged in
	// users roaming around
	store.Options(sessions.Options{
		MaxAge: 0,
	})
	m.Use(sessions.Sessions("sessionid", store))
	m.Use(sessionauth.SessionUser(GenerateAnonymousUser))
	sessionauth.RedirectUrl = "/login"
	sessionauth.RedirectParam = "redirect_url"

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})

	m.Get("/signup", func(r render.Render) {
		r.HTML(200, "signup", nil)
	})

	m.Post("/signup", binding.Bind(PosterUserModel{}), func(session sessions.Session, postedUser PosterUserModel, r render.Render, req *http.Request) {
		// You should verify credentials against a database or some other mechanism at this point.
		// Then you can authenticate this session.
		user := PosterUserModel{}
		err := dbmap.SelectOne(&user, "SELECT * FROM users WHERE username = $1", postedUser.Username)
		if err != nil {
			insertUser(dbmap, postedUser.Username, postedUser.Password)
			_ = dbmap.SelectOne(&user, "SELECT * FROM users WHERE username = $1", postedUser.Username)
			err = sessionauth.AuthenticateSession(session, &user)
			if err != nil {
				r.JSON(500, err)
			}
			r.Redirect("/")
			go procSignup(postedUser.Username)
			//str := fmt.Sprintf(welcomeEmail, postedUser.Username)
			//fmt.Println(str)
			return
		} else {
			r.Redirect("/signup")
			return
		}
	})

	m.Post("/login", binding.Bind(PosterUserModel{}), func(session sessions.Session, postedUser PosterUserModel, r render.Render, req *http.Request) {
		user := PosterUserModel{}
		err := dbmap.SelectOne(&user, "SELECT * FROM users WHERE username = $1 and password = $2", postedUser.Username, postedUser.Password)
		if err != nil {
			r.Redirect(sessionauth.RedirectUrl)
			return
		} else {
			err := sessionauth.AuthenticateSession(session, &user)
			if err != nil {
				r.JSON(500, err)
			}
			//params := req.URL.Query()
			//redirect := params.Get(sessionauth.RedirectParam)
			//r.Redirect(redirect)
			r.Redirect("/")
			return
		}
	})

	m.Get("/logout", sessionauth.LoginRequired, func(session sessions.Session, user sessionauth.User, r render.Render) {
		sessionauth.Logout(session, user)
		r.Redirect("/")
	})

	m.Get("/add_poster", sessionauth.LoginRequired, func(r render.Render, user sessionauth.User) {
		//r.HTML(200, "add_poster", user.(*PosterUserModel))
		r.HTML(200, "add_poster", user.(*PosterUserModel))
	})

	m.Get("/view_poster", sessionauth.LoginRequired, func(r render.Render, user sessionauth.User) {
		data := ViewRenderModel{Username: user.(*PosterUserModel).Username}
		poster_dbmap.Select(&data.PosterList, "select * from posters where author=\""+data.Username+"\"")
		fmt.Println("namename => " + data.Username)
		for _, val := range data.PosterList {
			fmt.Println("ent =>" + val.Title)
		}
		r.HTML(200, "view_poster", data)
	})

	m.Post("/add_poster", binding.Bind(PosterItem{}), func(session sessions.Session, poster PosterItem, user sessionauth.User, r render.Render, req *http.Request) {
		// FIXME not authenticattion here
		fmt.Println("hehe => " + user.(*PosterUserModel).Username)
		fmt.Println("wocao => " + poster.Title)
		imgPath := uuid.New()
		fmt.Println(imgPath)
		if poster.ImageUpload != nil {
			file, _ := poster.ImageUpload.Open()
			defer file.Close()
			outputFile, _ := os.Create("public/img/" + imgPath)
			defer outputFile.Close()
			_, _ = io.Copy(outputFile, file)
		}
		tagString := ""
		hostUrl := "http://yyj.cloudapp.net"
		poster.Image = hostUrl + "/img/" + imgPath

		tmpl, _ := template.ParseFiles("templates/newsletter.tmpl")
		buf := new(bytes.Buffer)
		_ = tmpl.Execute(buf, poster)
		tmplString := buf.String()
		fmt.Println(tmplString)

		if poster.Talk != "" {
			tagString += "talk "
			newsID := mailJet.CreateNews(tagMailList["talk"])
			mailJet.AddHtml(newsID, tmplString)
			mailJet.SendToGroup(newsID)
		}
		if poster.Show != "" {
			tagString += "show "
			newsID := mailJet.CreateNews(tagMailList["show"])
			mailJet.AddHtml(newsID, tmplString)
			mailJet.SendToGroup(newsID)
		}
		if poster.Athletics != "" {
			tagString += "athletics "

			newsID := mailJet.CreateNews(tagMailList["athletics"])
			mailJet.AddHtml(newsID, tmplString)
			mailJet.SendToGroup(newsID)

		}
		if poster.Recruit != "" {
			tagString += "recruit "

			newsID := mailJet.CreateNews(tagMailList["recruit"])
			mailJet.AddHtml(newsID, tmplString)
			mailJet.SendToGroup(newsID)
		}
		if poster.LostFound != "" {
			tagString += "lost_found "

			newsID := mailJet.CreateNews(tagMailList["lost_found"])
			mailJet.AddHtml(newsID, tmplString)
			mailJet.SendToGroup(newsID)
		}
		if poster.FreeFood != "" {
			tagString += "free_food"

			newsID := mailJet.CreateNews(tagMailList["free_food"])
			mailJet.AddHtml(newsID, tmplString)
			mailJet.SendToGroup(newsID)
		}
		fmt.Println(poster.Start + " ====" + poster.End)
		startTime := TransformTime(poster.Start)
		endTime := TransformTime(poster.End)
		logStr := CreateEvent(poster.Title, poster.Location, startTime, endTime)
		_ = InsertPoster(poster_dbmap, poster.Title, user.(*PosterUserModel).Username, startTime, endTime, poster.Location, tagString, poster.Info, "/img/"+imgPath, logStr)

		fmt.Println("================>" + logStr)
		//poster.Author = user.(*PosterUserModel).Username
		//_ = InsertPoster(poster_dbmap, &poster)
		//r.Redirect("/view_poster")
		r.Redirect("/")
	})

	m.Post("/delete_poster", binding.Bind(Poster{}), func(session sessions.Session, poster Poster, user sessionauth.User, r render.Render, req *http.Request) {
		// FIXME not authenticattion here
		curPoster := &Poster{}
		_ = dbmap.SelectOne(&curPoster, "SELECT * FROM posters WHERE id=$1", poster.Id)
		//		fmt.Println("name=>" + curPoster.Image)
		realPath := "public" + curPoster.Image
		if _, inErr := os.Stat(realPath); inErr == nil {
			_ = os.Remove(realPath)
		}
		_ = DeletePoster(poster_dbmap, poster.Id)
		//		fmt.Println("id====" + strconv.FormatInt(poster.Id, 10))
		r.Redirect("/view_poster")
	})

	// to retrieve poster information
	m.Get("/posters", func(r render.Render, req *http.Request) {
		params := req.URL.Query()
		tag := params.Get("tag")
		//		page := params.Get("page")
		var poster []Poster
		poster_dbmap.Select(&poster, "select * from posters where tag like \"%"+tag+"%\" order by id desc")
		//	fmt.Println(tag + "=>>>>>>>>>" + page)
		r.JSON(200, poster)
	})

	m.Get("/search", func(r render.Render, req *http.Request) {
		params := req.URL.Query()
		tag := params.Get("q")
		//		page := params.Get("page")
		var poster []Poster
		poster_dbmap.Select(&poster, "select * from posters where title like \"%"+tag+"%\" order by id desc")
		//	fmt.Println(tag + "=>>>>>>>>>" + page)
		r.JSON(200, poster)
	})

	m.Get("/sub", func(r render.Render) {
		//r.HTML(200, "add_poster", user.(*PosterUserModel))
		r.HTML(200, "sub", nil)
	})
	m.Post("/subscribe", binding.Bind(SubscribeModel{}), func(sub SubscribeModel, r render.Render, req *http.Request) {
		go procSubscribe(sub)
		r.Redirect("/")
	})

	m.Run()
}
