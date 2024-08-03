package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"richingm/LocalDocumentManager/internal/application"
	"strconv"
)

const (
	FIleSuffix = "md"
)

func InitRouter(r *gin.Engine) {
	// 加载静态目录
	r.Static("/static", "static")

	// 加载模板路径
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.tmpl",
			gin.H{},
		)
	})

	// 根据note的key获取脑图数据
	r.GET("/categories", func(c *gin.Context) {
		type response struct {
			Status   int    `json:"status"`
			Content  string `json:"content"`
			ErrorMsg string `json:"error_msg"`
		}
		var res response
		res.Status = http.StatusOK
		categoryService := application.NewCategoryService(c.Request.Context())
		content, err := categoryService.ListHtml(c.Request.Context())
		if err != nil {
			res.ErrorMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		res.Content = content
		c.JSON(http.StatusOK, res)
	})

	// 根据note的key获取脑图数据
	r.GET("/mind/:note_id", func(c *gin.Context) {
		type response struct {
			Status   int                 `json:"status"`
			Content  application.NodeDto `json:"content"`
			ErrorMsg string              `json:"error_msg"`
		}
		var res response
		res.Status = http.StatusOK

		noteIdStr := c.Param("note_id")
		noteId, err := strconv.Atoi(noteIdStr)
		if err != nil {
			res.ErrorMsg = "参数错误"
			c.JSON(http.StatusOK, res)
			return
		}

		articleService := application.NewArticleService(c.Request.Context())
		nodes, err := articleService.Nodes(c.Request.Context(), noteId)
		if err != nil {
			res.ErrorMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		res.Content = nodes
		c.JSON(http.StatusOK, res)
	})

	// 根据book的key和note的id获取内容
	r.GET("/note/:article_id", func(c *gin.Context) {
		type response struct {
			Title    string `json:"title"`
			Status   int    `json:"status"`
			Content  string `json:"content"`
			ErrorMsg string `json:"error_msg"`
		}
		var res response

		articleIdStr := c.Param("article_id")
		articleId, err := strconv.Atoi(articleIdStr)
		if err != nil {
			res.ErrorMsg = "参数错误"
			c.JSON(http.StatusOK, res)
			return
		}
		articleService := application.NewArticleService(c.Request.Context())
		dto, err := articleService.Get(c.Request.Context(), articleId)
		if err != nil {
			res.ErrorMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		res.Title = dto.Title
		res.Status = http.StatusOK
		res.Content = dto.Content
		c.JSON(http.StatusOK, res)
	})

	r.POST("/article", func(c *gin.Context) {
		type response struct {
			Status   int    `json:"status"`
			ErrorMsg string `json:"error_msg"`
		}

		var res response
		res.Status = http.StatusOK

		type articleParams struct {
			Cid     string `json:"cid" form:"cid"`
			Pid     string `json:"pid" form:"pid"`
			Title   string `form:"title" json:"title"`
			Content string `form:"content" json:"content"`
		}

		var param articleParams
		err := c.ShouldBindJSON(&param)
		if err != nil {
			res.ErrorMsg = "参数错误"
			c.JSON(http.StatusOK, res)
			return
		}

		cid, err := strconv.Atoi(param.Cid)
		if err != nil {
			res.ErrorMsg = "参数错误"
			c.JSON(http.StatusOK, res)
			return
		}

		pid, err := strconv.Atoi(param.Pid)
		if err != nil {
			res.ErrorMsg = "参数错误"
			c.JSON(http.StatusOK, res)
			return
		}

		articleService := application.NewArticleService(c.Request.Context())
		err = articleService.Create(c.Request.Context(), cid, pid, param.Title, param.Content)
		if err != nil {
			res.ErrorMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
