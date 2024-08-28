package worker

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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

type StakeStrategyOperatorType struct {
	MantaStake    *big.Int
	Reward        *big.Int
	ClaimedAmount *big.Int
	Timestamp     uint64
}
