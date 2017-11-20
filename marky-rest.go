package main

import (
	"context"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/serkansipahi/corm"
	"github.com/serkansipahi/marky"
	"net/http"
	"strings"
)

// Empty struct
var EmptyStruct = map[string]string{}

// MarkdownHtmlDoc is needed for storing documents
// in couchDB and for json http responses
type MarkdownHtmlDoc struct {
	Id       string `json:"_id,omitempty"`
	Rev      string `json:"_rev,omitempty"`
	Type     string `json:"type,omitempty"`
	Html     string `json:"html,omitempty"`
	Markdown string `json:"markdown,omitempty"`
	Error    string `json:"error,omitempty"`
}

func RenderMarkdown(markdownTemplate string) string {
	if markdownTemplate == "" {
		return markdownTemplate
	}

	// render markdown markup
	markdown := marky.NewMarkdown(markdownTemplate)
	markdownHtml := markdown.Compile()

	return markdownHtml
}

func main() {

	// create couchDB instance
	dbName := "markdown"
	ctx := context.TODO()
	db, err := corm.New(ctx, corm.Config{
		DBName: dbName,
	})

	if err != nil {
		panic("Database: " + dbName + " not found!")
	}

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Extensions: []string{".html"},
		Directory:  "./",
	}))

	m.Post("/markdown/preview", func(res render.Render, req *http.Request) error {

		// response with an empty json when no markdown received
		markdownHtml := RenderMarkdown(req.Header.Get("markdown"))
		if markdownHtml == "" {
			res.JSON(200, EmptyStruct)
			return nil
		}

		res.JSON(200, MarkdownHtmlDoc{
			Html: markdownHtml,
		})

		return nil
	})

	m.Post("/markdown/save", func(res render.Render, req *http.Request) error {

		// response with an empty json when no markdown received
		markdownTemplate := req.Header.Get("markdown")
		markdownHtml := RenderMarkdown(markdownTemplate)
		if markdownHtml == "" {
			res.JSON(200, EmptyStruct)
			return nil
		}

		// save rendered html markdown markup
		docId, _, err := db.Save(ctx, MarkdownHtmlDoc{
			Html:     markdownHtml,
			Markdown: markdownTemplate,
		})

		// something gone wrong
		if err != nil {
			res.JSON(200, MarkdownHtmlDoc{
				Error: "Something gone wrong while saving html document!",
			})
			return nil
		}

		// fine, return success response
		res.JSON(200, MarkdownHtmlDoc{
			Id:   docId,
			Html: markdownHtml,
		})
		return nil

	})

	m.Get("/markdown/get/:docId", func(res render.Render, params martini.Params) error {

		docId := params["docId"]

		// something gone wrong
		docId = strings.Trim(docId, " ")
		if docId == "" {
			res.JSON(200, MarkdownHtmlDoc{
				Error: "No docId received!",
			})
			return nil
		}

		// read document
		var markdown MarkdownHtmlDoc
		_, err = db.Read(ctx, docId, &markdown)
		if err != nil {
			res.JSON(200, MarkdownHtmlDoc{
				Error: "docId: '" + docId + "' not found!",
			})
			return nil
		}

		// fine, return success response
		res.JSON(200, markdown)
		return nil
	})

	m.Run()
}
