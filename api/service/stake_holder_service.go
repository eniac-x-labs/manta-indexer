package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/worker"
	"strings"
)

func (h HandlerSvc) GetStakeHolder(staker string) (*worker.StakeHolder, error) {
	stakerAddress := strings.ToLower(staker)
	stakeHolder, err := h.stakeHolderView.GetStakeHolder(stakerAddress)
	if err != nil {
		return &worker.StakeHolder{}, err
	}
	return stakeHolder, err
}

func (h HandlerSvc) ListStakeHolder(params *models.QueryDTParams) (*models.StakeHolderListResponse, error) {
	stakeHolderList, total := h.stakeHolderView.ListStakeHolder(params.Page, params.PageSize, params.Order)
	return &models.StakeHolderListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: stakeHolderList,
	}, nil
}
