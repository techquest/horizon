package db

import (
	"database/sql"
	"github.com/stellar/horizon/txsub"
	"golang.org/x/net/context"
)

type ResultProvider struct {
	Ctx     context.Context
	Core    *sql.DB
	History *sql.DB
}

func (rp *ResultProvider) ResultByHash(hash string) txsub.Result {
	var (
		hr TransactionRecord
	)

	hq := TransactionByHashQuery{
		SqlQuery: SqlQuery{rp.History},
		Hash:     hash,
	}

	err := Get(rp.Ctx, hq, &hr)
	if err == nil {
		return txResultFromTransactionRecord(hr)
	}

	if err != ErrNoResults {
		return txsub.Result{Err: err}
	}

	//TODO: check stellar-core for results as well

	return txsub.Result{Err: txsub.ErrNoResults}
}

func (rp *ResultProvider) ResultByAddressAndSequence(addr string, seq uint64) txsub.Result {
	var hr TransactionRecord

	hq := TransactionByAddressAndSequence{
		SqlQuery: SqlQuery{rp.History},
		Address:  addr,
		Sequence: seq,
	}

	err := Get(rp.Ctx, hq, &hr)
	if err == nil {
		return txResultFromTransactionRecord(hr)
	}

	if err != ErrNoResults {
		return txsub.Result{Err: err}
	}

	//TODO: check stellar-core for results as well

	return txsub.Result{Err: txsub.ErrNoResults}
}

func txResultFromTransactionRecord(hr TransactionRecord) txsub.Result {
	return txsub.Result{
		Hash:           hr.TransactionHash,
		LedgerSequence: hr.LedgerSequence,
		EnvelopeXDR:    hr.TxEnvelope.String,
		ResultXDR:      hr.TxResult.String,
		ResultMetaXDR:  hr.TxMeta.String,
	}
}
