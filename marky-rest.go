package main

import (
	"context"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/serkansipahi/corm"
	"github.com/serkansipahi/marky"
	"net/http"
)

type Markdown struct {
	Id       string `json:"_id,omitempty"`
	Rev      string `json:"_rev,omitempty"`
	Type     string `json:"type"`
	Markdown string `json:"Markdown"`
}

type DocId struct {
	id string
}

func main() {

	// create couchdb instance
	dbName := "markdown"
	ctx := context.TODO()
	db, err := corm.New(ctx, corm.Config{
		DBName: dbName,
	})

	if err != nil {
		panic("Database: " + dbName + " not found!")
	}

	m := martini.Classic()
	m.Use(render.Renderer())

	m.Post("/markdown/save", func(res render.Render, req *http.Request) {
		if markdown := req.Header.Get("markdown"); markdown != "" {

			markdown := marky.NewMarkdown(markdown)
			html := markdown.Compile()
			if html == "" {
				res.JSON(200, map[string]string{"id": ""})
			}

			docId, _, err := db.Save(ctx, Markdown{
				Markdown: markdown.Compile(),
			})
			if err != nil {
				panic(err)
			}
			res.JSON(200, map[string]string{"id": docId})
		}

	})

	m.Get("/markdown/get/:docId", func(res render.Render, params martini.Params) {

		docId := params["docId"]
		var markdown Markdown
		_, err = db.Read(ctx, docId, &markdown)
		res.JSON(200, markdown)
	})

	m.Run()
}
