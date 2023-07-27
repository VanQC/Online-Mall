package controller

import (
	"cloudrestaurant/model"
	"cloudrestaurant/service"
	"cloudrestaurant/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FoodController struct {
	service service.FoodService
}

func (fc *FoodController) FoodRouter(engine *gin.Engine) {
	food := engine.Group("/food/")
	{
		food.GET("category", fc.foodCategory)
		food.POST("category/add", fc.addFoodCategory)
	}
}

// foodCategory 获取食品类别
func (fc *FoodController) foodCategory(c *gin.Context) {
	categories, err := fc.service.GetCategories()
	if err != nil {
		tool.Fail(c, err.Error(), nil)
		return
	}
	tool.Success(c, "食品类型查询成功", gin.H{"data": categories})
}

// addFoodCategory 批量添加食品信息
func (fc *FoodController) addFoodCategory(c *gin.Context) {
	var categories []model.FoodCategory
	if err := c.ShouldBindJSON(&categories); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := fc.service.AddCategories(categories); err != nil {
		tool.Fail(c, "批量添加食品信息失败", nil)
		return
	}
	tool.Success(c, "批量添加食品信息成功", categories)
}
