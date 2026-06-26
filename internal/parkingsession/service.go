package parkingsession

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
	repo            Repository
	plateProccessor OCRService
	imageStore      ImageStore
	payment         PaymentService
}

func NewService(repo Repository, plateProccessor OCRService, imageStore ImageStore, payment PaymentService) *Service {
	return &Service{repo: repo, plateProccessor: plateProccessor, imageStore: imageStore, payment: payment}
}

type CheckInParams struct {
	CardUID     string
	ImgPlateIn  []byte
	ImgPersonIn []byte
	PlateIn     string // optional, used for mock cards
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
	//  Get the plate number
	var plateIn string
	if strings.HasPrefix(params.CardUID, "NFC-MOCK-") {
		plateIn = params.PlateIn
	} else {
		plateIn, err = s.plateProccessor.ExtractPlate(ctx, params.ImgPlateIn)
		if err != nil {
			return nil, fmt.Errorf("recognizing plate: %w", err)
		}

		isCasual, err := s.repo.IsCasualCard(ctx, params.CardUID)
		if err != nil {
			return nil, fmt.Errorf("checking card type: %w", err)
		}
		if !isCasual {
			vehiclePlate, err := s.repo.GetVehiclePlateByCard(ctx, params.CardUID)
			if err != nil {
				return nil, fmt.Errorf("getting vehicle for card: %w", err)
			}
			if vehiclePlate == "" || !plateMatches(vehiclePlate, plateIn) {
				return nil, fmt.Errorf("plate %s does not match registered vehicle %s for this card", plateIn, vehiclePlate)
			}
		}
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

func plateMatches(registered, ocr string) bool {
	// Normalize both plates: remove dashes, spaces, uppercase
	norm := func(s string) string {
		result := ""
		for _, c := range s {
			if c != '-' && c != ' ' {
				if c >= 'a' && c <= 'z' {
					result += string(c - 32)
				} else {
					result += string(c)
				}
			}
		}
		return result
	}
	nr := norm(registered)
	no := norm(ocr)
	return nr == no || (len(no) >= 4 && (strings.Contains(nr, no) || strings.Contains(no, nr)))
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

	var plateOut string
	if strings.HasPrefix(session.CardUID, "NFC-MOCK-") {
		if session.PlateIn != nil {
			plateOut = *session.PlateIn
		}
	} else {
		plateOut, err = s.plateProccessor.ExtractPlate(ctx, params.ImgPlateOut)
		if err != nil {
			return fmt.Errorf("recognizing plate: %w", err)
		}

		isCasual, err := s.repo.IsCasualCard(ctx, session.CardUID)
		if err != nil {
			return fmt.Errorf("checking card type: %w", err)
		}
		if isCasual && session.PlateIn != nil && !plateMatches(*session.PlateIn, plateOut) {
			return fmt.Errorf("plate %s does not match check-in plate %s", plateOut, *session.PlateIn)
		}
	}

	session.PlateOut = &plateOut
	session.ImgPlateOutPath = &imgPlateOutPath
	session.ImgPersonOutPath = &imgPersonOutPath

	const parkingFee = 5000.0
	if s.payment != nil && !strings.HasPrefix(session.CardUID, "NFC-MOCK-") {
		if err := s.payment.ChargeParking(ctx, session.CardUID, parkingFee, id); err != nil {
			return fmt.Errorf("charging parking fee: %w", err)
		}
	}

	if err := s.repo.CheckOut(ctx, id, session); err != nil {
		return fmt.Errorf("checking out session: %w", err)
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*ParkingSession, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) LookupByPlate(ctx context.Context, plate string) (*ParkingSession, error) {
	return s.repo.GetOngoingByPlate(ctx, plate)
}

func (s *Service) ListByCard(ctx context.Context, cardUID string) ([]*ParkingSession, error) {
	return s.repo.ListByCard(ctx, cardUID)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
