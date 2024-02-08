package server_streaming

import (
	"context"
	"errors"
	"log"
	svc "server-streaming/service"
	"strings"
)

type userService struct {
	svc.UnimplementedUsersServer
}

type repoServer struct {
	svc.UnimplementedRepoServer
}

func (s *userService) GetUser(ctx context.Context, in *svc.UserGetRequest) (*svc.UserGetReply, error) {
	log.Printf(
		"Received request for user with Email: %s Id %s\n",
		in.Email,
		in.Id,
	)
	componments := strings.Split(in.Email, "@")
	if len(componments) != 2 {
		return nil, errors.New("invalid email address")
	}
	u := &svc.User{
		Id:        in.Id,
		FirstName: componments[0],
		LastName:  componments[1],
		Age:       36,
	}
	return &svc.UserGetReply{User: u}, nil
}
