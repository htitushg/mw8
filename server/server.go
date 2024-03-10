package server

import (
	"html/template"
	"log"
	database "mw8/database"
	"mw8/internal/models"
	"mw8/internal/utils"
	"mw8/router"
	"net/http"
	"os"
)

func Run() {
	// Initializing database
	/* models.Conf_Mysql.Host = "127.0.0.1"
	models.Conf_Mysql.Port = "3306"
	models.Conf_Mysql.Database = "basemw8"
	models.Conf_Mysql.User = "henry"
	models.Conf_Mysql.Password = "11nhri04p" */

	models.Conf_Mysql.Host = os.Getenv("MYSQL_HOST")
	models.Conf_Mysql.Port = os.Getenv("MYSQL_PORT")
	models.Conf_Mysql.Database = os.Getenv("MYSQL_DATABASE")
	models.Conf_Mysql.User = os.Getenv("MYSQL_USER")
	models.Conf_Mysql.Password = os.Getenv("MYSQL_PASSWORD")
	database.InitDB(models.Conf_Mysql.Database)
	// Initialize templates
	tmplPath := utils.Path + "templates/"
	models.Tmpl["index"] = template.Must(template.ParseFiles(tmplPath+"index2.gohtml", tmplPath+"layouts/base.layout.gohtml"))
	models.Tmpl["login"] = template.Must(template.ParseFiles(tmplPath+"login2.gohtml", tmplPath+"layouts/base.layout.gohtml"))
	models.Tmpl["modifuser"] = template.Must(template.ParseFiles(tmplPath+"modifuser2.gohtml", tmplPath+"layouts/base.layout.gohtml"))
	models.Tmpl["register"] = template.Must(template.ParseFiles(tmplPath+"register2.gohtml", tmplPath+"layouts/base.layout.gohtml"))
	models.Tmpl["errormethod"] = template.Must(template.ParseFiles(tmplPath+"errormethod2.gohtml", tmplPath+"layouts/base.layout.gohtml"))
	models.Tmpl["error404"] = template.Must(template.ParseFiles(tmplPath+"error4042.gohtml", tmplPath+"layouts/base.layout.gohtml"))

	// Initializing the routes
	router.Init()

	// Sending the assets to the clients
	fs := http.FileServer(http.Dir(utils.Path + "assets"))
	router.Mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Running the goroutine to automatically remove expired sessions every given time
	go utils.MonitorSessions()

	// Running the goroutine to change log file every given time
	go utils.LogInit()

	// Running the goroutine to automatically remove old TempUsers
	go utils.ManageTempUsers()

	// Running the goroutine to automatically remove old ManageTempModifUsers
	go utils.ManageTempModifUsers()

	// Running the server
	log.Fatalln(http.ListenAndServe(models.Port, router.Mux))
}
