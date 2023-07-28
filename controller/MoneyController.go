package controller

import (
	"cloudrestaurant/middlewares"
	"cloudrestaurant/service"
	"cloudrestaurant/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MoneyController struct {
}

func (pc *MoneyController) PayRouter(engine *gin.Engine) {
	auth := engine.Group("/api/v1/", middlewares.JwtAuth())

	{
		auth.POST("money", pc.ShowMoney)  // 显示金额
		auth.POST("paydown", pc.OrderPay) // 支付功能
	}
}

func (pc *MoneyController) ShowMoney(c *gin.Context) {
	showMoneyService := service.ShowMoneyService{}

	tokenString := c.GetHeader("Authorization")[7:]
	_, claim, _ := tool.ParseToken(tokenString)

	if err := c.ShouldBind(&showMoneyService); err == nil {
		res := showMoneyService.Show(c.Request.Context(), claim.UserId)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}

func (pc *MoneyController) OrderPay(c *gin.Context) {
	orderPay := service.OrderPayService{}

	tokenString := c.GetHeader("Authorization")[7:]
	_, claim, _ := tool.ParseToken(tokenString)

	if err := c.ShouldBind(&orderPay); err == nil {
		res := orderPay.PayDown(c.Request.Context(), claim.UserId)
		c.JSON(http.StatusOK, res)
	} else {
		tool.LogrusObj.Infoln(err)
		c.JSON(http.StatusBadRequest, err)
	}
}
