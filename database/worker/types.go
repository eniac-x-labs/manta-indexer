package worker

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type OperatorsType struct {
	Socket                   string
	EarningsReceiver         common.Address
	DelegationApprover       common.Address
	StakerOptoutWindowBlocks *big.Int
	TotalMantaStake          *big.Int
	TotalStakeReward         *big.Int
	RateReturn               string
	Status                   uint8
}
