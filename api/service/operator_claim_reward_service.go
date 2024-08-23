package service

import (
	"github.com/eniac-x-labs/manta-indexer/api/models"
	"github.com/eniac-x-labs/manta-indexer/database/event"
)

func (h HandlerSvc) GetOperatorClaimReward(guid string) (*event.OperatorClaimReward, error) {
	reward, err := h.operatorClaimRewardView.GetOperatorClaimReward(guid)
	if err != nil {
		return &event.OperatorClaimReward{}, err
	}
	return reward, err
}

func (h HandlerSvc) ListOperatorClaimReward(params *models.QueryDTParams) (*models.OperatorClaimRewardListResponse, error) {
	rewardList, total := h.operatorClaimRewardView.ListOperatorClaimReward(params.Page, params.PageSize, params.Order)
	return &models.OperatorClaimRewardListResponse{
		ListResponse: models.ListResponse{
			Current: params.Page,
			Size:    params.PageSize,
			Total:   total,
		},
		Records: rewardList,
	}, nil
}
