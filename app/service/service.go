package service

import (
	"github.com/google/wire"
	"importer/api"
	"importer/app/dao"
)

var Provider = wire.NewSet(New, wire.Bind(new(api.NegtServer), new(*Service)))

// Service 服务层接口
type Service struct {
	dao dao.Dao // 数据层接口
}

func New(d dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		dao: d,
	}
	return
}

func (s *Service) Ping() error {
	return nil
}

func (s *Service) Start() error {
	return s.dao.Start()
}

func (s *Service) Close() {
}

func (s *Service) ImportRps(date string) error {
	return s.dao.ImportRps(date)
}

func (s *Service) ImportFxj() error {
	return s.dao.ImportFxj()
}
