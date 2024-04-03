package service

import (
	"context"
	"fmt"
	"lms_back/api/models"
	"lms_back/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type paymentService struct {
	storage storage.IStorage
}

func NewPaymentService(storage storage.IStorage) paymentService {
	return paymentService{
		storage: storage,
	}
}

func (u paymentService) Create(ctx context.Context, payment models.CreatePayment) (resp models.Payment, err error) {
	pKey, err := u.storage.Payment().Create(ctx, payment)
	if err != nil {
		return models.Payment{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if pKey == (models.Payment{}) {
		return models.Payment{}, status.Error(codes.InvalidArgument, "Failed to create payment")
	}

	resp, err = u.storage.Payment().GetByID(ctx, pKey.Id)
	if err != nil {
		return models.Payment{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if payment.Student_id != "" && payment.Price > 0 {
		student, err := u.storage.Student().GetByID(ctx, payment.Student_id)
		if err != nil {
			return models.Payment{}, status.Error(codes.InvalidArgument, err.Error())
		}

		if student == (models.Student{}) {
			return models.Payment{}, status.Error(codes.InvalidArgument, "Student not found")
		}

		student.PaidSum += payment.Price

		_, err = u.storage.Student().Update(ctx, student)
		if err != nil {
			return models.Payment{}, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	return resp, nil
}

func (u paymentService) Update(ctx context.Context, payment models.Payment) (models.Payment, error) {

	pKey, err := u.storage.Payment().Update(ctx, payment)
	if err != nil {
		fmt.Println("ERROR in service layer while updating payment", err.Error())
		return models.Payment{}, err
	}

	return pKey, nil
}

func (u paymentService) GetByID(ctx context.Context, id string) (models.Payment, error) {

	pKey, err := u.storage.Payment().GetByID(ctx, id)
	if err != nil {
		fmt.Println("ERROR in service layer while getbyid payment", err.Error())
		return models.Payment{}, err
	}

	return pKey, nil
}

func (u paymentService) GetAll(ctx context.Context, req models.GetAllPaymentsRequest) (models.GetAllPaymentsResponse, error) {

	pKey, err := u.storage.Payment().GetAll(ctx, req)
	if err != nil {
		fmt.Println("ERROR in service layer while GetAll payment", err.Error())
		return models.GetAllPaymentsResponse{}, err
	}

	return pKey, nil
}

func (u paymentService) Delete(ctx context.Context, id string) error {

	err := u.storage.Payment().Delete(ctx, id)
	if err != nil {
		fmt.Println("ERROR in service layer while deleting payment", err.Error())
		return err
	}

	return nil
}
