package service

import (
	"cloudrestaurant/dao"
	"cloudrestaurant/e"
	"cloudrestaurant/model"
	"cloudrestaurant/serializer"
	"context"
	logging "github.com/sirupsen/logrus"
)

type FavoritesService struct {
	ProductId  uint `form:"product_id" json:"product_id"`
	BossId     uint `form:"boss_id" json:"boss_id"`
	FavoriteId uint `form:"favorite_id" json:"favorite_id"`
	PageNum    int  `form:"pageNum"`
	PageSize   int  `form:"pageSize"`
}

// Show 商品收藏夹
func (fs *FavoritesService) Show(ctx context.Context, uId uint) serializer.Response {
	fDao := dao.NewFavoritesDao(ctx)
	code := e.SUCCESS
	if fs.PageSize == 0 {
		fs.PageSize = 15
	}
	favorites, total, err := fDao.ListFavoriteByUserId(uId, fs.PageSize, fs.PageNum)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildFavorites(ctx, favorites), uint(total))
}

// Create 创建收藏夹
func (fs *FavoritesService) Create(ctx context.Context, uId uint) serializer.Response {
	code := e.SUCCESS
	fd := dao.NewFavoritesDao(ctx)
	exist, _ := fd.FavoriteExistOrNot(fs.ProductId, uId)
	if exist {
		code = e.ErrorExistFavorite
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	md := dao.NewMemberDao(ctx)
	user, err := md.QueryById(uId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	boss, err := md.QueryById(fs.BossId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	pd := dao.NewProductDao(ctx)
	product, err := pd.GetProductById(fs.ProductId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	favorite := &model.Favorite{
		MemberID:  uId,
		Member:    *user,
		ProductID: fs.ProductId,
		Product:   *product,
		BossID:    fs.BossId,
		Boss:      *boss,
	}
	fd = dao.NewFavoritesDaoByDB(fd.DB)
	err = fd.CreateFavorite(favorite)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// Delete 删除收藏夹
func (fs *FavoritesService) Delete(ctx context.Context) serializer.Response {
	code := e.SUCCESS

	favoriteDao := dao.NewFavoritesDao(ctx)
	err := favoriteDao.DeleteFavoriteById(fs.FavoriteId)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Data:   e.GetMsg(code),
	}
}
