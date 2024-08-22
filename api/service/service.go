package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"strconv"
)

type Service interface {
	GetDepositTokensList(*models.QueryDTParams) (*models.DepositTokensResponse, error)

	QueryDTListParams(page string, pageSize string, order string) (*models.QueryDTParams, error)
}

type HandlerSvc struct {
	v *Validator
}

func (h HandlerSvc) GetDepositTokensList(params *models.QueryDTParams) (*models.DepositTokensResponse, error) {
	panic("implement me")
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

func New(v *Validator) Service {
	return &HandlerSvc{
		v: v,
	}
}

func (h HandlerSvc) GetDepositList(params *models.QueryDTParams) (*models.DepositTokensResponse, error) {
	return &models.DepositTokensResponse{
		Current: params.Page,
		Size:    params.PageSize,
		Total:   int64(100),
	}, nil
}
