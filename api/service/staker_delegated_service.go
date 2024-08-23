package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

func (h HandlerSvc) GetStakerDelegated(guid string) (*event.StakerDelegated, error) {
	stakerDelegated, err := h.stakerDelegatedView.GetStakerDelegated(guid)
	if err != nil {
		return &event.StakerDelegated{}, err
	}
	return stakerDelegated, err
}

func (h HandlerSvc) ListStakerDelegated(params *models.QueryDTParams) (*models.StakerDelegatedListResponse, error) {
	stakerDelegatedList, total := h.stakerDelegatedView.ListStakerDelegated(params.Page, params.PageSize, params.Order)
	return &models.StakerDelegatedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: stakerDelegatedList,
	}, nil
}
