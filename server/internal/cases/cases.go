package cases

import (
	"taskFive/server/internal/entity"
	"taskFive/server/internal/repository"
)

type (
	Usecase interface {
		CreateItem(*entity.Item, int) (*entity.Item, error)
		PatchItem(*entity.Item, int, int) (*entity.Item, error)
		DeleteItem(int, int) (*entity.DelR, error)
		GetAll() ([]*entity.Item, error)
	}

	DataBase interface {
		CreateItem(*entity.Item) (*entity.Item, error)
		PatchItem(*entity.Item, int, int) (*entity.Item, error)
		DeleteItem(int, int) (*entity.DelR, error)
		GetAll() ([]*entity.Item, error)
	}
)

type usecase struct {
	rep *repository.Repository
}

func NewUseCase(rep *repository.Repository) *usecase {
	return &usecase{
		rep: rep,
	}
}

func (u *usecase) CreateItem(item *entity.Item, cId int) (*entity.Item, error) {
	data, err := u.rep.CreateItem(item, cId)
	return data, err
}

func (u *usecase) PatchItem(item *entity.Item, cId, iId int) (*entity.Item, error) {
	data, err := u.rep.PatchItem(item, cId, iId)
	return data, err
}

func (u *usecase) DeleteItem(iId, cId int) (*entity.DelR, error) {
	data, err := u.rep.DeleteItem(iId, cId)
	if err != nil {
		return data, err
	}
	return data, err
}

func (u *usecase) GetAll() ([]*entity.Item, error) {
	data, err := u.rep.GetAll()
	return data, err
}
