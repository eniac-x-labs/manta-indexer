package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
	"strconv"
	"strings"
)

type Service interface {
	RegisterOperatorList(*models.QueryDTParams) (*models.RegisterOperatorListResponse, error)
	RegisterOperator(operator string) (*event.OperatorRegistered, error)

	QueryDTListParams(page string, pageSize string, order string) (*models.QueryDTParams, error)
}

type HandlerSvc struct {
	v                      *Validator
	operatorRegisteredView event.OperatorRegisteredView
}

func New(v *Validator, rgv event.OperatorRegisteredView) Service {
	return &HandlerSvc{
		v:                      v,
		operatorRegisteredView: rgv,
	}
}

func (h HandlerSvc) RegisterOperator(operator string) (*event.OperatorRegistered, error) {
	addressToLower := strings.ToLower(operator)
	operatorRegistered, err := h.operatorRegisteredView.QueryOperatorRegistered(addressToLower)
	if err != nil {
		return &event.OperatorRegistered{}, err
	}
	return operatorRegistered, err
}

func (h HandlerSvc) RegisterOperatorList(params *models.QueryDTParams) (*models.RegisterOperatorListResponse, error) {
	operatorRegisteredList, total := h.operatorRegisteredView.QueryOperatorRegisteredList(params.Page, params.PageSize, params.Order)
	return &models.RegisterOperatorListResponse{
		Current: params.Page,
		Size:    params.PageSize,
		Total:   total,
		Records: operatorRegisteredList,
	}, nil
}

func (h HandlerSvc) QueryDTListParams(page string, pageSize string, order string) (*models.QueryDTParams, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	pageVal := h.v.ValidatePage(pageInt)

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, err
	}
	pageSizeVal := h.v.ValidatePageSize(pageSizeInt)
	orderBy := h.v.ValidateOrder(order)

	return &models.QueryDTParams{
		Page:     pageVal,
		PageSize: pageSizeVal,
		Order:    orderBy,
	}, nil
}
