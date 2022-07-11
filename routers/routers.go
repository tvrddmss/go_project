package routers

import (

    "github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
	//"github.com/swaggo/swag"

    _ "go_project/docs"
    "go_project/pkg/setting"
    "go_project/routers/api"
    "go_project/routers/api/v1"
    "go_project/middleware/jwt"

	
)

type Option func(*gin.Engine)

var options = []Option{}

// 注册app的路由配置
func Include(opts ...Option) {
    options = append(options, opts...)
}
// 初始化
func Init() *gin.Engine {

    r := gin.New()

    r.Use(gin.Logger())

    r.Use(gin.Recovery())

    gin.SetMode(setting.RunMode)

    for _, opt := range options {
        opt(r)
    }
	r.LoadHTMLGlob("templates/**/*")
	r.LoadHTMLGlob("app/**/templates/*")

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

    r.GET("/auth", api.GetAuth)
    r.GET("/clearredis", api.ClearRedis)

    apiv1 := r.Group("/api/v1")
    apiv1.Use(jwt.JWT())
    {
        //获取标签列表
        apiv1.GET("/tags", v1.GetTags)
        //新建标签
        apiv1.POST("/tags", v1.AddTag)
        //更新指定标签
        apiv1.PUT("/tags/:id", v1.EditTag)
        //删除指定标签
        apiv1.DELETE("/tags/:id", v1.DeleteTag)

        //获取文章列表
        apiv1.GET("/articles", v1.GetArticles)
        //获取指定文章
        apiv1.GET("/articles/:id", v1.GetArticle)
        //新建文章
        apiv1.POST("/articles", v1.AddArticle)
        //更新指定文章
        apiv1.PUT("/articles/:id", v1.EditArticle)
        //删除指定文章
        apiv1.DELETE("/articles/:id", v1.DeleteArticle)
    }

    return r
}