// Auth example is an example application which requires a login
// to view a private link. The username is "testuser" and the password
// is "password". This will require GORP and an SQLite3 database.
package main

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	//	"strconv"
)

type ViewRenderModel struct {
	Username   string
	PosterList []Poster
}

var dbmap *gorp.DbMap
var poster_dbmap *gorp.DbMap

func initDb() (*gorp.DbMap, *gorp.DbMap) {
	// Delete our SQLite database if it already exists so we have a clean start
	dbName := "nuposter.db"
	_, err := os.Open(dbName)
	if err == nil {
		os.Remove(dbName)
	}

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

func main() {
	store := sessions.NewCookieStore([]byte("secret123"))
	dbmap, poster_dbmap = initDb()

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
			return
		} else {
			r.Redirect("/signup")
			return
		}
	})

	m.Post("/login", binding.Bind(PosterUserModel{}), func(session sessions.Session, postedUser PosterUserModel, r render.Render, req *http.Request) {
		// You should verify credentials against a database or some other mechanism at this point.
		// Then you can authenticate this session.
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

	m.Post("/add_poster", binding.Bind(Poster{}), func(session sessions.Session, poster Poster, user sessionauth.User, r render.Render, req *http.Request) {
		// FIXME not authenticattion here
		fmt.Println("hehe => " + user.(*PosterUserModel).Username)
		fmt.Println("wocao => " + poster.Title)
		_ = InsertPoster(poster_dbmap, poster.Title, user.(*PosterUserModel).Username, poster.Date, poster.Location, poster.Tag, poster.Info, poster.Image)
		//poster.Author = user.(*PosterUserModel).Username
		//_ = InsertPoster(poster_dbmap, &poster)
		r.Redirect("/view_poster")
	})

	m.Post("/delete_poster", binding.Bind(Poster{}), func(session sessions.Session, poster Poster, user sessionauth.User, r render.Render, req *http.Request) {
		// FIXME not authenticattion here
		_ = DeletePoster(poster_dbmap, poster.Id)
		//		fmt.Println("id====" + strconv.FormatInt(poster.Id, 10))
		r.Redirect("/view_poster")
	})

	m.Run()
}
