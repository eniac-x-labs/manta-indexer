package service

import (
	"math/big"

	"github.com/eniac-x-labs/manta-indexer/database/event/finality"
)

func (h HandlerSvc) GetFinalityVerified(l2BlockNumber *big.Int) (*finality.FinalityVerified, error) {
	strategies, err := h.finalityVerifiedView.GetFinalityVerifiedByL2Block(l2BlockNumber)
	if err != nil {
		return &finality.FinalityVerified{}, err
	}
	return strategies, err
}
