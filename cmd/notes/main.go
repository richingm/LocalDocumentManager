package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"richingm/LocalDocumentManager/configs"
	"richingm/LocalDocumentManager/internal/application"
	"runtime"
	"strings"
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
		c.HTML(
			http.StatusOK,
			"index.tmpl",
			gin.H{},
		)
	})

	// 根据note的key获取脑图数据
	r.GET("/menus", func(c *gin.Context) {
		type response struct {
			Status   int    `json:"status"`
			Content  string `json:"content"`
			ErrorMsg string `json:"error_msg"`
		}
		var res response
		res.Status = http.StatusOK
		menuService := application.NewMenuService()
		menuDtoList, err := getMenuList(menuService)
		if err != nil {
			res.ErrorMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}

		res.Content = menuService.ConvertMenusToString(menuDtoList)
		c.JSON(http.StatusOK, res)
	})

	// 根据note的key获取脑图数据
	r.GET("/mind/:note_key", func(c *gin.Context) {
		type response struct {
			Status   int                 `json:"status"`
			Content  application.NodeDto `json:"content"`
			ErrorMsg string              `json:"error_msg"`
		}
		var res response
		res.Status = http.StatusOK

		noteKey := c.Param("note_key")
		noteTreeService := application.NewNoteTreeService()

		menuService := application.NewMenuService()
		menuDtoList, err := getMenuList(menuService)
		note := noteTreeService.GetNote(menuDtoList, noteKey)
		if err != nil {
			res.ErrorMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		_, err = os.Stat(note.DirPath)
		if err != nil {
			res.ErrorMsg = "知识库不存在"
			c.JSON(http.StatusOK, res)
			return
		}

		nodeService := application.NewNodeService()
		nodes, err := nodeService.GetMind(note.DirPath, note.MenuName, getDisplayLevel(3), FIleSuffix)
		if err != nil {
			res.ErrorMsg = "知识库不存在"
			c.JSON(http.StatusOK, res)
			return
		}
		res.Content = nodes
		c.JSON(http.StatusOK, res)
	})

	// 根据book的key和note的id获取内容
	r.GET("/note/:note_key/node/:node_id", func(c *gin.Context) {
		type response struct {
			Title    string `json:"title"`
			Status   int    `json:"status"`
			Content  string `json:"content"`
			ErrorMsg string `json:"error_msg"`
		}
		var res response
		noteKey := c.Param("note_key")
		nodeId := c.Param("node_id")
		menuService := application.NewMenuService()
		menuDtoList, err := getMenuList(menuService)
		noteTreeService := application.NewNoteTreeService()
		note := noteTreeService.GetNote(menuDtoList, noteKey)
		if err != nil {
			res.ErrorMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}
		nodeService := application.NewNodeService()
		title, content, err := nodeService.GetContentAndTitle(note.DirPath, note.MenuName, nodeId, FIleSuffix)
		if err != nil {
			res.ErrorMsg = err.Error()
			c.JSON(http.StatusOK, res)
			return
		}

		// 处理图片
		content = nodeService.ExtractImagePaths(note.DirPath, content, getRootPath())

		res.Title = title
		res.Status = http.StatusOK
		res.Content = content
		c.JSON(http.StatusOK, res)
	})

	r.Run(configs.ConfigXx.Server.HTTP.Addr)
}

func getMenuList(menuService *application.MenuService) ([]application.MenuDto, error) {
	rootPath := getRootPath()
	menuDtoList := make([]application.MenuDto, 0)
	for _, docPath := range configs.ConfigXx.Docs {
		realPath := fmt.Sprintf("%s/../../%s/%s", strings.TrimRight(rootPath, "/"),
			strings.TrimRight(strings.TrimLeft(docPath, "/"), "/"),
			configs.ConfigXx.MenusYamlFile)
		menuDtos, err := menuService.GetMenus(realPath)
		if err != nil {
			return nil, err
		}
		menuDtoList = append(menuDtoList, menuDtos...)
	}
	return menuDtoList, nil
}

func getDisplayLevel(level int64) int64 {
	if level > 0 {
		return level
	}
	return configs.ConfigXx.DefaultDisplayLevel
}

// 获取main文件的执行路径
func getRootPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("")
	}
	return filepath.Dir(filename)
}
