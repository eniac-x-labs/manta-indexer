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


