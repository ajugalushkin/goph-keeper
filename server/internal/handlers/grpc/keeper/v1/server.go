package v1

import (
	"context"
	"errors"
	"strconv"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
	"github.com/ajugalushkin/goph-keeper/server/internal/dto/models"

	"github.com/ajugalushkin/goph-keeper/server/internal/services/auth"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

type Keeper interface {
	ListItem(
		ctx context.Context,
		since int64,
	) (list *models.ListItem, err error)
	SaveItem(ctx context.Context, item *models.Item) (serverUpdateAt string, err error)
}

type serverAPI struct {
	v1.UnimplementedKeeperServiceV1Server
	auth   Auth
	keeper Keeper
}

func Register(gRPC *grpc.Server, auth Auth, keeper Keeper) {
	v1.RegisterKeeperServiceV1Server(gRPC, &serverAPI{
		auth:   auth,
		keeper: keeper,
	})
}

func (s *serverAPI) RegisterV1(
	ctx context.Context,
	req *v1.RegisterRequestV1,
) (*v1.RegisterResponseV1, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register new user")
	}

	return &v1.RegisterResponseV1{UserId: user}, nil
}

func (s *serverAPI) LoginV1(
	ctx context.Context,
	req *v1.LoginRequestV1,
) (*v1.LoginResponseV1, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := validator.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")

		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &v1.LoginResponseV1{
		Token: token,
	}, nil
}

func (s *serverAPI) ListItemsV1(
	ctx context.Context,
	req *v1.ListItemRequestV1,
) (*v1.ListItemResponseV1, error) {
	_, err := s.keeper.ListItem(ctx, req.GetSince())
	if err != nil {
		return nil, err
	}

	return &v1.ListItemResponseV1{}, nil
}

func (s *serverAPI) SetItemV1(
	ctx context.Context,
	req *v1.SetItemRequestV1,
) (*v1.SetItemResponseV1, error) {
	_, err := s.keeper.SaveItem(ctx, &models.Item{
		ID:              req.GetItem().GetId(),
		Name:            req.GetItem().GetName(),
		Type:            req.GetItem().GetType().String(),
		Value:           req.GetItem().GetValue(),
		ServerUpdatedAt: strconv.FormatInt(req.GetItem().ServerUpdatedAt, 10),
		IsDeleted:       req.GetItem().IsDeleted,
	})
	if err != nil {
		return nil, err
	}

	return &v1.SetItemResponseV1{
		//ServerUpdatedAt: updatedAt,
	}, nil
}
