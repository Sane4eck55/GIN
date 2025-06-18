package services

import (
    "context"
    "log"
    "myapp/api"
    "myapp/models"
    "gorm.io/gorm"
)

type GRPCServer struct {
    proto.UnimplementedUserServiceServer
    DB *gorm.DB
}

func (s *GRPCServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
    var user models.User
    if err := s.DB.First(&user, req.Id).Error; err != nil {
        return nil, err
    }
    return &proto.GetUserResponse{
        Id:   int32(user.ID),
        Name: user.Name,
    }, nil
}
