package controller

import (
	"cloudrestaurant/middlewares"
	"cloudrestaurant/service"
	"cloudrestaurant/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FavoriteController struct {
}

// 收藏夹
func (fc *FavoriteController) FavoriteRouter(engine *gin.Engine) {
	auth := engine.Group("/api/v1/", middlewares.JwtAuth())
	{
		auth.GET("favorites", fc.ShowFavorites)
		auth.POST("favorites", fc.CreateFavorite)
		auth.DELETE("favorites/:id", fc.DeleteFavorite)
	}
}

// 创建收藏
func (fc *FavoriteController) CreateFavorite(c *gin.Context) {
	fs := service.FavoritesService{}

	tokenString := c.GetHeader("Authorization")[7:]
	_, claim, _ := tool.ParseToken(tokenString)

	if err := c.ShouldBind(&fs); err == nil {
		res := fs.Create(c.Request.Context(), claim.UserId)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}

// 收藏夹详情接口
func (fc *FavoriteController) ShowFavorites(c *gin.Context) {
	fs := service.FavoritesService{}
	tokenString := c.GetHeader("Authorization")[7:]
	_, claim, _ := tool.ParseToken(tokenString)

	if err := c.ShouldBind(&fs); err == nil {
		res := fs.Show(c.Request.Context(), claim.UserId)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}

func (fc *FavoriteController) DeleteFavorite(c *gin.Context) {
	fs := service.FavoritesService{}
	if err := c.ShouldBind(&fs); err == nil {
		res := fs.Delete(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}
