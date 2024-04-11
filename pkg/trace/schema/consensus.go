package schema

import (
	"github.com/cometbft/cometbft/pkg/trace"
	"github.com/cometbft/cometbft/types"
)

// ConsensusTables returns the list of tables that are used for consensus
// tracing.
func ConsensusTables() []string {
	return []string{
		RoundStateTable,
		BlockPartsTable,
		BlockTable,
		VoteTable,
		ConsensusStateTable,
		ProposalTable,
	}
}

// Schema constants for the consensus round state tracing database.
const (
	// RoundStateTable is the name of the table that stores the consensus
	// state traces.
	RoundStateTable = "consensus_round_state"
)

// RoundState describes schema for the "consensus_round_state" table.
type RoundState struct {
	Height int64 `json:"height"`
	Round  int32 `json:"round"`
	Step   uint8 `json:"step"`
}

// Table returns the table name for the RoundState struct.
func (r RoundState) Table() string {
	return RoundStateTable
}

func (r RoundState) GetHeight() int64 {
	return r.Height
}

func (r RoundState) GetRound() int32 {
	return r.Round
}

func (r RoundState) GetStep() uint8 {
	return r.Step
}

// WriteRoundState writes a tracing point for a tx using the predetermined
// schema for consensus state tracing.
func WriteRoundState(client trace.Tracer, height int64, round int32, step uint8) {
	client.Write(RoundState{Height: height, Round: round, Step: step})
}

// Schema constants for the "consensus_block_parts" table.
const (
	// BlockPartsTable is the name of the table that stores the consensus block
	// parts.
	BlockPartsTable = "consensus_block_parts"
)

// BlockPart describes schema for the "consensus_block_parts" table.
type BlockPart struct {
	Height       int64        `json:"height"`
	Round        int32        `json:"round"`
	Index        int32        `json:"index"`
	Catchup      bool         `json:"catchup"`
	Peer         string       `json:"peer"`
	TransferType TransferType `json:"transfer_type"`
}

// Table returns the table name for the BlockPart struct.
func (b BlockPart) Table() string {
	return BlockPartsTable
}

func (b BlockPart) GetHeight() int64 {
	return b.Height
}

func (b BlockPart) GetRound() int32 {
	return b.Round
}

func (b BlockPart) GetIndex() int32 {
	return b.Index
}

func (b BlockPart) GetPeer() string {
	return b.Peer
}

func (b BlockPart) GetTransferType() TransferType {
	return b.TransferType
}

func (b BlockPart) GetCatchup() bool {
	return b.Catchup
}

// WriteBlockPart writes a tracing point for a BlockPart using the predetermined
// schema for consensus state tracing.
func WriteBlockPart(
	client trace.Tracer,
	height int64,
	round int32,
	index uint32,
	catchup bool,
	peer string,
	transferType TransferType,
) {
	// this check is redundant to what is checked during client.Write, although it
	// is an optimization to avoid allocations from the map of fields.
	if !client.IsCollecting(BlockPartsTable) {
		return
	}
	client.Write(BlockPart{
		Height:       height,
		Round:        round,
		Index:        int32(index),
		Catchup:      catchup,
		Peer:         peer,
		TransferType: transferType,
	})
}

// Schema constants for the consensus votes tracing database.
const (
	// VoteTable is the name of the table that stores the consensus
	// voting traces.
	VoteTable = "consensus_vote"
)

// Vote describes schema for the "consensus_vote" table.
type Vote struct {
	Height                   int64        `json:"height"`
	Round                    int32        `json:"round"`
	VoteType                 string       `json:"vote_type"`
	VoteHeight               int64        `json:"vote_height"`
	VoteRound                int32        `json:"vote_round"`
	VoteMillisecondTimestamp int64        `json:"vote_unix_millisecond_timestamp"`
	ValidatorAddress         string       `json:"vote_validator_address"`
	Peer                     string       `json:"peer"`
	TransferType             TransferType `json:"transfer_type"`
}

func (v Vote) Table() string {
	return VoteTable
}

func (v Vote) GetHeight() int64 {
	return v.Height
}

func (v Vote) GetRound() int32 {
	return v.Round
}

func (v Vote) GetVoteType() string {
	return v.VoteType
}

func (v Vote) GetPeer() string {
	return v.Peer
}

func (v Vote) GetTransferType() TransferType {
	return v.TransferType
}

// WriteVote writes a tracing point for a vote using the predetermined
// schema for consensus vote tracing.
func WriteVote(client trace.Tracer,
	height int64, // height of the current peer when it received/sent the vote
	round int32, // round of the current peer when it received/sent the vote
	vote *types.Vote, // vote received by the current peer
	peer string, // the peer from which it received the vote or the peer to which it sent the vote
	transferType TransferType, // download (received) or upload(sent)
) {
	client.Write(Vote{
		Height:                   height,
		Round:                    round,
		VoteType:                 vote.Type.String(),
		VoteHeight:               vote.Height,
		VoteRound:                vote.Round,
		VoteMillisecondTimestamp: vote.Timestamp.UnixMilli(),
		ValidatorAddress:         vote.ValidatorAddress.String(),
		Peer:                     peer,
		TransferType:             transferType,
	})
}

const (
	// BlockTable is the name of the table that stores metadata about consensus blocks.
	BlockTable = "consensus_block"
)

// BlockSummary describes schema for the "consensus_block" table.
type BlockSummary struct {
	Height                   int64  `json:"height"`
	UnixMillisecondTimestamp int64  `json:"unix_millisecond_timestamp"`
	TxCount                  int    `json:"tx_count"`
	SquareSize               uint64 `json:"square_size"`
	BlockSize                int    `json:"block_size"`
	Proposer                 string `json:"proposer"`
	LastCommitRound          int32  `json:"last_commit_round"`
}

func (b BlockSummary) Table() string {
	return BlockTable
}

func (b BlockSummary) GetHeight() int64 {
	return b.Height
}

func (b BlockSummary) GetPoposer() string {
	return b.Proposer
}

// WriteBlockSummary writes a tracing point for a block using the predetermined
func WriteBlockSummary(client trace.Tracer, block *types.Block, size int) {
	client.Write(BlockSummary{
		Height:                   block.Height,
		UnixMillisecondTimestamp: block.Time.UnixMilli(),
		TxCount:                  len(block.Data.Txs),
		SquareSize:               block.SquareSize,
		BlockSize:                size,
		Proposer:                 block.ProposerAddress.String(),
		LastCommitRound:          block.LastCommit.Round,
	})
}

const (
	ConsensusStateTable = "consensus_state"
)

type ConsensusStateUpdateType string

const (
	ConsensusNewValidBlock      ConsensusStateUpdateType = "new_valid_block"
	ConsensusNewRoundStep       ConsensusStateUpdateType = "new_round_step"
	ConsensusVoteSetBits        ConsensusStateUpdateType = "vote_set_bits"
	ConsensusVoteSet23Prevote   ConsensusStateUpdateType = "vote_set_23_prevote"
	ConsensusVoteSet23Precommit ConsensusStateUpdateType = "vote_set_23_precommit"
	ConsensusHasVote            ConsensusStateUpdateType = "has_vote"
	ConsensusPOL                ConsensusStateUpdateType = "pol"
)

type ConsensusState struct {
	Height       int64        `json:"height"`
	Round        int32        `json:"round"`
	UpdateType   string       `json:"update_type"`
	Peer         string       `json:"peer"`
	TransferType TransferType `json:"transfer_type"`
}

func (c ConsensusState) Table() string {
	return ConsensusStateTable
}

func (c ConsensusState) GetHeight() int64 {
	return c.Height
}

func (c ConsensusState) GetRound() int32 {
	return c.Round
}

func (c ConsensusState) GetPeer() string {
	return c.Peer
}

func (c ConsensusState) GetTransferType() TransferType {
	return c.TransferType
}

func (c ConsensusState) GetUpdateType() string {
	return c.UpdateType
}

func WriteConsensusState(
	client trace.Tracer,
	height int64,
	round int32,
	step uint8,
	peer string,
	updateType ConsensusStateUpdateType,
	transferType TransferType,
) {
	client.Write(ConsensusState{
		Height:     height,
		Round:      round,
		Peer:       peer,
		UpdateType: string(updateType),
	})
}

const (
	ProposalTable = "consensus_proposal"
)

type Proposal struct {
	Height       int64        `json:"height"`
	Round        int32        `json:"round"`
	Peer         string       `json:"peer"`
	TransferType TransferType `json:"transfer_type"`
}

func (p Proposal) Table() string {
	return ProposalTable
}

func (p Proposal) GetHeight() int64 {
	return p.Height
}

func (p Proposal) GetRound() int32 {
	return p.Round
}

func (p Proposal) GetPeer() string {
	return p.Peer
}

func (p Proposal) GetTransferType() TransferType {
	return p.TransferType
}

func WriteProposal(
	client trace.Tracer,
	height int64,
	round int32,
	peer string,
	transferType TransferType,
) {
	client.Write(Proposal{
		Height:       height,
		Round:        round,
		Peer:         peer,
		TransferType: transferType,
	})
}
