package main

import (
	"flag"
	"net/http"

	"golang/handlers"
	"golang/models"
	"golang/tools"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configFile := flag.String("config", "config.yml", "конфиг сервера")
	sync := flag.Bool("sync", false, "синкануть скрипты")
	flag.Parse()
	tools.SetupConfig(*configFile)
	//sshio := make(chan []byte)

	e := echo.New()
	e.Debug = true
	db := tools.InitDB(tools.Config.Db.ProdPath)
	tools.Migrate(db)

	if *sync {
		sshConfig := tools.SSHConfig{
			Host:     tools.Config.SSH.Host,
			Port:     tools.Config.SSH.Port,
			Login:    tools.Config.SSH.Login,
			Password: tools.Config.SSH.Password,
			Keypath:  tools.Config.SSH.Keypath,
		}

		scripts, err := tools.GetAnsibleScripts(&sshConfig)
		tools.CheckErr(err)
		updates, err := tools.StoreScriptsToDB(&scripts, db)
		e.Logger.Printf("Rows Updated/Inserted: %d\n", updates)
		return

	}

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		user, err := models.GetUser(username, db)
		if err != nil {
			return false, echo.NewHTTPError(http.StatusUnauthorized, "Invalid Creds")
		}
		if user.Password != password {
			return false, echo.NewHTTPError(http.StatusUnauthorized, "Invalid Password")
		}
		c.Set("user", username)
		return true, nil
	}))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("2ri0LwhR62EVYlh16rWQJz5AHVA8npsyL"))))
	e.File("/", "public/index.html")
	e.Static("/static", "public/static")
	e.GET("/ansiblescripts", handlers.GetAnsibleScripts(db))
	e.PUT("/ansiblescripts", handlers.RunAnsibleScriptsAsync(db))
	e.GET("/ansiblescript_out/:uuid", handlers.GetScriptOut())
	e.GET("/logoff", handlers.Logoff())
	e.Logger.Fatal(e.Start(tools.Config.Network.Listen))
}
