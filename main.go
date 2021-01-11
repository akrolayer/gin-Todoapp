package main

import(
	"gin-todo/db"
	"gin-todo/router"
)

func main(){
	dbConn := db.Init()
	UserDB := db.UserInit()
	router.Router(dbConn, UserDB)
}