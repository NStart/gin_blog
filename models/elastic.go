package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"

	"project/helpers"
)

type EArticle struct {
	BaseModel
	Link       string
	DocumentId string
	Title      string `gorm:"type:varchar(100);not null"`
	Content    string `gorm:"not null"`
}

func (article *EArticle) GetArticleList(pageIndex, pageSize int, keyword string) ([]*EArticle, int, int, error) {
	var err error
	var articleList []*EArticle
	var common helpers.Common

	var iTotalCount, totalPage int

	caCert, err := ioutil.ReadFile("./cmd/import-elk/http_ca.crt")
	if err != nil {
		common.GetErrorLocation()
		fmt.Printf("Error reading CA certificate:  %s %d", err, common.GetErrorLocation()["line"])
		return articleList, iTotalCount, totalPage, err
	}
	//fmt.Println(caCert)

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
		fmt.Printf("Error creating the Elasticsearch client:  %s %d", err, common.GetErrorLocation()["line"])
		return articleList, iTotalCount, totalPage, err
	}

	// 定义查询条件
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": make([]interface{}, 0), // 初始化为空切片
			},
		},
		"size": pageSize,
		"from": (pageIndex - 1) * pageSize,
	}

	// 如果 keyword 不为空，则添加 title 和 content 的查询条件
	if keyword != "" {
		must := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{}) // 将 interface{} 转换为 []interface{}
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"title": map[string]interface{}{
					"query": keyword,
				},
			},
		})
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"content": map[string]interface{}{
					"query": keyword,
				},
			},
		})

		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = must
	}

	body, err := json.Marshal(query)
	if err != nil {
		fmt.Printf("Error marshaling JSON:  %s %d", err, common.GetErrorLocation()["line"])
		return articleList, iTotalCount, totalPage, err
	}

	req := esapi.SearchRequest{
		Index: []string{"indexname2"}, // 替换为您的索引名称
		Body:  bytes.NewReader(body),
	}

	// 执行搜索请求
	res, err := req.Do(context.Background(), es)
	if err != nil {
		fmt.Printf("Error searching documents:  %s %d", err, common.GetErrorLocation()["line"])
		return articleList, iTotalCount, totalPage, err
	}
	defer res.Body.Close()

	// 检查搜索结果是否出错
	if res.IsError() {
		fmt.Printf("Error searching documents: %s", res.Status())
		return articleList, iTotalCount, totalPage, err
	}

	// 解析搜索结果
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		fmt.Printf("Error parsing the response body:  %s %d", err, common.GetErrorLocation()["line"])
		return articleList, iTotalCount, totalPage, err
	}

	// 输出搜索结果
	hits, _ := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		fmt.Printf("Document ID: %s, Title: %s\n", hit.(map[string]interface{})["_id"], source["title"])

		documentId := hit.(map[string]interface{})["_id"].(string)
		seoLink := source["seolink"].(string)
		parts := strings.Split(documentId, "_")
		fmt.Println(parts)
		table := parts[0]
		articleIdStr := parts[1]
		var link string
		if table == "article" {
			link = "/posts/detail/" + seoLink
		} else {
			articleIdStr = parts[2]
			link = "/course-article/detail/" + articleIdStr
		}
		article = &EArticle{
			Link:       link,
			DocumentId: documentId,
			Title:      source["title"].(string),
			Content:    source["content"].(string),
		}
		//fmt.Println(article)
		articleList = append(articleList, article)
	}

	eArticle := &EArticle{}
	iTotalCount, err = eArticle.GetArticleCount(keyword)
	if err != nil {
		fmt.Printf("Get total count error:  %s %d", err, common.GetErrorLocation()["line"])
		return articleList, iTotalCount, totalPage, err
	}
	totalPage = int(math.Ceil(float64(iTotalCount) / float64(pageSize)))

	//fmt.Printf("%+v %+v", articleList, pageSize)
	//fmt.Println(articleList)
	return articleList, iTotalCount, totalPage, nil
}

func (article *EArticle) GetArticleCount(keyword string) (int, error) {
	var err error
	var common helpers.Common

	caCert, err := ioutil.ReadFile("./cmd/import-elk/http_ca.crt")
	if err != nil {
		fmt.Printf("Error reading CA certificate:  %s %d", err, common.GetErrorLocation()["line"])
		return 0, nil
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
		fmt.Printf("Error creating the Elasticsearch client:  %s %d", err, common.GetErrorLocation()["line"])
		return 0, nil
	}

	// 定义查询条件
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": make([]interface{}, 0), // 初始化为空切片
			},
		},
	}

	// 如果 keyword 不为空，则添加 title 和 content 的查询条件
	if keyword != "" {
		must := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{}) // 将 interface{} 转换为 []interface{}
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"title": map[string]interface{}{
					"query": keyword,
				},
			},
		})
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"content": map[string]interface{}{
					"query": keyword,
				},
			},
		})

		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = must
	}

	body, err := json.Marshal(query)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %s %d", err, common.GetErrorLocation()["line"])
		return 0, nil
	}

	req := esapi.CountRequest{
		Index: []string{"indexname2"}, // 替换为您的索引名称
		Body:  bytes.NewReader(body),
	}

	// 执行搜索请求
	res, err := req.Do(context.Background(), es)
	if err != nil {
		fmt.Printf("Error searching documents:  %s %d", err, common.GetErrorLocation()["line"])
		return 0, nil
	}
	defer res.Body.Close()

	// 检查搜索结果是否出错
	if res.IsError() {
		fmt.Printf("Error searching documents: %s %d", res.Status(), common.GetErrorLocation()["line"])
		return 0, nil
	}

	// 解析搜索结果
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		fmt.Printf("Error parsing the response body: %s %d", err, common.GetErrorLocation()["line"])
		return 0, nil
	}

	fmt.Println(11111)
	fmt.Println(result)
	// 获取匹配文档的总数量
	totalHits := int(result["count"].(float64))

	//fmt.Printf("%+v %+v", articleList, pageSize)
	return totalHits, nil
}
