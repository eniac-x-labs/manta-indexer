package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

func (h HandlerSvc) GetStakerUndelegated(guid string) (*event.StakerUndelegated, error) {
	stakerUndelegated, err := h.stakerUndelegatedView.GetStakerUndelegated(guid)
	if err != nil {
		return &event.StakerUndelegated{}, err
	}
	return stakerUndelegated, err
}

func (h HandlerSvc) ListStakerUndelegated(params *models.QueryDTParams) (*models.StakerUndelegatedListResponse, error) {
	stakerUndelegatedList, total := h.stakerUndelegatedView.ListStakerUndelegated(params.Page, params.PageSize, params.Order)
	return &models.StakerUndelegatedListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: stakerUndelegatedList,
	}, nil
}
