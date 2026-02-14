package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/shnaki/studytrack-api/internal/application/port"
	"github.com/shnaki/studytrack-api/internal/domain"
)

type SubjectUsecase struct {
	subjectRepo port.SubjectRepository
	userRepo    port.UserRepository
}

func NewSubjectUsecase(subjectRepo port.SubjectRepository, userRepo port.UserRepository) *SubjectUsecase {
	return &SubjectUsecase{
		subjectRepo: subjectRepo,
		userRepo:    userRepo,
	}
}

func (u *SubjectUsecase) CreateSubject(ctx context.Context, userID, name string) (*domain.Subject, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	id := uuid.New().String()
	subject, err := domain.NewSubject(id, userID, name)
	if err != nil {
		return nil, err
	}
	if err := u.subjectRepo.Create(ctx, subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (u *SubjectUsecase) ListSubjects(ctx context.Context, userID string) ([]*domain.Subject, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	return u.subjectRepo.FindByUserID(ctx, userID)
}

func (u *SubjectUsecase) UpdateSubject(ctx context.Context, id, name string) (*domain.Subject, error) {
	subject, err := u.subjectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := subject.UpdateName(name); err != nil {
		return nil, err
	}
	if err := u.subjectRepo.Update(ctx, subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (u *SubjectUsecase) DeleteSubject(ctx context.Context, id string) error {
	if _, err := u.subjectRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return u.subjectRepo.Delete(ctx, id)
}
