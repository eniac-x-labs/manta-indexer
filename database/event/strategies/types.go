package strategies

import (
	"math/big"
)

type StrategyType struct {
	Strategy string
	Tvl      *big.Int
}
