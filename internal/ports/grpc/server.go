package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/app"
)

type Server struct {
	a app.App
	UnimplementedAdServiceServer
}

func (s Server) CreateAd(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	_, err := s.a.GetUser(ctx, request.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	ad, err := s.a.CreateAd(ctx, request.Title, request.Text, request.UserId)
	if err != nil {
		if errors.Is(err, app.ErrValidationFail) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	newAd := AdResponse{Id: ad.ID, Title: ad.Title, Text: ad.Text, AuthorId: ad.AuthorID, Published: false}
	return &newAd, nil
}

func (s Server) ChangeAdStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	_, err := s.a.GetUser(ctx, request.UserId)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	ad, err := s.a.ChangeAdStatus(ctx, request.AdId, request.UserId, request.Published)
	if err != nil {
		if errors.Is(err, app.ErrValidationFail) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, app.ErrWrongUser) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
	}
	newAd := AdResponse{Id: ad.ID, Title: ad.Title, Text: ad.Text, AuthorId: ad.AuthorID, Published: ad.Published}
	return &newAd, nil
}

func (s Server) UpdateAd(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	_, err := s.a.GetUser(ctx, request.UserId)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	ad, err := s.a.UpdateAd(ctx, request.AdId, request.UserId, request.Title, request.Text)
	if err != nil {
		if errors.Is(err, app.ErrValidationFail) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, app.ErrWrongUser) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
	}
	newAd := AdResponse{Id: ad.ID, Title: ad.Title, Text: ad.Text, AuthorId: ad.AuthorID, Published: ad.Published}
	return &newAd, nil
}

func (s Server) ListAds(ctx context.Context, empty *emptypb.Empty) (*ListAdResponse, error) {
	ads, err := s.a.GetAds(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	adsList := make([]*AdResponse, 0)
	for _, Ad := range ads {
		ad := &AdResponse{Id: Ad.ID, Title: Ad.Title, Text: Ad.Text, AuthorId: Ad.AuthorID, Published: Ad.Published}
		adsList = append(adsList, ad)
	}
	return &ListAdResponse{List: adsList}, nil
}

func (s Server) CreateUser(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	user, err := s.a.CreateUser(ctx, request.Nickname, request.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	userReq := &UserResponse{Id: user.ID, Nickname: user.Nickname, Email: user.Email}
	return userReq, nil
}

func (s Server) GetUser(ctx context.Context, request *GetUserRequest) (*UserResponse, error) {
	user, err := s.a.GetUser(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	User := UserResponse{Id: user.ID, Nickname: user.Nickname, Email: user.Email}
	return &User, nil
}

func (s Server) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.a.DeleteUser(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s Server) DeleteAd(ctx context.Context, request *DeleteAdRequest) (*emptypb.Empty, error) {
	_, err := s.a.GetUser(ctx, request.AdId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = s.a.DeleteAd(ctx, request.AdId, request.AuthorId)
	if err != nil {
		if errors.Is(err, app.ErrWrongUser) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func NewService(a app.App) AdServiceServer {
	return &Server{a: a}
}
