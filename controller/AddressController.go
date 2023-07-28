package controller

import (
	"cloudrestaurant/middlewares"
	"cloudrestaurant/service"
	"cloudrestaurant/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddressController struct {
}

func (adc *AddressController) AdderssRouter(engine *gin.Engine) {
	auth := engine.Group("/api/v1/", middlewares.JwtAuth())
	{
		auth.POST("addresses", adc.CreateAddress)
		auth.GET("addresses/:id", adc.GetAddress)
		auth.GET("addresses", adc.ListAddress)
		auth.PUT("addresses/:id", adc.UpdateAddress)
		auth.DELETE("addresses/:id", adc.DeleteAddress)
	}
}

// CreateAddress 新增收货地址
func (adc *AddressController) CreateAddress(c *gin.Context) {
	ads := service.AddressService{}
	tokenString := c.GetHeader("Authorization")[7:]
	_, claim, _ := tool.ParseToken(tokenString)
	if err := c.ShouldBind(&ads); err == nil {
		res := ads.Create(c.Request.Context(), claim.UserId)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}

// GetAddress 展示某个收货地址
func (adc *AddressController) GetAddress(c *gin.Context) {
	ads := service.AddressService{}
	res := ads.Show(c.Request.Context(), c.Param("id"))
	c.JSON(http.StatusOK, res)
}

// ListAddress 展示收货地址
func (adc *AddressController) ListAddress(c *gin.Context) {
	ads := service.AddressService{}
	tokenString := c.GetHeader("Authorization")[7:]
	_, claim, _ := tool.ParseToken(tokenString)
	if err := c.ShouldBind(&ads); err == nil {
		res := ads.List(c.Request.Context(), claim.UserId)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}

// UpdateAddress 修改收货地址
func (adc *AddressController) UpdateAddress(c *gin.Context) {
	ads := service.AddressService{}
	tokenString := c.GetHeader("Authorization")[7:]
	_, claim, _ := tool.ParseToken(tokenString)
	if err := c.ShouldBind(&ads); err == nil {
		res := ads.Update(c.Request.Context(), claim.UserId, c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}

// DeleteAddress 删除收获地址
func (adc *AddressController) DeleteAddress(c *gin.Context) {
	ads := service.AddressService{}
	if err := c.ShouldBind(&ads); err == nil {
		res := ads.Delete(c.Request.Context(), c.Param("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}
