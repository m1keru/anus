package handlers

import (
	"database/sql"
	"janesface/models"
	"janesface/tools"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo"
)

type H map[string]interface{}

func GetAnsibleScripts(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, func() models.AnsibleScriptCollection {
			var result models.AnsibleScriptCollection
			var err error
			if c.Get("user") == "osa" {
				result, err = models.GetAnsibleScriptsFiltered(db, "case")
			} else {
				result, err = models.GetAnsibleScripts(db)
			}
			if err != nil {
				c.Logger().Error(err)
				return result
			}

			sess, _ := session.Get("session", c)
			sess.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   86400 * 7,
				HttpOnly: true,
			}
			sess.Save(c.Request(), c.Response())
			return result
		}())
	}
}

func RunAnsibleScripts(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		sess, _ := session.Get("session", c)
		busy, _ := sess.Values["busy"]
		if busy == "true" {
			c.Logger().Debug("Job Already Running")
		}

		// создаём новую задачу
		var ansibleScript models.AnsibleScript
		// привязываем пришедший JSON в новую задачу
		c.Bind(&ansibleScript)
		// добавим задачу с помощью модели
		sshConfig := tools.SSHConfig{
			Host:     tools.Config.SSH.Host,
			Port:     tools.Config.SSH.Port,
			Login:    tools.Config.SSH.Login,
			Password: tools.Config.SSH.Password,
			Keypath:  tools.Config.SSH.Keypath}
		out, err := tools.RunOverSSH(&sshConfig, "ansible-playbook -C "+ansibleScript.Path)
		if err == nil {
			return c.JSON(http.StatusCreated, H{
				"<div class='alert alert-success' role='alert'>result</div>": string(out),
			})
			// обработка ошибок
		} else {
			return c.JSON(http.StatusCreated, H{
				"<div class='alert alert-danger' role='alert'>error</div>": err.Error(),
			})

		}
	}
}

func RunAnsibleScriptsAsync(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		// создаём новую задачу
		var ansibleScript models.AnsibleScript
		// привязываем пришедший JSON в новую задачу
		c.Bind(&ansibleScript)
		// добавим задачу с помощью модели
		sshConfig := tools.SSHConfig{
			Host:     tools.Config.SSH.Host,
			Port:     tools.Config.SSH.Port,
			Login:    tools.Config.SSH.Login,
			Password: tools.Config.SSH.Password,
			Keypath:  tools.Config.SSH.Keypath}
		chanId, err := tools.RunOverSshChan(&sshConfig, "ansible-playbook -C "+ansibleScript.Path)
		if err == nil {
			return c.JSON(http.StatusOK, H{"ChanID": chanId})
		} else {
			return c.JSON(http.StatusNotFound, H{"Error": err.Error()})
		}

	}
}

func GetScriptOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		chanId := c.Param("uuid")
		//time.Sleep(time.Millisecond*100)
		return c.JSON(http.StatusOK, H{"cmd": string(<-tools.ChanStack[chanId])})
	}
}
