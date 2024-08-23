package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

func (h HandlerSvc) GetStakeHolderClaimReward(guid string) (*event.StakeHolderClaimReward, error) {
	stakeHolderClaimReward, err := h.stakeHolderClaimRewardView.GetStakeHolderClaimReward(guid)
	if err != nil {
		return &event.StakeHolderClaimReward{}, err
	}
	return stakeHolderClaimReward, err
}

func (h HandlerSvc) ListStakeHolderClaimReward(params *models.QueryDTParams) (*models.StakeHolderClaimRewardListResponse, error) {
	stakeHolderClaimRewardList, total := h.stakeHolderClaimRewardView.ListStakeHolderClaimReward(params.Page, params.PageSize, params.Order)
	return &models.StakeHolderClaimRewardListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: stakeHolderClaimRewardList,
	}, nil
}
