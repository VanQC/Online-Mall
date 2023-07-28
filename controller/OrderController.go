package controller

import (
	"cloudrestaurant/middlewares"
	"cloudrestaurant/service"
	"cloudrestaurant/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderController struct {
}

func (odc *OrderController) OrderRouter(engine *gin.Engine) {
	auth := engine.Group("/api/v1/", middlewares.JwtAuth())
	{
		auth.POST("orders", odc.CreateOrder)
		auth.GET("orders", odc.ListOrders)
		auth.GET("orders/:id", odc.ShowOrder)
		auth.DELETE("orders/:id", odc.DeleteOrder)
	}
}

func (odc *OrderController) CreateOrder(c *gin.Context) {
	ordServ := service.OrderService{}

	tokenString := c.GetHeader("Authorization")[7:]
	_, claim, _ := tool.ParseToken(tokenString)

	if err := c.ShouldBind(&ordServ); err == nil {
		res := ordServ.Create(c.Request.Context(), claim.UserId)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}

func (odc *OrderController) ListOrders(c *gin.Context) {
	ordServ := service.OrderService{}

	tokenString := c.GetHeader("Authorization")[7:]
	_, claim, _ := tool.ParseToken(tokenString)

	if err := c.ShouldBind(&ordServ); err == nil {
		res := ordServ.List(c.Request.Context(), claim.UserId)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}

// 订单详情
func (odc *OrderController) ShowOrder(c *gin.Context) {
	ordServ := service.OrderService{}
	if err := c.ShouldBind(&ordServ); err == nil {
		res := ordServ.Show(c.Request.Context(), c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}

func (odc *OrderController) DeleteOrder(c *gin.Context) {
	ordServ := service.OrderService{}
	if err := c.ShouldBind(&ordServ); err == nil {
		res := ordServ.Delete(c.Request.Context(), c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}
