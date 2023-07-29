package service

import (
	"cloudrestaurant/dao"
	"cloudrestaurant/ero"
	"cloudrestaurant/model"
	"cloudrestaurant/serializer"
	"cloudrestaurant/tool"
	"context"
	"errors"
	"fmt"
	logging "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

type MoneyService struct {
	Key   string `json:"key" form:"key"`
	Money string `json:"money" form:"money"`
}

func (service *MoneyService) Add(ctx context.Context, uId uint) serializer.Response {
	code := ero.SUCCESS

	err := dao.NewMemberDao(ctx).Transaction(
		func(tx *gorm.DB) error {
			tool.Encrypt.SetKey(service.Key)

			userDao := dao.NewMemberDaoByDB(tx)
			user, err := userDao.QueryById(uId)
			if err != nil {
				logging.Info(err)
				code = ero.ErrorDatabase
				return err
			}

			// 对钱进行解密。加上新增的。再进行加密。
			var moneyStr string
			if user.Money != "" {
				moneyStr = tool.Encrypt.AesDecoding(user.Money)
			} else {
				moneyStr = "0"
			}

			moneyCount, _ := strconv.ParseFloat(moneyStr, 64)
			moneyAdd, _ := strconv.ParseFloat(service.Money, 64)
			finMoney := fmt.Sprintf("%f", moneyCount+moneyAdd)

			user.Money = tool.Encrypt.AesEncoding(finMoney)

			err = userDao.UpdateUserById(uId, user)
			if err != nil { // 更新用户金额失败，回滚
				logging.Info(err)
				code = ero.ErrorDatabase
				return err
			}
			return nil
		})

	if err != nil {
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

func (service *MoneyService) Show(ctx context.Context, uId uint) serializer.Response {
	code := ero.SUCCESS
	md := dao.NewMemberDao(ctx)
	user, err := md.QueryById(uId)
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
		Data:   serializer.BuildMoney(user, service.Key), //对用户的金额进行了对称加密
		Msg:    ero.GetMsg(code),
	}
}

type OrderPayService struct {
	OrderId   uint    `form:"order_id" json:"order_id"`
	Money     float64 `form:"money" json:"money"`
	OrderNo   string  `form:"orderNo" json:"orderNo"`
	ProductID int     `form:"product_id" json:"product_id"`
	PayTime   string  `form:"payTime" json:"payTime" `
	Sign      string  `form:"sign" json:"sign" `
	BossID    int     `form:"boss_id" json:"boss_id"`
	BossName  string  `form:"boss_name" json:"boss_name"`
	Num       int     `form:"num" json:"num"`
	Key       string  `form:"key" json:"key"`
}

func (service *OrderPayService) PayDown(ctx context.Context, uId uint) serializer.Response {
	code := ero.SUCCESS

	err := dao.NewOrderDao(ctx).Transaction(
		func(tx *gorm.DB) error {
			tool.Encrypt.SetKey(service.Key)
			orderDao := dao.NewOrderDaoByDB(tx)

			order, err := orderDao.GetOrderById(service.OrderId)
			if err != nil {
				logging.Info(err)
				return err
			}
			money := order.Money
			num := order.Num
			money = money * float64(num)

			userDao := dao.NewMemberDaoByDB(tx)
			user, err := userDao.QueryById(uId)
			if err != nil {
				logging.Info(err)
				code = ero.ErrorDatabase
				return err
			}

			// 对钱进行解密。减去订单。再进行加密。
			var moneyStr string
			if user.Money != "" {
				moneyStr = tool.Encrypt.AesDecoding(user.Money)
			} else {
				moneyStr = "0"
			}
			moneyFloat, _ := strconv.ParseFloat(moneyStr, 64)
			if moneyFloat-money < 0.0 { // 金额不足进行回滚
				logging.Info(err)
				code = ero.ErrorDatabase
				return errors.New("金币不足")
			}

			finMoney := fmt.Sprintf("%f", moneyFloat-money)
			user.Money = tool.Encrypt.AesEncoding(finMoney)

			err = userDao.UpdateUserById(uId, user)
			if err != nil { // 更新用户金额失败，回滚
				logging.Info(err)
				code = ero.ErrorDatabase
				return err
			}

			boss := new(model.Member)
			boss, err = userDao.QueryById(uint(service.BossID))

			if boss.Money != "" {
				moneyStr = tool.Encrypt.AesDecoding(boss.Money)
			} else {
				moneyStr = "0"
			}
			moneyFloat, _ = strconv.ParseFloat(moneyStr, 64)
			finMoney = fmt.Sprintf("%f", moneyFloat+money)
			boss.Money = tool.Encrypt.AesEncoding(finMoney)

			err = userDao.UpdateUserById(uint(service.BossID), boss)
			if err != nil { // 更新boss金额失败，回滚
				logging.Info(err)
				code = ero.ErrorDatabase
				return err
			}

			product := new(model.Product)
			productDao := dao.NewProductDaoByDB(tx)
			product, err = productDao.GetProductById(uint(service.ProductID))
			if err != nil {
				return err
			}
			product.Num -= num
			err = productDao.UpdateProduct(uint(service.ProductID), product)
			if err != nil { // 更新商品数量减少失败，回滚
				logging.Info(err)
				code = ero.ErrorDatabase
				return err
			}

			// 更新订单状态
			order.Type = 2
			err = orderDao.UpdateOrderById(service.OrderId, order)
			if err != nil { // 更新订单失败，回滚
				logging.Info(err)
				code = ero.ErrorDatabase
				return err
			}

			productUser := model.Product{
				Name:          product.Name,
				CategoryID:    product.CategoryID,
				Title:         product.Title,
				Info:          product.Info,
				ImgPath:       product.ImgPath,
				Price:         product.Price,
				DiscountPrice: product.DiscountPrice,
				Num:           num,
				OnSale:        false,
				BossID:        uId,
				BossName:      user.UserName,
				BossAvatar:    user.Avatar,
			}

			err = productDao.CreateProduct(&productUser)
			if err != nil { // 买完商品后创建成了自己的商品失败。订单失败，回滚
				logging.Info(err)
				code = ero.ErrorDatabase
				return err
			}
			return nil
		})

	if err != nil {
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
