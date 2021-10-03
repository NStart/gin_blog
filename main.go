package main

import (
	"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/redis"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"project/config"
	"project/controllers/admin"
	"project/controllers/index"
	"project/handler"
)

func main() {
	r := gin.Default()

	//基于cookie
	store := cookie.NewStore([]byte(config.SESSION_SECRET))

	// 初始化基于redis的存储引擎
	// 参数说明：
	//    第1个参数 - redis最大的空闲连接数
	//    第2个参数 - 数通信协议tcp或者udp
	//    第3个参数 - redis地址, 格式，host:port
	//    第4个参数 - redis密码
	//    第5个参数 - session加密密钥
	//store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte(config.SESSION_SECRET))

	r.Use(sessions.Sessions(config.SESSION_NAME, store))
	r.Use(handler.CheckLogin())

	r.LoadHTMLGlob("views/***/**/*")
	r.Static("/static", "./static")

	var indexController index.IndexController
	r.GET("/", indexController.Index)
	r.GET("/index", indexController.Index)

	var articleController index.ArticleController
	r.GET("/posts", articleController.Posts)
	r.GET("/posts/detail/:seoLink", articleController.PostsDetail)
	r.GET("/about", articleController.About)
	r.GET("/categories", articleController.Categories)
	r.GET("/tags", articleController.Tags)

	var adminLogin admin.LoginController
	var adminArticle admin.ArticleController
	var adminTags admin.TagsController
	var adminCategories admin.CategoriesController
	v2 := r.Group("/ikebackend")
	{
		v2.GET("/login", adminLogin.Index)
		v2.POST("/dologin", adminLogin.DoLogin)
		v2.GET("/logout", adminLogin.Logout)

		v2.GET("/article/index", adminArticle.Index)
		v2.GET("/article/add", adminArticle.Add)
		v2.POST("/article/doadd", adminArticle.Doadd)
		v2.GET("/article/edit", adminArticle.Edit)
		v2.POST("/article/doedit", adminArticle.Doedit)
		v2.GET("/article/delete", adminArticle.Delete)

		v2.GET("/tags/index", adminTags.Index)
		v2.GET("/tags/add", adminTags.Add)
		v2.POST("/tags/doadd", adminTags.Doadd)
		v2.GET("/tags/edit", adminTags.Edit)
		v2.POST("/tags/doedit", adminTags.Doedit)
		v2.GET("/tags/delete", adminTags.Delete)

		v2.GET("/categories/index", adminCategories.Index)
		v2.GET("/categories/add", adminCategories.Add)
		v2.POST("/categories/doadd", adminCategories.Doadd)
		v2.GET("/categories/edit", adminCategories.Edit)
		v2.POST("/categories/doedit", adminCategories.Doedit)
		v2.GET("/categories/delete", adminCategories.Delete)
	}

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
