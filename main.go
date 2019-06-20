package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	"html/template"
	"net/http"
)

type Post struct {
	Title       string `json:"title"`
	Date        string `json:"date"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Data        string `json:"data"`
}
type IndexPost struct {
	Id          string
	Title       string
	Date        string
	Author      string
	Description string
	Data        string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var posts []IndexPost

	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
	)
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Search().
		Index("test_index").
		Pretty(true).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	for _, hit := range res.Hits.Hits {
		var localPost IndexPost
		localPost.Id = hit.Id
		err := json.Unmarshal(hit.Source, &localPost)
		if err != nil {
			fmt.Println(err)
		}
		posts = append(posts, localPost)
	}
	tmpl, _ := template.ParseFiles("templates/index.html")
	_ = tmpl.Execute(w, posts)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		title := r.FormValue("title")
		date := r.FormValue("date")
		author := r.FormValue("author")
		description := r.FormValue("description")
		data := r.FormValue("data")

		client, err := elastic.NewClient(
			elastic.SetURL("http://localhost:9200"),
		)
		if err != nil {
			fmt.Println(err)
		}
		_, err = client.Index().
			Index("test_index").
			Id(id).
			BodyJson(&Post{
				Title:       title,
				Date:        date,
				Author:      author,
				Description: description,
				Data:        data,
			}).
			Refresh("true").
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		http.Redirect(w, r, "/index", 301)
		return
	}
	http.ServeFile(w, r, "templates/add.html")

}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		client, err := elastic.NewClient(
			elastic.SetURL("http://localhost:9200"),
		)
		if err != nil {
			fmt.Println(err)
		}
		q := elastic.NewMultiMatchQuery(id, "_id").Type("phrase")

		res, err := client.Search().
			Index("test_index").
			Pretty(true).
			Query(q).
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		var post IndexPost
		post.Id = id
		for _, hit := range res.Hits.Hits {
			err := json.Unmarshal(hit.Source, &post)
			if err != nil {
			}
			break
		}

		if err != nil {
			fmt.Println(err)
		}
		tmpl, _ := template.ParseFiles("templates/update.html")
		_ = tmpl.Execute(w, post)
		return
	}
	http.Redirect(w, r, "/index", 301)
}
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		client, err := elastic.NewClient(
			elastic.SetURL("http://localhost:9200"),
		)
		if err != nil {
			fmt.Println(err)
		}
		id := r.FormValue("id")
		ctx := context.Background()
		var params = map[string]interface{}{
			"title":       r.FormValue("title"),
			"author":      r.FormValue("author"),
			"date":        r.FormValue("date"),
			"data":        r.FormValue("data"),
			"description": r.FormValue("description"),
		}
		_, err = client.Update().Index("test_index").Id(id).Refresh("true").
			Script(elastic.NewScriptInline(
				"ctx._source.title = params.title;" +
					"ctx._source.description = params.description;" +
					"ctx._source.author = params.author;" +
					"ctx._source.date = params.date;" +
					"ctx._source.data = params.data;").
				Params(params)).
			Do(ctx)
		if err != nil {
			fmt.Println(err)
		}

	}
	http.Redirect(w, r, "/index", 301)
}
func IdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
	)
	if err != nil {
		fmt.Println(err)
	}
	q := elastic.NewMultiMatchQuery(id, "_id").Type("phrase")
	res, err := client.Search().
		Index("test_index").
		Pretty(true).
		Query(q).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	var post IndexPost
	post.Id = id
	for _, hit := range res.Hits.Hits {
		err := json.Unmarshal(hit.Source, &post)
		if err != nil {
		}
		break
	}
	jsonId, _ := json.Marshal(post)
	_, _ = fmt.Fprint(w, string(jsonId))
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		client, err := elastic.NewClient(
			elastic.SetURL("http://localhost:9200"),
		)
		if err != nil {
			fmt.Println(err)
		}
		id := r.FormValue("id")
		_, err = client.Delete().
			Index("test_index").
			Id(id).
			Refresh("true").
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}
	http.Redirect(w, r, "/index", 301)
	return
}
func main() {
	http.HandleFunc("/add", AddHandler)
	http.HandleFunc("/edit", EditHandler)
	http.HandleFunc("/index", IndexHandler)
	http.HandleFunc("/delete", DeleteHandler)
	http.HandleFunc("/update", UpdateHandler)

	router := mux.NewRouter()
	router.HandleFunc("/id/{id:[0-9]+}", IdHandler)
	http.Handle("/", router)
	_ = http.ListenAndServe(":8182", nil)
}
