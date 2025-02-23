package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	api "github.com/regen-network/regen-ledger/api/regen/ecocredit/basket/v1"
	baseapi "github.com/regen-network/regen-ledger/api/regen/ecocredit/v1"
	types "github.com/regen-network/regen-ledger/x/ecocredit/basket/types/v1"
)

func TestKeeper_BasketBalance(t *testing.T) {
	t.Parallel()
	s := setupBase(t)

	// add a basket
	basketDenom := testBasketDenom
	batchDenom := "bar"
	balance := "5.3"
	id, err := s.stateStore.BasketTable().InsertReturningID(s.ctx, &api.Basket{
		BasketDenom: basketDenom,
	})
	require.NoError(t, err)

	err = s.baseStore.BatchTable().Insert(s.ctx, &baseapi.Batch{
		Denom: batchDenom,
	})
	require.NoError(t, err)

	// add a balance
	require.NoError(t, s.stateStore.BasketBalanceTable().Insert(s.ctx, &api.BasketBalance{
		BasketId:   id,
		BatchDenom: batchDenom,
		Balance:    balance,
	}))

	// query
	res, err := s.k.BasketBalance(s.ctx, &types.QueryBasketBalanceRequest{
		BasketDenom: basketDenom,
		BatchDenom:  batchDenom,
	})
	require.NoError(t, err)
	require.Equal(t, balance, res.Balance)

	// bad query
	_, err = s.k.BasketBalance(s.ctx, &types.QueryBasketBalanceRequest{
		BasketDenom: batchDenom,
		BatchDenom:  basketDenom,
	})
	require.Error(t, err)

	// add another basket
	basketDenom = "foo1"
	basketName := "foo1.bar"
	err = s.stateStore.BasketTable().Insert(s.ctx, &api.Basket{
		BasketDenom: basketDenom,
		Name:        basketName,
	})
	require.NoError(t, err)

	// expect empty basket balance
	res, err = s.k.BasketBalance(s.ctx, &types.QueryBasketBalanceRequest{
		BasketDenom: basketDenom,
		BatchDenom:  batchDenom,
	})
	require.NoError(t, err)
	require.Equal(t, res.Balance, "0")
}
