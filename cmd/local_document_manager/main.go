package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"richingm/LocalDocumentManager/configs"
)

func main() {
	configs.InitConfig()
	r := gin.Default()

	// 加载静态目录
	r.Static("/static", "static")

	// 加载模板路径
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.tmpl",
			gin.H{
				"treeData": []string{},
			},
		)
	})

	// 根据book的key获取脑图数据
	r.GET("/mind/:doc_key", func(c *gin.Context) {

	})

	// 根据book的key和note的id获取内容
	r.GET("/note/:note_key/node/:node_id", func(c *gin.Context) {

	})

	r.Run(configs.ConfigXx.Server.HTTP.Addr)
}
