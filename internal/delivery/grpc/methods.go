package grpc

import (
	customError "User/internal/errors"
	"User/internal/models"
	"User/internal/usecase"
	log_c "User/pkg/logger"
	"context"
	userProtobuf "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.user"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
	"time"
)

type UsergRPC struct {
	userProtobuf.UnimplementedUserServer
	userUseCase  usecase.UserUseCase
	cloudUseCase usecase.CloudUseCase
}

func Register(gRPC *grpc.Server, userUseCase usecase.UserUseCase, cloudUseCase usecase.CloudUseCase) {
	userProtobuf.RegisterUserServer(gRPC, &UsergRPC{userUseCase: userUseCase, cloudUseCase: cloudUseCase})
}

func (u UsergRPC) UpdateUser(g grpc.ClientStreamingServer[userProtobuf.UpdateUserRequest, userProtobuf.UpdateUserResponse]) error {

	var userData *userProtobuf.UserData
	var photo []byte

	for {
		chunk, err := g.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if ud := chunk.GetUserData(); ud != nil {
			userData = ud
		}

		if uf := chunk.GetUserPhoto(); uf != nil {
			photo = append(photo, uf...)
		}
	}

	var url string
	if userData != nil {
		photoUrl, err := u.cloudUseCase.PutObject(context.TODO(), photo, userData.Id)
		url = photoUrl
		if err != nil {
			return err
		}
	} else {
		log_c.ErrorLogger.Printf("Feiled put object to cloud: void user data")
		return customError.VoidUserData
	}

	userDataToUpdate := &models.User{
		ID:          uuid.MustParse(userData.Id),
		Name:        userData.Name,
		Email:       userData.Email,
		Role:        userData.Role,
		Photo:       url,
		UpdatedTime: time.Now(),
	}

	user, err := u.userUseCase.UpdateUser(context.TODO(), userDataToUpdate)
	if err != nil {
		return err
	}

	uur := &userProtobuf.UpdateUserResponse{
		Id:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
		Photo: user.Photo,
	}

	err = g.SendAndClose(uur)
	if err != nil {
		log_c.ErrorLogger.Printf("Failed to send user update response: %v", err)
		return err
	}

	return nil
}
