package serializer

import (
	"cloudrestaurant/model"
	"cloudrestaurant/tool"
)

type Money struct {
	UserID    uint   `json:"user_id" form:"user_id"`
	UserName  string `json:"user_name" form:"user_name"`
	UserMoney string `json:"user_money" form:"user_money"`
}

func BuildMoney(item *model.Member, key string) Money {
	if item.Money == "" {
		return Money{
			UserID:    item.Id,
			UserName:  item.UserName,
			UserMoney: "0",
		}
	}
	tool.Encrypt.SetKey(key)
	return Money{
		UserID:    item.Id,
		UserName:  item.UserName,
		UserMoney: tool.Encrypt.AesDecoding(item.Money),
	}
}
