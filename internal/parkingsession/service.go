package parkingsession

import (
	"context"
	"fmt"
)

type Service struct {
	repo            Repository
	plateProccessor OCRService
	imageStore      ImageStore
}

func NewService(repo Repository, plateProccessor OCRService, imageStore ImageStore) *Service {
	return &Service{repo: repo, plateProccessor: plateProccessor, imageStore: imageStore}
}

type CheckInParams struct {
	CardUID     string
	ImgPlateIn  []byte
	ImgPersonIn []byte
}

type CheckOutParams struct {
	ImgPlateOut  []byte
	ImgPersonOut []byte
}

func (s *Service) CheckIn(ctx context.Context, params CheckInParams) (*ParkingSession, error) {
	ongoing, err := s.repo.GetOngoingSessionByCard(ctx, params.CardUID)
	if err != nil {
		return nil, fmt.Errorf("checking ongoing session: %w", err)
	}
	if ongoing != nil {
		return nil, fmt.Errorf("card %s already has an ongoing session (id: %d)", params.CardUID, ongoing.ID)
	}
	//  Save the two images using the image storage service and get the paths
	ImgPlateInPath, err := s.imageStore.SaveImage(ctx, params.ImgPlateIn)
	if err != nil {
		return nil, fmt.Errorf("saving plate in image: %w", err)
	}
	ImgPersonInPath, err := s.imageStore.SaveImage(ctx, params.ImgPersonIn)
	if err != nil {
		return nil, fmt.Errorf("saving person in image: %w", err)
	}
	//  Get the plate number from the plate image using the OCR service
	plateIn, err := s.plateProccessor.ExtractPlate(ctx, params.ImgPlateIn)
	if err != nil {
		return nil, fmt.Errorf("recognizing plate: %w", err)
	}
	session := &ParkingSession{
		CardUID:         params.CardUID,
		PlateIn:         &plateIn,
		ImgPlateInPath:  &ImgPlateInPath,
		ImgPersonInPath: &ImgPersonInPath,
		Status:          "ongoing",
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("creating parking session: %w", err)
	}

	return session, nil
}

func (s *Service) CheckOut(ctx context.Context, id int64, params CheckOutParams) error {
	session, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("fetching session: %w", err)
	}
	if session.Status != "ongoing" {
		return fmt.Errorf("session %d is already completed", id)
	}

	imgPlateOutPath, err := s.imageStore.SaveImage(ctx, params.ImgPlateOut)
	if err != nil {
		return fmt.Errorf("saving plate out image: %w", err)
	}
	imgPersonOutPath, err := s.imageStore.SaveImage(ctx, params.ImgPersonOut)
	if err != nil {
		return fmt.Errorf("saving person out image: %w", err)
	}

	plateOut, err := s.plateProccessor.ExtractPlate(ctx, params.ImgPlateOut)
	if err != nil {
		return fmt.Errorf("recognizing plate: %w", err)
	}

	session.PlateOut = &plateOut
	session.ImgPlateOutPath = &imgPlateOutPath
	session.ImgPersonOutPath = &imgPersonOutPath

	if err := s.repo.CheckOut(ctx, id, session); err != nil {
		return fmt.Errorf("checking out session: %w", err)
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*ParkingSession, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByCard(ctx context.Context, cardUID string) ([]*ParkingSession, error) {
	return s.repo.ListByCard(ctx, cardUID)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
