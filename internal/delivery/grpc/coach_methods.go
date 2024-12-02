package grpc

/*

import (
	"User/internal/dtos"
	customErrors "User/internal/errors"
	"User/internal/usecase"
	"User/pkg/logger"
	"context"
	"errors"
	coachProtobuf "github.com/DanKo-code/FitnessCenter-Protobuf/gen/FitnessCenter.protobuf.coach"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

var _ coachProtobuf.CoachServer = (*CoachgRPC)(nil)

type CoachgRPC struct {
	coachProtobuf.UnimplementedCoachServer

	coachUseCase usecase.CoachUseCase
	cloudUseCase usecase.CloudUseCase
}

func RegisterCoachServer(gRPC *grpc.Server, coachUseCase usecase.CoachUseCase, cloudUseCase usecase.CloudUseCase) {
	coachProtobuf.RegisterCoachServer(gRPC, &CoachgRPC{coachUseCase: coachUseCase, cloudUseCase: cloudUseCase})
}

func (u *CoachgRPC) UpdateCoach(
	g grpc.ClientStreamingServer[coachProtobuf.UpdateCoachRequest, coachProtobuf.UpdateCoachResponse],
) error {
	var coachData *coachProtobuf.CoachData
	var photo []byte

	for {
		chunk, err := g.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if ud := chunk.GetCoachData(); ud != nil {
			coachData = ud
		}

		if uf := chunk.GetCoachPhoto(); uf != nil {
			photo = append(photo, uf...)
		}
	}

	var url string
	if coachData != nil {
		photoUrl, err := u.cloudUseCase.PutObject(context.TODO(), photo, coachData.Id)
		url = photoUrl
		if err != nil {
			return err
		}
	} else {
		logger.ErrorLogger.Printf("Feiled put object to cloud: void coach data")
		return customErrors.VoidCoachData
	}

	coachDataToUpdate := &dtos.UpdateCoachCommand{
		Id:          uuid.MustParse(coachData.Id),
		Name:        coachData.Name,
		Description: coachData.Description,
		Photo:       url,
	}

	coach, err := u.coachUseCase.UpdateCoach(context.TODO(), coachDataToUpdate)
	if err != nil {
		return err
	}

	coachData = &coachProtobuf.CoachData{
		Id:          coach.Id.String(),
		Name:        coach.Name,
		Description: coach.Description,
	}

	response := &coachProtobuf.UpdateCoachResponse{
		CoachData: coachData,
		Url:       coach.Photo,
	}

	err = g.SendAndClose(response)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send coach update response: %v", err)
		return err
	}

	return nil
}

func (u *CoachgRPC) CreateCoach(
	ctx context.Context,
	request *coachProtobuf.CreateCoachRequest,
) (*coachProtobuf.CreateCoachResponse, error) {
	cmd := &dtos.CreateCoachCommand{
		Name:        request.Name,
		Description: request.Description,
	}

	coach, err := u.coachUseCase.CreateCoach(ctx, cmd)
	if err != nil {
		return nil, err
	}

	response := &coachProtobuf.CreateCoachResponse{
		Id:    coach.Id.String(),
		Email: coach.,
		Name:  coach.Name,
		Role:  coach.Role,
		Photo: coach.Photo,
	}

	return response, nil
}

func (u *CoachgRPC) GetCoachById(
	ctx context.Context,
	request *coachProtobuf.GetCoachByIdRequest,
) (*coachProtobuf.GetCoachByIdResponse, error) {
	coach, err := u.coachUseCase.GetCoachById(ctx, uuid.MustParse(request.Id))
	if err != nil {

		if errors.Is(err, customErrors.CoachNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	response := &coachProtobuf.get{
		Id:           coach.ID.String(),
		Email:        coach.Email,
		Name:         coach.Name,
		Role:         coach.Role,
		Photo:        coach.Photo,
		PasswordHash: coach.PasswordHash,
		CreatedTime:  coach.CreatedTime.String(),
		UpdatedTime:  coach.UpdatedTime.String(),
	}

	return response, nil
}
*/
