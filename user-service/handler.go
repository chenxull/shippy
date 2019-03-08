package main

/*
handler.go 文件中实现proto中定义所提供的一些服务，在这里UserService定义了如下rpc服务
service UserService {
    rpc Create (User) returns (Response) {}
    rpc Get (User) returns (Response) {}
    rpc GetAll (Request) returns (Response) {}
    rpc Auth (User) returns (Token) {}
    rpc ValidateToken (Token) returns (Token) {}
} */

import (
	"context"

	pd "github.com/chenxull/shippy/user-service/proto/user"
)

type service struct {
	repo         Repository
	tokenService Authable
}

func (s *service) Get(ctx context.Context, req *pd.User, resp *pd.Response) error {
	user, err := s.repo.Get(req.Id)
	if err != nil {
		return err
	}
	resp.User = user
	return nil
}

func (s *service) GetAll(ctx context.Context, req *pd.Request, resp *pd.Response) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	resp.Users = users
	return nil
}

func (s *service) Auth(ctx context.Context, req *pd.User, resp *pd.Token) error {
	_, err := s.repo.GetByEmailAndPassword(req)
	if err != nil {
		return err
	}
	resp.Token = "testingabc"
	return nil
}

func (s *service) Create(ctx context.Context, req *pd.User, resp *pd.Response) error {
	if err := s.repo.Create(req); err != nil {
		return err
	}
	resp.User = req
	return nil
}

func (s *service) ValidateToken(ctx context.Context, req *pd.Token, resp *pd.Token) error {
	return nil
}
