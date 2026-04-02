package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
  	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)


// 数据库model
type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Status bool `json:"status"` 
}

// 声明全局数据库变量DB
var DB *gorm.DB


func initMySQL() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/projlist?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) 
	if err != nil {
		return
	}
	return 
}

func main() {
	//1.mysql中创建数据库 
	//2.连接数据库 
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	sqlDB,err := DB.DB() //获取数据库底层对象
	if err != nil {
		panic("获取数据库对象失败")
	}
	err = sqlDB.Ping() //测试连接
	if err != nil {
		panic("数据库连接失败")
	}
	//3.模型绑定
	DB.AutoMigrate(&Todo{})
	//4.程序结束前关闭数据库连接
	defer sqlDB.Close()
//——————————————————————————————————————————————————————————————————
	r := gin.Default()
	// 告诉gin框架模板文件引用的静态文件去哪里找
	// HTML中请求/static下的文件，但后端没有/static的路由，需要告诉他/static去static下找
	r.Static("/static","static")

	// 告诉gin框架去哪里找模板文件
	r.LoadHTMLGlob("templates/*")

	r.GET("/",func(c *gin.Context) {
		c.HTML(http.StatusOK,"index.html",nil)
	})

	// API分组
	v1Group := r.Group("v1")
	{
		//待办事项
		//添加待办事项
		v1Group.POST("/todo", func(c *gin.Context) {
			//前端页面填写待办事项，点击提交，发请求到这里
			//1.从请求中把数据拿出来(参数绑定自动提取请求中的数据，自动填充到结构体中)
			var todo Todo
			if err := c.ShouldBindJSON(&todo); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			//2.存入数据库
			//3.返回响应
			if err := DB.Create(&todo).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}	
			c.JSON(http.StatusOK,todo)
		})
		//查看所有待办事项
		v1Group.GET("/todo", func(c *gin.Context) {
			//查询projlist表中的所有数据
			var todoList []Todo
			if err := DB.Find(&todoList).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK,todoList)	
		})

		//查看某一个待办事项(暂时用不到，前端页面没有这个功能)
		v1Group.GET("/todo/:id", func(c *gin.Context) {})

		//修改待办事项
		//基于主键查询更新，先查到数据，然后直接覆写，最后保存的时候由于主键有值，因此update（否则insert）
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id,ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
				return
			}
			var todo Todo
			if err := DB.Where("id=?",id).First(&todo).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			if err := c.ShouldBindJSON(&todo); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}	
			if err := DB.Save(&todo).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}	
			c.JSON(http.StatusOK,todo)
		})
		//删除待办事项
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id,ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
				return
			}	
			if err := DB.Where("id=?",id).Delete(&Todo{}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}	
			c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
		})

	}

	r.Run(":8090") // listen and serve on 0.0.0.0:8090
}