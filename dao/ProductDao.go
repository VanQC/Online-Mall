package dao

import (
	"cloudrestaurant/model"
	"context"
	"gorm.io/gorm"
)

type ProductDao struct {
	*gorm.DB
}

func NewProductDao(ctx context.Context) *ProductDao {
	return &ProductDao{DB.WithContext(ctx)}
}

func NewProductDaoByDB(db *gorm.DB) *ProductDao {
	return &ProductDao{db}
}

func (dao *ProductDao) CreateProducts(product *model.Product) error {
	return dao.DB.Model(&model.Product{}).Create(&product).Error
}

// 利用事务机制批量添加食品类别信息
func (dao *ProductDao) AddProducts(products []model.Product) error {
	// 开启事务
	tx := dao.DB.Begin()
	// 捕获任何可能发生的panic，并在回滚事务后重新抛出异常
	// 这是一种常用的做法，可以确保即使函数内部出现了异常，也可以保证事务的完整性
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 先检查事务对象是否有错误，如果有则直接返回
	if err := tx.Error; err != nil {
		return err
	}
	// 在事务中执行一系列写入操作
	for _, product := range products {
		if err := tx.Create(&product).Error; err != nil {
			tx.Rollback() // 遇到错误时进行回滚操作
			return err
		}
	}
	// 提交事务
	return tx.Commit().Error
}

// CreateProduct 创建商品
func (dao *ProductDao) CreateProduct(product *model.Product) error {
	return dao.DB.Model(&model.Product{}).Create(&product).Error
}

// UpdateProduct 更新商品
func (dao *ProductDao) UpdateProduct(pId uint, product *model.Product) error {
	return dao.DB.Model(&model.Product{}).Where("id=?", pId).
		Updates(&product).Error
}

// DeleteProduct 删除商品
func (dao *ProductDao) DeleteProduct(pId uint) error {
	return dao.DB.Model(&model.Product{}).Where("id=?", pId).
		Delete(&model.Product{}).Error
}

// GetProductById 通过 id 获取product
func (dao *ProductDao) GetProductById(id uint) (product *model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).Where("id=?", id).First(&product).Error
	return
}

// CountProductByCondition 根据情况获取商品的数量
func (dao *ProductDao) CountProductByCondition(condition map[string]interface{}) (total int64, err error) {
	err = dao.DB.Model(&model.Product{}).Where(condition).Count(&total).Error
	return
}

// ListProductByCondition 获取商品列表
func (dao *ProductDao) ListProductByCondition(condition map[string]interface{}, page model.BasePage) (products []*model.Product, err error) {
	err = dao.DB.Where(condition).
		Offset((page.PageNum - 1) * page.PageSize).
		Limit(page.PageSize).Find(&products).Error
	return
}

// SearchProduct 搜索商品
func (dao *ProductDao) SearchProduct(info string, page model.BasePage) (products []*model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where("name LIKE ? OR info LIKE ?", "%"+info+"%", "%"+info+"%").
		Offset((page.PageNum - 1) * page.PageSize).
		Limit(page.PageSize).Find(&products).Error
	return
}

type ProductImgDao struct {
	*gorm.DB
}

func NewProductImgDao(ctx context.Context) *ProductImgDao {
	return &ProductImgDao{DB.WithContext(ctx)}
}

func NewProductImgDaoByDB(db *gorm.DB) *ProductImgDao {
	return &ProductImgDao{db}
}

// CreateProductImg 创建商品图片
func (dao *ProductImgDao) CreateProductImg(productImg *model.ProductImg) (err error) {
	err = dao.DB.Model(&model.ProductImg{}).Create(&productImg).Error
	return
}

// ListProductImgByProductId 根据商品id获取商品图片
func (dao *ProductImgDao) ListProductImgByProductId(pId uint) (products []*model.ProductImg, err error) {
	err = dao.DB.Model(&model.ProductImg{}).
		Where("product_id=?", pId).Find(&products).Error
	return
}

type CategoryDao struct {
	*gorm.DB
}

func NewCategoryDao(ctx context.Context) *CategoryDao {
	return &CategoryDao{DB.WithContext(ctx)}
}

// ListCategory 分类列表
func (dao *CategoryDao) ListCategory() (category []*model.Category, err error) {
	err = dao.DB.Model(&model.Category{}).Find(&category).Error
	return
}
