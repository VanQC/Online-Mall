package controller

import (
	"cloudrestaurant/middlewares"
	"cloudrestaurant/model"
	"cloudrestaurant/service"
	"cloudrestaurant/tool"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProductController struct {
}

func (pc *ProductController) ProductRouter(engine *gin.Engine) {
	product := engine.Group("/api/v1/")
	{
		product.GET("products", pc.ListProducts)
		product.GET("product/:id", pc.ShowProduct)
		product.POST("products", pc.SearchProducts)
		product.GET("imgs/:id", pc.ListProductImg)   // 商品图片
		product.GET("categories", pc.ListCategories) // 商品分类

		auth := product.Group("", middlewares.JwtAuth())
		{
			auth.POST("product", pc.CreateProduct)
			auth.PUT("product/:id", pc.UpdateProduct)
			auth.DELETE("product/:id", pc.DeleteProduct)
		}
	}
}

// 批量添加商品信息
func (pc *ProductController) AddProducts(c *gin.Context) {
	var products []model.Product
	if err := c.ShouldBindJSON(&products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ps := service.ProductService{}
	if err := ps.Adds(c, products); err != nil {
		tool.Fail(c, "批量添加食品信息失败", nil)
		return
	}
	tool.Success(c, "批量添加食品信息成功", products)
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]

	tokenString := c.GetHeader("Authorization")[7:]
	_, claims, _ := tool.ParseToken(tokenString)

	ps := service.ProductService{}
	if err := c.ShouldBind(&ps); err == nil {
		res := ps.Create(c, claims.UserId, files)
		c.JSON(http.StatusOK, res) // todo 这块成功创建后返回值是null
	} else {
		tool.BadRequest(c, gin.H{"参数绑定错误:": err})
		tool.LogrusObj.Infoln(err)
	}
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	ps := service.ProductService{}
	if err := c.ShouldBind(&ps); err == nil {
		res := ps.Update(c.Request.Context(), c.Param("id"))
		c.JSON(200, res)
	} else {
		c.JSON(400, gin.H{"err:": err})
	}
}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	ps := service.ProductService{}
	res := ps.Delete(c.Request.Context(), c.Param("id"))
	c.JSON(200, res)
}

func (pc *ProductController) ListProducts(c *gin.Context) {
	ps := service.ProductService{}
	if err := c.ShouldBind(&ps); err == nil {
		res := ps.List(c.Request.Context())
		c.JSON(200, res)
	} else {
		c.JSON(400, gin.H{"err:": err})
	}
}

func (pc *ProductController) ShowProduct(c *gin.Context) {
	ps := service.ProductService{}
	res := ps.Show(c.Request.Context(), c.Param("id"))
	c.JSON(200, res)
}

func (pc *ProductController) SearchProducts(c *gin.Context) {
	ps := service.ProductService{}
	if err := c.ShouldBind(&ps); err == nil {
		res := ps.Search(c.Request.Context())
		c.JSON(200, res)
	} else {
		c.JSON(400, gin.H{"err:": err})
	}
}

func (pc *ProductController) ListProductImg(c *gin.Context) {
	pis := service.ProductImgService{}
	if err := c.ShouldBind(&pis); err == nil {
		res := pis.List(c.Request.Context(), c.Param("id"))
		c.JSON(200, res)
	} else {
		c.JSON(400, gin.H{"err:": err})
	}
}

func (pc *ProductController) ListCategories(c *gin.Context) {
	lcs := service.ListCategoriesService{}
	if err := c.ShouldBind(&lcs); err == nil {
		res := lcs.List(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
		tool.LogrusObj.Infoln(err)
	}
}
