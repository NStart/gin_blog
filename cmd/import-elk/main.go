package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"project/models"
	"strconv"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Document struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Table   string `json:"table"`
	SeoLink string `json:"seolink"`
}

func main() {
	article := models.Article{}
	pageSize := 5000
	articleList := make([]*models.Article, pageSize)
	articleList, _, _, err := article.GetArticleList(1, pageSize, 0, 0)
	if err != nil {
		fmt.Println("get data fail")
		return
	}

	article = *articleList[0]
	fmt.Println(articleList[0].Title)

	courseArticle := models.CourseArticle{}
	courseArticleList := make([]*models.CourseArticle, pageSize)
	courseArticleList, _, _, err = courseArticle.GetArticleList(0, 0, 1, pageSize)
	if err != nil {
		fmt.Println("get data fail")
		return
	}

	courseArticle = *courseArticleList[0]
	fmt.Println(courseArticleList[0].ArticleName)

	caCert, err := ioutil.ReadFile("./cmd/import-elk/http_ca.crt")
	if err != nil {
		log.Fatalf("Error reading CA certificate: %s", err)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "EY67+jcZ3S75-u1cYRxu",
		CACert:   caCert,
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}

	var wg sync.WaitGroup
	indexElastic(&wg, es, "article", articleList)
	indexElastic(&wg, es, "course_article", courseArticleList)

	wg.Wait()

}

func indexElastic(wg *sync.WaitGroup, es *elasticsearch.Client, table string, list interface{}) {
	switch v := list.(type) {
	case []*models.Article:
		fmt.Println("article")
		ch := make(chan struct{}, 100)
		for _, item := range v {

			go func(item *models.Article) {
				ch <- struct{}{}
				wg.Add(1)
				articleItem := item
				doc := Document{
					ID:      table + "_" + strconv.Itoa(articleItem.ID),
					Title:   articleItem.Title,
					Content: articleItem.Content,
					SeoLink: articleItem.SeoLink,
				}

				importData(es, doc)

				<-ch
				wg.Done()
			}(item)
		}

	case []*models.CourseArticle:
		fmt.Println("course_article")
		ch := make(chan struct{}, 100)
		for _, item := range v {

			go func(item *models.CourseArticle) {
				ch <- struct{}{}
				wg.Add(1)
				articleItem := item
				doc := Document{
					ID:      table + "_" + strconv.Itoa(articleItem.ID),
					Title:   articleItem.ArticleName,
					Content: articleItem.ArticleContent,
					SeoLink: "",
				}

				importData(es, doc)

				<-ch
				wg.Done()
			}(item)
		}
	}

}

func importData(es *elasticsearch.Client, doc Document) (bool, error) {

	body, err := json.Marshal(doc)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %s", err)
		return false, err
	}

	req := esapi.IndexRequest{
		Index:      "indexname2",
		DocumentID: doc.ID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		fmt.Printf("Error indexing document: %s", err)
		return false, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			fmt.Printf("Error pasing the response body: %s", err)
			return false, err
		}
		log.Fatalf("Error indexing document: %s : %s",
			res.Status(), e["error"].(map[string]interface{})["reason"])
	}

	fmt.Println("ok")
	return true, nil
}
