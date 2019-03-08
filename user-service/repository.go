package main

/*
repository用来和数据库进行交互
在Repository接口中定义了handler调用的方法。对外只暴露各种服务，然后通过服务来调用内部的函数对数据库进行操作

*/
import (
	pd "github.com/chenxull/shippy/user-service/proto/user"
	"github.com/jinzhu/gorm"
)

type Repository interface {
	GetAll() ([]*pd.User, error)
	Get(id string) (*pd.User, error)
	Create(user *pd.User) error
	GetByEmailAndPassword(user *pd.User) (*pd.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func (repo *UserRepository) GetAll() ([]*pd.User, error) {
	var users []*pd.User
	//TODO
	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *UserRepository) Get(id string) (*pd.User, error) {
	var user *pd.User
	user.Id = id
	if err := repo.db.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) Create(user *pd.User) error {
	if err := repo.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetByEmailAndPassword(user *pd.User) (*pd.User, error) {
	if err := repo.db.First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
