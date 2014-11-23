package main

import (
	"github.com/coopernurse/gorp"
)

type Poster struct {
	Id       int64  `form:"id" db:"id" json:"id"`
	Title    string `form:"title" db:"title" json:"title"`
	Author   string `form:"Author", db:"author" json:"author"`
	Start    int64  `form:"start" db:"start" json:"start"`
	End      int64  `form:"end" db:"end" json:"end"`
	Location string `form:"location" db:"location" json:"location"`
	Tag      string `form:"tag" db:"tag" json:"tag"`
	Info     string `form:"info" db:"info" json:"info"`
	Image    string `form:"image" db:"image" json:"image"`
	EventID  string `form:"event" db:"event" json:"event"`
}

func InsertPoster(dbmap *gorp.DbMap, title string, author string, start int64, end int64, loc string, tag string, info string, img string, e string) error {
	poster := Poster{Title: title, Author: author, Start: start, End: end, Location: loc, Tag: tag, Info: info, Image: img, EventID: e}
	return dbmap.Insert(&poster)
}

/*
func InsertPoster(dmmap *gorp.DbMap, poster *Poster) error {
	return dbmap.Insert(&poster)
}
*/

func DeletePoster(dbmap *gorp.DbMap, id int64) error {
	_, err := dbmap.Delete(&Poster{Id: id})
	return err
}
