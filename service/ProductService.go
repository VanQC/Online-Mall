package service

import (
	"cloudrestaurant/dao"
	"cloudrestaurant/ero"
	"cloudrestaurant/model"
	"cloudrestaurant/serializer"
	"cloudrestaurant/tool"
	"context"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
	"mime/multipart"
	"strconv"
	"sync"
)

// 更新商品的服务
type ProductService struct {
	ID            uint   `form:"id" json:"id"`
	Name          string `form:"name" json:"name"`
	CategoryID    int    `form:"category_id" json:"category_id"`
	Title         string `form:"title" json:"title" `
	Info          string `form:"info" json:"info" `
	ImgPath       string `form:"img_path" json:"img_path"`
	Price         string `form:"price" json:"price"`
	DiscountPrice string `form:"discount_price" json:"discount_price"`
	OnSale        bool   `form:"on_sale" json:"on_sale"`
	Num           int    `form:"num" json:"num"`
	model.BasePage
}

// 批量添加食品信息
func (ps *ProductService) Adds(c *gin.Context, products []model.Product) error {
	pd := dao.NewProductDao(c)
	return pd.AddProducts(products)
}

// 创建商品
func (ps *ProductService) Create(c *gin.Context, id uint, files []*multipart.FileHeader) []serializer.Response {
	var boss *model.Member
	code := ero.SUCCESS

	md := dao.NewMemberDao(c.Request.Context())
	boss, _ = md.QueryById(id)
	// 以第一张作为封面图
	path, err := tool.UploadProductToLocalStatic(c, files[0], id, ps.Name)
	if err != nil {
		code = ero.ErrorUploadFile
		return []serializer.Response{
			{Status: code, Data: ero.GetMsg(code), Error: path},
		}
	}
	product := &model.Product{
		Name:          ps.Name,
		CategoryID:    uint(ps.CategoryID),
		Title:         ps.Title,
		Info:          ps.Info,
		ImgPath:       path,
		Price:         ps.Price,
		DiscountPrice: ps.DiscountPrice,
		Num:           ps.Num,
		OnSale:        true,
		BossID:        id,
		BossName:      boss.UserName,
		BossAvatar:    boss.Avatar,
	}
	productDao := dao.NewProductDao(c.Request.Context())
	err = productDao.CreateProduct(product)
	if err != nil {
		logging.Info(err)
		code = ero.ErrorDatabase
		return []serializer.Response{
			{Status: code, Msg: ero.GetMsg(code), Error: err.Error()},
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan serializer.Response)
	for index, file := range files {
		wg.Add(1)
		go func(index int, file *multipart.FileHeader, errChan chan<- serializer.Response) {
			defer wg.Done()

			num := strconv.Itoa(index)
			productImgDao := dao.NewProductImgDaoByDB(productDao.DB)

			path, err := tool.UploadProductToLocalStatic(c, file, id, ps.Name+num)
			if err != nil {
				errChan <- serializer.Response{
					Status: ero.ErrorUploadFile,
					Data:   ero.GetMsg(code),
					Error:  path,
				}
				return
			}

			productImg := &model.ProductImg{
				ProductID: product.ID,
				ImgPath:   path,
			}
			err = productImgDao.CreateProductImg(productImg)
			if err != nil {
				errChan <- serializer.Response{
					Status: ero.ERROR,
					Msg:    ero.GetMsg(code),
					Error:  err.Error(),
				}
				return
			}
		}(index, file, errChan)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// 处理错误信息
	var res []serializer.Response
	for resp := range errChan {
		res = append(res, resp)
	}
	return res
}

// 更新商品
func (ps *ProductService) Update(ctx context.Context, pId string) serializer.Response {
	code := ero.SUCCESS
	productDao := dao.NewProductDao(ctx)

	productId, _ := strconv.Atoi(pId)
	product := &model.Product{
		Name:       ps.Name,
		CategoryID: uint(ps.CategoryID),
		Title:      ps.Title,
		Info:       ps.Info,
		// ImgPath:       ps.ImgPath,
		Price:         ps.Price,
		DiscountPrice: ps.DiscountPrice,
		OnSale:        ps.OnSale,
	}
	err := productDao.UpdateProduct(uint(productId), product)
	if err != nil {
		logging.Info(err)
		code = ero.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    ero.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    ero.GetMsg(code),
	}
}

// 删除商品
func (ps *ProductService) Delete(ctx context.Context, pId string) serializer.Response {
	code := ero.SUCCESS

	productDao := dao.NewProductDao(ctx)
	productId, _ := strconv.Atoi(pId)
	err := productDao.DeleteProduct(uint(productId))
	if err != nil {
		logging.Info(err)
		code = ero.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    ero.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    ero.GetMsg(code),
	}
}

// 商品
func (ps *ProductService) Show(ctx context.Context, id string) serializer.Response {
	code := ero.SUCCESS
	pId, _ := strconv.Atoi(id)

	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(uint(pId))
	if err != nil {
		logging.Info(err)
		code = ero.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    ero.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   product,
		Msg:    ero.GetMsg(code),
	}
}

func (ps *ProductService) List(ctx context.Context) serializer.Response {
	var products []*model.Product
	code := ero.SUCCESS

	if ps.PageSize == 0 {
		ps.PageSize = 15
	}
	condition := make(map[string]interface{})
	if ps.CategoryID != 0 {
		condition["category_id"] = ps.CategoryID
	}
	productDao := dao.NewProductDao(ctx)
	total, err := productDao.CountProductByCondition(condition)
	if err != nil {
		code = ero.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    ero.GetMsg(code),
		}
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		productDao = dao.NewProductDaoByDB(productDao.DB)
		products, _ = productDao.ListProductByCondition(condition, ps.BasePage)
		wg.Done()
	}()
	wg.Wait()

	return serializer.BuildListResponse(products, uint(total))
}

// 搜索商品
func (ps *ProductService) Search(ctx context.Context) serializer.Response {
	code := ero.SUCCESS
	if ps.PageSize == 0 {
		ps.PageSize = 15
	}
	productDao := dao.NewProductDao(ctx)
	products, err := productDao.SearchProduct(ps.Info, ps.BasePage)
	if err != nil {
		logging.Info(err)
		code = ero.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    ero.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(products, uint(len(products)))
}

type ProductImgService struct {
}

// List 获取商品列表图片
func (service *ProductImgService) List(ctx context.Context, pId string) serializer.Response {
	pis := dao.NewProductImgDao(ctx)
	productId, _ := strconv.Atoi(pId)
	productImgs, _ := pis.ListProductImgByProductId(uint(productId))
	return serializer.BuildListResponse(productImgs, uint(len(productImgs)))
}

type ListCategoriesService struct {
}

func (service *ListCategoriesService) List(ctx context.Context) serializer.Response {
	code := ero.SUCCESS
	categoryDao := dao.NewCategoryDao(ctx)
	categories, err := categoryDao.ListCategory()
	if err != nil {
		logging.Info(err)
		code = ero.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    ero.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    ero.GetMsg(code),
		Data:   categories,
	}
}
