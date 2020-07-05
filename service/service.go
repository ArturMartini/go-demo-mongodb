package service

import (
	"go-demo-mongodb/canonical"
	"go-demo-mongodb/repository"
)

type Service interface {
	Add(*canonical.Player) error
	Update(*canonical.Player) error
	Get(id string) (canonical.Player, error)
	GetAll(offset int, limit int) ([]canonical.Player, error)
	Delete(id string) error
	Upload(id string, bin []byte) error
	Download(id string) ([]byte, error)
}

type service struct {
	repo repository.Repository
}

var instance Service

func NewService() Service {
	if instance == nil {
		instance = &service{
			repo: repository.NewRepository(),
		}
	}
	return instance
}

func (r service) Add(player *canonical.Player) error {
	return r.repo.Add(player)
}

func (r service)  Update(player *canonical.Player) error {
	return r.repo.Update(player)
}

func (r service) Get(id string) (canonical.Player, error) {
	return r.repo.Get(id)
}

func (r service) GetAll(offset int, limit int) ([]canonical.Player, error) {
	return r.repo.GetAll(offset, limit)
}

func (r service) Delete(id string) error {
	return r.repo.Delete(id)
}

func (r service) Upload(id string, bin []byte) error {
	return r.repo.Upload(id, bin)
}

func (r service) Download(id string) ([]byte, error) {
	return r.repo.Download(id)
}