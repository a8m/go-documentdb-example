package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"io/ioutil"
)

type Config struct {
	Url       string `json:"url"`
	MasterKey string `json:"masterKey"`
}

func main() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	config := &Config{}
	if err = json.Unmarshal(data, config); err != nil {
		panic(err)
	}

	// Init Martini, Middleware
	m := martini.Classic()
	m.Use(render.Renderer())

	// Routes, Controller
	m.Get("/:id", func() {})
	m.Post("/", binding.Bind(User{}), func(user User, r render.Render) {
		r.JSON(200, user)
	})
	m.Run()
}
