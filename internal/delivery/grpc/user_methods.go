package grpc

import (
	"User/internal/dtos"
	customErrors "User/internal/errors"
	"User/internal/usecase"
	"User/pkg/logger"
	"context"
	"errors"
	userProtobuf "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"time"
)

var _ userProtobuf.UserServer = (*UsergRPC)(nil)

type UsergRPC struct {
	userProtobuf.UnimplementedUserServer

	userUseCase  usecase.UserUseCase
	cloudUseCase usecase.CloudUseCase
}

func Register(
	gRPC *grpc.Server,
	userUseCase usecase.UserUseCase,
	cloudUseCase usecase.CloudUseCase,
) {
	userProtobuf.RegisterUserServer(gRPC, &UsergRPC{userUseCase: userUseCase, cloudUseCase: cloudUseCase})
}

func (u *UsergRPC) CreateUser(
	g grpc.ClientStreamingServer[
		userProtobuf.CreateUserRequest,
		userProtobuf.CreateUserResponse],
) error {

	userData, userPhoto, err := GetUserData(
		&g,
		func(chunk *userProtobuf.CreateUserRequest) interface{} {
			return chunk.GetUserDataForCreate()
		},
		func(chunk *userProtobuf.CreateUserRequest) []byte {
			return chunk.GetUserPhoto()
		},
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid request data")
	}

	if userData == nil {
		logger.ErrorLogger.Printf("user data is empty")
		return status.Error(codes.InvalidArgument, "user data is empty")
	}

	castedUserData, ok := userData.(userProtobuf.UserDataForCreate)
	if !ok {
		logger.ErrorLogger.Printf("user data is not of type UserProtobuf.UserDataForCreate")
		return status.Error(codes.InvalidArgument, "user data is not of type UserProtobuf.UserDataForCreate")
	}

	cmd := &dtos.CreateUserCommand{
		ID:       uuid.New(),
		Name:     castedUserData.Name,
		Email:    castedUserData.Email,
		Role:     castedUserData.Role,
		Password: castedUserData.Password,
	}

	var photoURL string
	if userPhoto != nil {
		url, err := u.cloudUseCase.PutObject(context.TODO(), userPhoto, "user/"+cmd.ID.String())
		photoURL = url
		if err != nil {
			logger.ErrorLogger.Printf("Failed to create user photo in cloud: %v", err)
			return status.Error(codes.Internal, "Failed to create user photo in cloud")
		}
	}

	cmd.Photo = photoURL

	user, err := u.userUseCase.CreateUser(context.TODO(), cmd)
	if err != nil {

		if photoURL == "" {
			err := u.cloudUseCase.DeleteObject(context.TODO(), "user/"+cmd.ID.String())
			if err != nil {
				logger.ErrorLogger.Printf("Failed to delete user photo from cloud: %v", err)
				return status.Error(codes.Internal, "Failed to delete user photo in cloud")
			}
		}

		return status.Error(codes.Internal, "Failed to create user")
	}

	userObject := &userProtobuf.UserObject{
		Id:          user.ID.String(),
		Email:       user.Email,
		Role:        user.Role,
		Photo:       user.Photo,
		Name:        user.Name,
		CreatedTime: user.CreatedTime.String(),
		UpdatedTime: user.UpdatedTime.String(),
	}

	response := &userProtobuf.CreateUserResponse{
		UserObject: userObject,
	}

	err = g.SendAndClose(response)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send user create response: %v", err)
		return status.Error(codes.Internal, "Failed to send user create response")
	}

	return nil
}

func (u *UsergRPC) GetUserById(ctx context.Context,
	request *userProtobuf.GetUserByIdRequest,
) (
	*userProtobuf.GetUserByIdResponse,
	error,
) {
	user, err := u.userUseCase.GetUserById(ctx, uuid.MustParse(request.Id))
	if err != nil {

		if errors.Is(err, customErrors.UserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	userObject := &userProtobuf.UserObject{
		Id:          user.ID.String(),
		Email:       user.Email,
		Name:        user.Name,
		Role:        user.Role,
		Photo:       user.Photo,
		CreatedTime: user.CreatedTime.String(),
		UpdatedTime: user.UpdatedTime.String(),
	}

	response := &userProtobuf.GetUserByIdResponse{
		UserObject: userObject,
	}

	return response, nil
}

func (u *UsergRPC) UpdateUser(
	g grpc.ClientStreamingServer[
		userProtobuf.UpdateUserRequest,
		userProtobuf.UpdateUserResponse],
) error {

	userData, userPhoto, err := GetUserData(
		&g,
		func(chunk *userProtobuf.UpdateUserRequest) interface{} {
			return chunk.GetUserDataForUpdate()
		},
		func(chunk *userProtobuf.UpdateUserRequest) []byte {
			return chunk.GetUserPhoto()
		},
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid request data")
	}

	if userData == nil {
		logger.ErrorLogger.Printf("user data is empty")
		return status.Error(codes.InvalidArgument, "user data is empty")
	}

	castedUserData, ok := userData.(userProtobuf.UserDataForUpdate)
	if !ok {
		logger.ErrorLogger.Printf("user data is not of type UserProtobuf.UserDataForCreate")
		return status.Error(codes.InvalidArgument, "user data is not of type UserProtobuf.UserDataForCreate")
	}

	cmd := &dtos.UpdateUserCommand{
		Id:          uuid.MustParse(castedUserData.Id),
		Name:        castedUserData.Name,
		Email:       castedUserData.Email,
		Role:        castedUserData.Role,
		UpdatedTime: time.Now(),
	}

	_, err = u.userUseCase.GetUserById(context.TODO(), uuid.MustParse(castedUserData.Id))
	if err != nil {
		return status.Error(codes.NotFound, "user not found")
	}

	var photoURL string
	var previousPhoto []byte
	if userPhoto != nil {
		previousPhoto, err = u.cloudUseCase.GetObjectByName(context.TODO(), "user/"+cmd.Id.String())
		if err != nil {
			logger.ErrorLogger.Printf("Failed to get previos photo from cloud: %v", err)
			return err
		}

		url, err := u.cloudUseCase.PutObject(context.TODO(), userPhoto, "user/"+cmd.Id.String())
		photoURL = url
		if err != nil {
			logger.ErrorLogger.Printf("Failed to create user photo in cloud: %v", err)
			return status.Error(codes.Internal, "Failed to create user photo in cloud")
		}
	}

	cmd.Photo = photoURL

	user, err := u.userUseCase.UpdateUser(context.TODO(), cmd)
	if err != nil {

		_, err := u.cloudUseCase.PutObject(context.TODO(), previousPhoto, "user/"+cmd.Id.String())
		if err != nil {
			logger.ErrorLogger.Printf("Failed to set previous photo in cloud: %v", err)
			return status.Error(codes.Internal, "Failed to create user photo in cloud")
		}

		return status.Error(codes.Internal, "Failed to create user")
	}

	userObject := &userProtobuf.UserObject{
		Id:          user.ID.String(),
		Email:       user.Email,
		Role:        user.Role,
		Photo:       user.Photo,
		Name:        user.Name,
		CreatedTime: user.CreatedTime.String(),
		UpdatedTime: user.UpdatedTime.String(),
	}

	updateUserResponse := &userProtobuf.UpdateUserResponse{
		UserObject: userObject,
	}

	err = g.SendAndClose(updateUserResponse)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send user update response: %v", err)
		return err
	}

	return nil
}

func (u *UsergRPC) DeleteUserById(
	ctx context.Context,
	request *userProtobuf.DeleteUserByIdRequest,
) (*userProtobuf.DeleteUserByIdResponse,
	error,
) {
	deletedUser, err := u.userUseCase.DeleteUserById(ctx, uuid.MustParse(request.Id))
	if err != nil {
		return nil, err
	}

	userObject := &userProtobuf.UserObject{
		Id:          deletedUser.ID.String(),
		Email:       deletedUser.Email,
		Role:        deletedUser.Role,
		Photo:       deletedUser.Photo,
		Name:        deletedUser.Name,
		CreatedTime: deletedUser.CreatedTime.String(),
		UpdatedTime: deletedUser.UpdatedTime.String(),
	}

	deleteUserByIdResponse := &userProtobuf.DeleteUserByIdResponse{
		UserObject: userObject,
	}

	return deleteUserByIdResponse, nil
}

func (u *UsergRPC) GetUserByEmail(
	ctx context.Context,
	request *userProtobuf.GetUserByEmailRequest,
) (*userProtobuf.GetUserByEmailResponse,
	error,
) {
	user, err := u.userUseCase.GetUserByEmail(ctx, request.Email)
	if err != nil {

		if errors.Is(err, customErrors.UserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	userObject := &userProtobuf.UserObject{
		Id:          user.ID.String(),
		Email:       user.Email,
		Name:        user.Name,
		Role:        user.Role,
		Photo:       user.Photo,
		CreatedTime: user.CreatedTime.String(),
		UpdatedTime: user.UpdatedTime.String(),
	}

	response := &userProtobuf.GetUserByEmailResponse{
		UserObject: userObject,
	}

	return response, nil
}

func (u *UsergRPC) CheckPassword(
	ctx context.Context,
	request *userProtobuf.CheckPasswordRequest,
) (*emptypb.Empty,
	error,
) {

	checkPasswordCommand := &dtos.CheckPasswordCommand{
		UserId:   uuid.MustParse(request.UserId),
		Password: request.Password,
	}

	err := u.userUseCase.CheckPassword(ctx, checkPasswordCommand)
	if err != nil {

		if errors.Is(err, customErrors.UserNotFound) || errors.Is(err, customErrors.InvalidPassword) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func GetUserData[T any, R any](
	g *grpc.ClientStreamingServer[T, R],
	extractUserData func(chunk *T) interface{},
	extractUserPhoto func(chunk *T) []byte,
) (interface{},
	[]byte,
	error,
) {
	var userData interface{}
	var userPhoto []byte

	for {
		chunk, err := (*g).Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.ErrorLogger.Printf("Error getting chunk: %v", err)
			return nil, nil, err
		}

		if ud := extractUserData(chunk); ud != nil {
			userData = ud
		}

		if uf := extractUserPhoto(chunk); uf != nil {
			userPhoto = append(userPhoto, uf...)
		}
	}

	return userData, userPhoto, nil
}
