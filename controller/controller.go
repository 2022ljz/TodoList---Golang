package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"LIST/dao"
	"LIST/models"
)


/*
url --> controller --> service --> model
请求来了 --> 控制器 --> 业务逻辑 --> 数据的增删改查
*/
//这个业务太简单了所以controller直接调用model了，service层就省略了


//注意首字母大写包导出
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK,"index.html",nil)
}

func CreateATodo(c *gin.Context) {
	//前端页面填写待办事项，点击提交，发请求到这里
	//1.从请求中把数据拿出来(参数绑定自动提取请求中的数据，自动填充到结构体中)
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//2.存入数据库
	//3.返回响应
	if err := dao.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	
	c.JSON(http.StatusOK,todo)
}

func GetTodoList(c *gin.Context) {
	//查询projlist表中的所有数据
	var todoList []models.Todo
	if err := dao.DB.Find(&todoList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK,todoList)	
}

//基于主键查询更新，先查到数据，然后直接覆写，最后保存的时候由于主键有值，因此update（否则insert）
func UpdateATodo(c *gin.Context) {
	id,ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}
	var todo models.Todo
	if err := dao.DB.Where("id=?",id).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	
	if err := dao.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	
	c.JSON(http.StatusOK,todo)
}

func DeleteATodo(c *gin.Context) {
	id,ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}	
	if err := dao.DB.Where("id=?",id).Delete(&models.Todo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	
	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}