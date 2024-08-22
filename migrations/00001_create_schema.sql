DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'uint256') THEN
        CREATE DOMAIN UINT256 AS NUMERIC
            CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
    ELSE
        ALTER DOMAIN UINT256 DROP CONSTRAINT uint256_check;
        ALTER DOMAIN UINT256 ADD
            CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
    END IF;
END $$;


CREATE TABLE IF NOT EXISTS block_headers (
    guid        VARCHAR PRIMARY KEY,
    hash        VARCHAR NOT NULL UNIQUE,
    parent_hash VARCHAR NOT NULL UNIQUE,
    number      UINT256 NOT NULL UNIQUE,
    timestamp   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0),
    rlp_bytes   VARCHAR NOT NULL
);
CREATE INDEX IF NOT EXISTS block_headers_timestamp ON block_headers(timestamp);
CREATE INDEX IF NOT EXISTS block_headers_number ON block_headers(number);


CREATE TABLE IF NOT EXISTS contract_events (
    guid             VARCHAR PRIMARY KEY,
    block_hash       VARCHAR NOT NULL REFERENCES block_headers(hash) ON DELETE CASCADE,
    contract_address VARCHAR NOT NULL,
    transaction_hash VARCHAR NOT NULL,
    log_index        INTEGER NOT NULL,
    event_signature  VARCHAR NOT NULL,
    timestamp        INTEGER NOT NULL CHECK (timestamp > 0),
    rlp_bytes        VARCHAR NOT NULL
);
CREATE INDEX IF NOT EXISTS contract_events_timestamp ON contract_events(timestamp);
CREATE INDEX IF NOT EXISTS contract_events_block_hash ON contract_events(block_hash);
CREATE INDEX IF NOT EXISTS contract_events_event_signature ON contract_events(event_signature);
CREATE INDEX IF NOT EXISTS contract_events_contract_address ON contract_events(contract_address);


CREATE TABLE IF NOT EXISTS event_blocks(
    guid        VARCHAR PRIMARY KEY,
    hash        VARCHAR NOT NULL UNIQUE,
    parent_hash VARCHAR NOT NULL UNIQUE,
    number      UINT256 NOT NULL UNIQUE,
    timestamp   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS event_blocks_timestamp ON event_blocks(timestamp);
CREATE INDEX IF NOT EXISTS event_blocks_number ON event_blocks(number);


CREATE TABLE IF NOT EXISTS operator_registered (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    operator                    VARCHAR NOT NULL,
    earnings_receiver           VARCHAR NOT NULL,
    delegation_approver         VARCHAR NOT NULL,
    staker_optout_window_blocks UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_registered_timestamp ON operator_registered(timestamp);
CREATE INDEX IF NOT EXISTS operator_registered_number ON operator_registered(number);
CREATE INDEX IF NOT EXISTS operator_registered_operator ON operator_registered(operator);


CREATE TABLE IF NOT EXISTS operator_modified (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    operator                    VARCHAR NOT NULL,
    earnings_receiver           VARCHAR NOT NULL,
    delegation_approver         VARCHAR NOT NULL,
    staker_optout_window_blocks UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_modified_timestamp ON operator_modified(timestamp);
CREATE INDEX IF NOT EXISTS operator_modified_number ON operator_modified(number);
CREATE INDEX IF NOT EXISTS operator_modified_operator ON operator_modified(operator);


CREATE TABLE IF NOT EXISTS operator_node_url_update (
     guid                        VARCHAR PRIMARY KEY,
     block_hash                  VARCHAR NOT NULL,
     number                      UINT256 NOT NULL,
     tx_hash                     VARCHAR NOT NULL UNIQUE,
     operator                    VARCHAR NOT NULL,
     metadata_uri                VARCHAR NOT NULL,
     timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_node_url_update_timestamp ON operator_node_url_update(timestamp);
CREATE INDEX IF NOT EXISTS operator_node_url_update_number ON operator_node_url_update(number);
CREATE INDEX IF NOT EXISTS operator_node_url_update_operator ON operator_node_url_update(operator);


CREATE TABLE IF NOT EXISTS operator_shares_increased (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    operator                    VARCHAR NOT NULL,
    staker                      VARCHAR NOT NULL,
    strategy                    VARCHAR NOT NULL,
    shares                      UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_shares_increased_timestamp ON operator_shares_increased(timestamp);
CREATE INDEX IF NOT EXISTS operator_shares_increased_number ON operator_shares_increased(number);
CREATE INDEX IF NOT EXISTS operator_shares_increased_operator ON operator_shares_increased(operator);


CREATE TABLE IF NOT EXISTS operator_shares_decreased (
     guid                        VARCHAR PRIMARY KEY,
     block_hash                  VARCHAR NOT NULL,
     number                      UINT256 NOT NULL,
     tx_hash                     VARCHAR NOT NULL UNIQUE,
     operator                    VARCHAR NOT NULL,
     staker                      VARCHAR NOT NULL,
     strategy                    VARCHAR NOT NULL,
     shares                      UINT256 NOT NULL,
     timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_shares_decreased_timestamp ON operator_shares_decreased(timestamp);
CREATE INDEX IF NOT EXISTS operator_shares_decreased_number ON operator_shares_decreased(number);
CREATE INDEX IF NOT EXISTS operator_shares_decreased_operator ON operator_shares_decreased(operator);

CREATE TABLE IF NOT EXISTS staker_delegated (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    operator                    VARCHAR NOT NULL,
    staker                      VARCHAR NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS staker_delegated_timestamp ON staker_delegated(timestamp);
CREATE INDEX IF NOT EXISTS staker_delegated_number ON staker_delegated(number);
CREATE INDEX IF NOT EXISTS staker_delegated_operator ON staker_delegated(operator);

CREATE TABLE IF NOT EXISTS staker_undelegated (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    operator                    VARCHAR NOT NULL,
    staker                      VARCHAR NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS staker_undelegated_timestamp ON staker_undelegated(timestamp);
CREATE INDEX IF NOT EXISTS staker_undelegated_number ON staker_undelegated(number);
CREATE INDEX IF NOT EXISTS staker_undelegated_operator ON staker_undelegated(operator);

CREATE TABLE IF NOT EXISTS withdrawal_queued (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    withdrawal_root             VARCHAR NOT NULL,
    staker                      VARCHAR NOT NULL,
    delegated_to                VARCHAR NOT NULL,
    withdrawer                  VARCHAR NOT NULL,
    nonce                       UINT256 NOT NULL,
    start_block                 UINT256 NOT NULL,
    strategies                  VARCHAR NOT NULL,
    shares                      VARCHAR NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS withdrawal_queued_number ON withdrawal_queued(number);
CREATE INDEX IF NOT EXISTS withdrawal_queued_staker ON withdrawal_queued(staker);

CREATE TABLE IF NOT EXISTS withdrawal_migrated (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    old_withdrawal_root         VARCHAR NOT NULL,
    new_withdrawal_root         VARCHAR NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS withdrawal_migrated_timestamp ON withdrawal_migrated(timestamp);
CREATE INDEX IF NOT EXISTS withdrawal_migrated_number ON withdrawal_migrated(number);
CREATE INDEX IF NOT EXISTS withdrawal_completed_old_withdrawal_root ON withdrawal_migrated(old_withdrawal_root);


CREATE TABLE IF NOT EXISTS min_withdrawal_delay_blocks_set (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    previous_value              UINT256 NOT NULL,
    new_value                   UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS min_withdrawal_delay_blocks_set_timestamp ON min_withdrawal_delay_blocks_set(timestamp);
CREATE INDEX IF NOT EXISTS min_withdrawal_delay_blocks_set_number ON min_withdrawal_delay_blocks_set(number);


CREATE TABLE IF NOT EXISTS strategy_withdrawal_delay_blocks_set (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    strategy                    VARCHAR NOT NULL,
    previous_value              UINT256 NOT NULL,
    new_value                   UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS strategy_withdrawal_delay_blocks_set_timestamp ON strategy_withdrawal_delay_blocks_set(timestamp);
CREATE INDEX IF NOT EXISTS strategy_withdrawal_delay_blocks_set_number ON strategy_withdrawal_delay_blocks_set(number);
CREATE INDEX IF NOT EXISTS strategy_withdrawal_delay_blocks_set_strategy ON strategy_withdrawal_delay_blocks_set(strategy);


CREATE TABLE IF NOT EXISTS strategy_deposit (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    staker                      VARCHAR NOT NULL,
    manta_token                 VARCHAR NOT NULL,
    strategy                    VARCHAR NOT NULL,
    shares                      UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS strategy_deposit_timestamp ON strategy_deposit(timestamp);
CREATE INDEX IF NOT EXISTS strategy_deposit_number ON strategy_deposit(number);
CREATE INDEX IF NOT EXISTS strategy_deposit_strategy ON strategy_deposit(strategy);


CREATE TABLE IF NOT EXISTS operator_and_stake_reward (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    strategy                    VARCHAR NOT NULL,
    operator                    VARCHAR NOT NULL,
    staker_fee                  UINT256 NOT NULL,
    operator_fee                UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_and_stake_reward_timestamp ON operator_and_stake_reward(timestamp);
CREATE INDEX IF NOT EXISTS operator_and_stake_reward_number ON operator_and_stake_reward(number);
CREATE INDEX IF NOT EXISTS operator_and_stake_reward_strategy ON operator_and_stake_reward(strategy);


CREATE TABLE IF NOT EXISTS operator_claim_reward (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    operator                    VARCHAR NOT NULL,
    amount                      UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_claim_reward_timestamp ON operator_claim_reward(timestamp);
CREATE INDEX IF NOT EXISTS operator_claim_reward_number ON operator_claim_reward(number);
CREATE INDEX IF NOT EXISTS operator_claim_reward_operator ON operator_claim_reward(operator);


CREATE TABLE IF NOT EXISTS stake_holder_claim_reward (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    stake_holder                VARCHAR NOT NULL,
    strategy                    VARCHAR NOT NULL,
    amount                      UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS stake_holder_claim_reward_timestamp ON stake_holder_claim_reward(timestamp);
CREATE INDEX IF NOT EXISTS stake_holder_claim_reward_number ON stake_holder_claim_reward(number);
CREATE INDEX IF NOT EXISTS stake_holder_claim_reward_stake_holder ON stake_holder_claim_reward(stake_holder);


CREATE TABLE IF NOT EXISTS operator (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    operator                    VARCHAR NOT NULL,
    socket                      VARCHAR NOT NULL,
    earnings_receiver           VARCHAR NOT NULL,
    delegation_approver         VARCHAR NOT NULL,
    staker_optout_window_blocks UINT256 NOT NULL,
    total_manta_stake           UINT256 NOT NULL,
    total_stake_reward          UINT256 NOT NULL,
    status                      SMALLINT NOT NULL DEFAULT 0,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_timestamp ON operator(timestamp);
CREATE INDEX IF NOT EXISTS operator_number ON operator(number);
CREATE INDEX IF NOT EXISTS operator_operator ON operator(operator);

CREATE TABLE IF NOT EXISTS operator_public_keys (
    guid                        VARCHAR PRIMARY KEY,
    operator                    VARCHAR NOT NULL,
    pubkey_hash                 VARCHAR NOT NULL,
    pubkey_g1                   UINT256 NOT NULL,
    pubkey_g2                   UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_public_keys_operator ON operator_public_keys(operator);


CREATE TABLE IF NOT EXISTS total_operator (
    guid                        VARCHAR PRIMARY KEY,
    to_block_number             UINT256 NOT NULL,
    op_count                    UINT256 NOT NULL,
    agg_pubkey_hash             VARCHAR NOT NULL,
    agg_pub_key                 UINT256 NOT NULL,
    op_index                    SMALLINT NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS total_operator_to_block_number ON total_operator(to_block_number);


CREATE TABLE IF NOT EXISTS operator_stake (
    guid                        VARCHAR PRIMARY KEY,
    staker                      VARCHAR NOT NULL,
    operator                    VARCHAR NOT NULL,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    manta_stake                 UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_stake_timestamp ON operator_stake(timestamp);
CREATE INDEX IF NOT EXISTS operator_stake_operator ON operator_stake(operator);
CREATE INDEX IF NOT EXISTS operator_stake_staker ON operator_stake(staker);


CREATE TABLE IF NOT EXISTS staker (
    guid                        VARCHAR PRIMARY KEY,
    staker                      VARCHAR NOT NULL,
    total_manta_stake           UINT256 NOT NULL,
    total_reward                UINT256 NOT NULL,
    total_claim_amount          UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_stake_staker ON operator_stake(staker);

CREATE TABLE IF NOT EXISTS staker_cliam (
      guid                        VARCHAR PRIMARY KEY,
      block_hash                  VARCHAR NOT NULL,
      number                      UINT256 NOT NULL,
      tx_hash                     VARCHAR NOT NULL UNIQUE,
      staker                      VARCHAR NOT NULL,
      claim_amount                UINT256 NOT NULL,
      timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS staker_cliam_staker ON staker_cliam(staker);
CREATE INDEX IF NOT EXISTS staker_cliam_timestamp ON staker_cliam(timestamp);

CREATE TABLE IF NOT EXISTS operator_cliam (
    guid                        VARCHAR PRIMARY KEY,
    block_hash                  VARCHAR NOT NULL,
    number                      UINT256 NOT NULL,
    tx_hash                     VARCHAR NOT NULL UNIQUE,
    operator                    VARCHAR NOT NULL,
    claim_amount                UINT256 NOT NULL,
    timestamp                   INTEGER NOT NULL UNIQUE CHECK (timestamp > 0)
);
CREATE INDEX IF NOT EXISTS operator_cliam_operator ON operator_cliam(operator);
CREATE INDEX IF NOT EXISTS operator_cliam_timestamp ON operator_cliam(timestamp);


