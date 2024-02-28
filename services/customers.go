package services

import (
	"autflow_back/models"
	"autflow_back/models/dto"
	"autflow_back/repositories"
	"autflow_back/utils"
	"context"
)

type Customer struct {
	customerRepository *repositories.Customers
	logger             utils.Logger
}

func NewCustomer(workflowRepository *repositories.Workflows,
	customerRepository *repositories.Customers,
	logger utils.Logger) *Customer {
	return &Customer{
		logger:             logger,
		customerRepository: customerRepository,
	}
}

func (r *Customer) Insert(ctx context.Context, dt *dto.CreateCustomerDTO) (string, error) {
	customer := dt.ToCustomer()

	r.logger.Debugf("Create Customer: %+v", customer)

	idCriado, erro := r.customerRepository.Insert(ctx, customer)
	if erro != nil {
		return "", erro
	}

	return idCriado, nil
}

func (r *Customer) Find(ctx context.Context, query string) ([]models.Customer, error) {
	customers, erro := r.customerRepository.Find(ctx, query)
	if erro != nil {
		return nil, erro
	}

	return customers, nil
}

func (r *Customer) FindId(ctx context.Context, id string) (models.Customer, error) {
	customer, erro := r.customerRepository.FindId(ctx, id, false)
	if erro != nil {
		return models.Customer{}, nil

	}

	return *customer, nil
}

func (r *Customer) Edit(ctx context.Context, id string, dt *dto.CreateCustomerDTO) error {
	newCustomer := dt.ToCustomer()

	r.logger.Debugf("Edit Customer: %+v", newCustomer)

	return r.customerRepository.Edit(ctx, id, newCustomer)
}

func (r *Customer) Delete(ctx context.Context, id string) error {
	return r.customerRepository.Delete(ctx, id)
}
