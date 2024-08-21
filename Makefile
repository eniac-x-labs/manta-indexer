GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

DM_ABI_ARTIFACT := ./abis/DelegationManager.sol/DelegationManager.json
RM_ABI_ARTIFACT := ./abis/RewardManager.sol/RewardManager.json
SM_ABI_ARTIFACT := ./abis/StrategyManager.sol/StrategyManager.json

manta-index:
	env GO111MODULE=on go build -v $(LDFLAGS) ./cmd/manta-index

clean:
	rm manta-index

test:
	go test -v ./...

lint:
	golangci-lint run ./...

bindings: binding-dm binding-rm binding-sm

binding-dm:
	 $(eval temp := $(shell mktemp))
	 cat $(DM_ABI_ARTIFACT) | jq -r .bytecode > $(temp)
	 cat $(DM_ABI_ARTIFACT) | jq .abi \
	 | abigen --pkg dm \
	 --abi - \
	 --out bindings/dm/delegation_manager.go \
	 --type DelegationManager \
	 --bin $(temp)
	 rm $(temp)

binding-rm:
	 $(eval temp := $(shell mktemp))
	 cat $(RM_ABI_ARTIFACT) | jq -r .bytecode > $(temp)
	 cat $(RM_ABI_ARTIFACT) | jq .abi \
	 | abigen --pkg rm \
	 --abi - \
	 --out bindings/rm/reward_manager.go \
	 --type RewardManager \
	 --bin $(temp)
	 rm $(temp)

binding-sm:
	 $(eval temp := $(shell mktemp))
	 cat $(SM_ABI_ARTIFACT) | jq -r .bytecode > $(temp)
	 cat $(SM_ABI_ARTIFACT) | jq .abi \
	 | abigen --pkg sm \
	 --abi - \
	 --out bindings/sm/strategy_manager.go \
	 --type StrategyManager \
	 --bin $(temp)
	 rm $(temp)

.PHONY: \
	 event-watcher \
	 bindings \
	 clean \
	 test \
	 lint