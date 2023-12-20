package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"richingm/LocalDocumentManager/configs"
	"richingm/LocalDocumentManager/internal/application"
)

const (
	FIleSuffix = "md"
)

func main() {
	configs.InitConfig()
	r := gin.Default()

	// 加载静态目录
	r.Static("/static", "static")

	// 加载模板路径
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		noteTreeService := application.NewNoteTreeService()
		treeList := noteTreeService.GetTree(configs.ConfigXx)
		c.HTML(
			http.StatusOK,
			"index.tmpl",
			gin.H{
				"treeList": treeList,
				"aaa":      "花花草草",
			},
		)
	})

	// 根据note的key获取脑图数据
	r.GET("/mind/:note_key", func(c *gin.Context) {
		noteKey := c.Param("note_key")
		noteTreeService := application.NewNoteTreeService()
		noteName, dir, err := noteTreeService.GetDirAndName(configs.ConfigXx, noteKey)
		if err != nil {
			c.JSON(http.StatusNotFound, err)
			return
		}
		nodeService := application.NewNodeService()
		nodes, err := nodeService.GetMind(dir, noteName, 3, FIleSuffix)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, nodes)
	})

	// 根据book的key和note的id获取内容
	r.GET("/note/:note_key/node/:node_id", func(c *gin.Context) {
		type response struct {
			Status   int    `json:"status"`
			Content  string `json:"content"`
			ErrorMsg string `json:"error_msg"`
		}
		var res response
		noteKey := c.Param("note_key")
		nodeId := c.Param("node_id")
		noteTreeService := application.NewNoteTreeService()
		noteName, dir, err := noteTreeService.GetDirAndName(configs.ConfigXx, noteKey)
		if err != nil {
			res.ErrorMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		nodeService := application.NewNodeService()
		content, err := nodeService.GetContent(dir, noteName, nodeId, FIleSuffix)
		if err != nil {
			res.ErrorMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		res.Status = http.StatusOK
		res.Content = content
		c.JSON(http.StatusOK, res)
	})

	r.Run(configs.ConfigXx.Server.HTTP.Addr)
}
