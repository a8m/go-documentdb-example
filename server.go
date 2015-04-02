// Example of a users REST api using DocumentDB
package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"io/ioutil"
	"net/http"
)

// DocumentDB config
type Config struct {
	Url       string `json:"url"`
	MasterKey string `json:"masterKey"`
}

func main() {
	// Read `url` and `masterKey` from config file
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	config := &Config{}
	if err = json.Unmarshal(data, config); err != nil {
		panic(err)
	}

	// UsersDB
	users := NewDB("test", "users", config)

	// Init Martini, Middleware
	m := martini.Classic()
	m.Use(render.Renderer())

	// Routes
	m.Get("/", func(r render.Render) {
		users, err := users.GetAll()
		if err != nil {
			r.JSON(http.StatusNotFound, err.Error())
		} else {
			r.JSON(200, users)
		}
	})
	m.Post("/", binding.Bind(User{}), func(user User, r render.Render) {
		if err := users.Add(&user); err != nil {
			r.JSON(http.StatusNotFound, err.Error())
		} else {
			r.JSON(201, user)
		}
	})
	m.Get("/:id", func(params martini.Params, r render.Render) {
		if user, err := users.Get(params["id"]); err != nil {
			r.JSON(http.StatusNotFound, err.Error())
		} else {
			r.JSON(http.StatusOK, user)
		}
	})
	m.Put("/:id", binding.Bind(User{}), func(user User, params martini.Params, r render.Render) {
		if err := users.Update(params["id"], &user); err != nil {
			r.JSON(http.StatusNotFound, err.Error())
		}
		r.JSON(http.StatusOK, user)
	})
	m.Delete("/:id", func(params martini.Params, r render.Render) {
		if err := users.Remove(params["id"]); err != nil {
			r.JSON(http.StatusNotFound, err.Error())
		}
		r.JSON(http.StatusNoContent, "")
	})
	m.Run()
}
