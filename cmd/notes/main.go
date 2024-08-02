package main

import (
	"github.com/gin-gonic/gin"
	"richingm/LocalDocumentManager/configs"
	"richingm/LocalDocumentManager/internal/infrastructure/mysql"
	"richingm/LocalDocumentManager/router"
)

func main() {
	configs.InitConfig()
	mysql.InitDb()
	r := gin.Default()
	router.InitRouter(r)
	r.Run(configs.ConfigXx.Server.HTTP.Addr)
}
