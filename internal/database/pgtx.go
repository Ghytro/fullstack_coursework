package database

import "context"

type TxContext struct {
	context.Context
	tx *TX
}

// ApplyTx прикрепляет транзакцию к контексту, только в том случае,
// если она к нему еще не прикреплена
func ApplyTx(ctx context.Context, tx *TX) context.Context {
	if _, ok := ctx.(*TxContext); ok {
		return ctx
	}
	return &TxContext{
		Context: ctx,
		tx:      tx,
	}
}

// GetTx получить транзакцию из конекста, если она там есть
func GetTx(ctx context.Context) *TX {
	if txCtx, ok := ctx.(*TxContext); ok {
		return txCtx.tx
	}
	return nil
}
