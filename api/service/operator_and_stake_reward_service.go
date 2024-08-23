package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

func (h HandlerSvc) GetOperatorAndStakeReward(guid string) (*event.OperatorAndStakeReward, error) {
	reward, err := h.operatorAndStakeRewardView.GetOperatorAndStakeReward(guid)
	if err != nil {
		return &event.OperatorAndStakeReward{}, err
	}
	return reward, err
}

func (h HandlerSvc) ListOperatorAndStakeReward(params *models.QueryDTParams) (*models.OperatorAndStakeRewardListResponse, error) {
	rewardList, total := h.operatorAndStakeRewardView.ListOperatorAndStakeReward(params.Page, params.PageSize, params.Order)
	return &models.OperatorAndStakeRewardListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: rewardList,
	}, nil
}
