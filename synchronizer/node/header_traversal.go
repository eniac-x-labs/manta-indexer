package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/eniac-x-labs/manta-indexer/common/bigint"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

var (
	ErrHeaderTraversalAheadOfProvider            = errors.New("the HeaderTraversal's internal state is ahead of the provider")
	ErrHeaderTraversalAndProviderMismatchedState = errors.New("the HeaderTraversal and provider have diverged in state")
	// ErrHeaderTraversalCheckHeaderByHashDelDbData Delete the previous batch of data
	ErrHeaderTraversalCheckHeaderByHashDelDbData = errors.New("the HeaderTraversal headerList[0].ParentHash != dbLatestHeader.Hash()")
)

type HeaderTraversal struct {
	ethClient EthClient
	chainId   uint

	latestHeader        *types.Header
	lastTraversedHeader *types.Header

	blockConfirmationDepth *big.Int
}

func NewHeaderTraversal(ethClient EthClient, fromHeader *types.Header, confDepth *big.Int, chainId uint) *HeaderTraversal {
	return &HeaderTraversal{
		ethClient:              ethClient,
		lastTraversedHeader:    fromHeader,
		blockConfirmationDepth: confDepth,
		chainId:                chainId,
	}
}

func (f *HeaderTraversal) LatestHeader() *types.Header {
	return f.latestHeader
}

func (f *HeaderTraversal) LastTraversedHeader() *types.Header {
	return f.lastTraversedHeader
}

func (f *HeaderTraversal) NextHeaders(maxSize uint64) ([]types.Header, error) {
	latestHeader, err := f.ethClient.BlockHeaderByNumber(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to query latest block: %w", err)
	} else if latestHeader == nil {
		return nil, fmt.Errorf("latest header unreported")
	} else {
		f.latestHeader = latestHeader
	}
	latestHeaderJson, _ := json.Marshal(latestHeader)
	log.Info("header_traversal dbLatestHeader: ", "info", string(latestHeaderJson))

	endHeight := new(big.Int).Sub(latestHeader.Number, f.blockConfirmationDepth)
	if endHeight.Sign() < 0 {
		// No blocks with the provided confirmation depth available
		return nil, nil
	}

	lastTraversedHeaderJson, _ := json.Marshal(f.lastTraversedHeader)
	log.Info("header_traversal lastTraversedHeaderJson: ", "info", string(lastTraversedHeaderJson))

	if f.lastTraversedHeader != nil {
		cmp := f.lastTraversedHeader.Number.Cmp(endHeight)
		if cmp == 0 {
			return nil, nil
		} else if cmp > 0 {
			return nil, ErrHeaderTraversalAheadOfProvider
		}
	}

	nextHeight := bigint.Zero
	if f.lastTraversedHeader != nil {
		nextHeight = new(big.Int).Add(f.lastTraversedHeader.Number, bigint.One)
	}

	endHeight = bigint.Clamp(nextHeight, endHeight, maxSize)
	headers, err := f.ethClient.BlockHeadersByRange(nextHeight, endHeight, f.chainId)
	if err != nil {
		return nil, fmt.Errorf("error querying blocks by range: %w", err)
	}
	if len(headers) == 0 {
		return nil, nil
	}
	err = f.checkHeaderListByHash(f.lastTraversedHeader, headers)
	if err != nil {
		log.Error("NextHeaders checkBlockListByHash", "error", err)
		return nil, err
	}

	numHeaders := len(headers)
	if numHeaders == 0 {
		return nil, nil
	} else if f.lastTraversedHeader != nil && headers[0].ParentHash != f.lastTraversedHeader.Hash() {
		fmt.Println(f.lastTraversedHeader.Number)
		fmt.Println(headers[0].Number)
		fmt.Println(len(headers))
		log.Error("ErrHeaderTraversalAndProviderMismatchedState", "headers[0].ParentHash = ", headers[0].ParentHash.String(), "lastTraversedHeader.Hash", f.lastTraversedHeader.Hash().String())
		return nil, ErrHeaderTraversalAndProviderMismatchedState
	}
	f.lastTraversedHeader = &headers[numHeaders-1]
	return headers, nil
}

func (f *HeaderTraversal) checkHeaderListByHash(dbLatestHeader *types.Header, headerList []types.Header) error {
	if len(headerList) == 0 {
		return nil
	}
	if len(headerList) == 1 {
		return nil
	}

	dbLatestHeaderJson, _ := json.Marshal(dbLatestHeader)
	log.Info("header_traversal checkHeaderListByHash dbLatestHeaderJson ", "info", string(dbLatestHeaderJson))
	headerFirstJson, _ := json.Marshal(headerList[0])
	log.Info("header_traversal checkHeaderListByHash headerFirstJson ", "info", string(headerFirstJson))

	// check input and db
	// input first ParentHash = dbLatestHeader.Hash
	if dbLatestHeader != nil && headerList[0].ParentHash != dbLatestHeader.Hash() {
		log.Error("checkHeaderListByHash", "headers[0].ParentHash = ", headerList[0].ParentHash.String(), "dbLatestHeader.Hash()", dbLatestHeader.Hash().String())
		return ErrHeaderTraversalCheckHeaderByHashDelDbData
	}

	// check input
	for i := 1; i < len(headerList); i++ {
		if headerList[i].ParentHash != headerList[i-1].Hash() {
			return fmt.Errorf("checkHeaderListByHash: headerList headerList[i].ParentHash != headerList[i-1].Hash()")
		}
	}
	return nil
}

func (f *HeaderTraversal) ChangeLastTraversedHeaderByDelAfter(dbLatestHeader *types.Header) {
	f.lastTraversedHeader = dbLatestHeader
}
