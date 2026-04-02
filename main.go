//main 包只能作为程序入口，不能被其他包 import
package main

//导入的模块名和go.mod中一致
import (
	"github.com/gin-gonic/gin"
	"LIST/controller"
	"LIST/dao"
	"LIST/models"
)

func main() {
	//1.mysql中创建数据库 
	//2.连接数据库 
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	sqlDB,err := dao.DB.DB() //获取数据库底层对象
	if err != nil {
		panic("获取数据库对象失败")
	}
	err = sqlDB.Ping() //测试连接
	if err != nil {
		panic("数据库连接失败")
	}
	//3.模型绑定
	dao.DB.AutoMigrate(&models.Todo{})
	//4.程序结束前关闭数据库连接
	defer sqlDB.Close()
//——————————————————————————————————————————————————————————————————
	r := gin.Default()
	// 告诉gin框架模板文件引用的静态文件去哪里找
	// HTML中请求/static下的文件，但后端没有/static的路由，需要告诉他/static去static下找
	r.Static("/static","static")

	// 告诉gin框架去哪里找模板文件
	r.LoadHTMLGlob("templates/*")

	r.GET("/",controller.IndexHandler)

	// API分组
	v1Group := r.Group("v1")
	{
		//待办事项
		//添加待办事项
		v1Group.POST("/todo", controller.CreateATodo)

		//查看所有待办事项
		v1Group.GET("/todo", controller.GetTodoList)

		//查看某一个待办事项(暂时用不到，前端页面没有这个功能)
		v1Group.GET("/todo/:id", func(c *gin.Context) {})

		//修改待办事项
		//基于主键查询更新，先查到数据，然后直接覆写，最后保存的时候由于主键有值，因此update（否则insert）
		v1Group.PUT("/todo/:id", controller.UpdateATodo)

		//删除待办事项
		v1Group.DELETE("/todo/:id", controller.DeleteATodo)
	}

	r.Run(":8090") // listen and serve on 0.0.0.0:8090
}