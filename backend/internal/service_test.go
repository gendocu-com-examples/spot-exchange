package internal_test

import (
	"context"
	spot_exchange "git.gendocu.com/gendocu/SpotExchange.git/sdk/go"
	"github.com/gendocu-com-examples/spot-exchange/backend/internal"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBalance_NewUser(t *testing.T) {
	ctx := context.Background()
	// Given
	s := internal.NewService()
	accountId := "2af6cb3f-0310-43a2-8e4c-5a400ef648a4"

	// When
	b, err := s.Balance(ctx, &spot_exchange.BalanceRequest{
		AccountId: accountId,
	})
	b2, err2 := s.Balance(ctx, &spot_exchange.BalanceRequest{
		AccountId: accountId,
	})

	// Then
	assert.NoError(t, err)
	assert.NoError(t, err2)
	if assert.Len(t, b.Assets, 1) &&
		assert.Len(t, b2.Assets, 1) {
		assert.Equal(t, b.Assets[0], b2.Assets[0])
	}
}

func TestExecuteOrders(t *testing.T) {
	t.Parallel()
	asset, asset2 := "A1", "A2"
	tst := []struct {
		name string
		orders []*spot_exchange.OrderRequest
		account *spot_exchange.BalanceResponse
		account2 *spot_exchange.BalanceResponse
	}{
		{name: "buy_sell",
			orders: []*spot_exchange.OrderRequest{
				{
					AccountId: "1",
					Order: &spot_exchange.Order{
						OrderType: spot_exchange.OrderType_Buy, AssetId: asset, Volume: 1, OrderId: uuid.NewString(), PriceLimit: 100},
				},
				{
					AccountId: "2",
					Order: &spot_exchange.Order{
						OrderType: spot_exchange.OrderType_Sell, AssetId: asset, Volume: 1, OrderId: uuid.NewString(), PriceLimit: 100},
				},
			}},
		{name: "sell_buy",
			orders: []*spot_exchange.OrderRequest{
				{
					AccountId: "1",
					Order: &spot_exchange.Order{
						OrderType: spot_exchange.OrderType_Sell, AssetId: asset, Volume: 1, OrderId: uuid.NewString(), PriceLimit: 100},
				},
				{
					AccountId: "2",
					Order: &spot_exchange.Order{
						OrderType: spot_exchange.OrderType_Buy, AssetId: asset, Volume: 1, OrderId: uuid.NewString(), PriceLimit: 100},
				},
			}},
			{name: "sell_half_buy",
			orders: []*spot_exchange.OrderRequest{
				{
					AccountId: "1",
					Order: &spot_exchange.Order{
						OrderType: spot_exchange.OrderType_Sell, AssetId: asset, Volume: 1, OrderId: uuid.NewString(), PriceLimit: 100},
				},
				{
					AccountId: "2",
					Order: &spot_exchange.Order{
						OrderType: spot_exchange.OrderType_Buy, AssetId: asset, Volume: 1, OrderId: uuid.NewString(), PriceLimit: 100},
				},
			}},
		{name: "different assets",
			orders: []*spot_exchange.OrderRequest{
			{
				AccountId: "1",
				Order: &spot_exchange.Order{
					OrderType: spot_exchange.OrderType_Buy, AssetId: asset, Volume: 1, OrderId: uuid.NewString(), PriceLimit: 100},
			},
			{
				AccountId: "2",
				Order: &spot_exchange.Order{
					OrderType: spot_exchange.OrderType_Sell, AssetId: asset2, Volume: 1, OrderId: uuid.NewString(), PriceLimit: 100},
			},
		}},
	}
	for _, test := range tst {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// Given
			u1, u2 := uuid.NewString(), uuid.NewString()
			s := internal.NewService()
			s.SetAssets(u1, nil)
			s.SetAssets(u2, nil)

			// When
			for _, order := range test.orders {
				if order.AccountId == "1" {
					order.AccountId = u1
				} else {
					order.AccountId = u2
				}
				s.ExecuteOrder(order)
			}
			b1, err := s.Balance(context.Background(), &spot_exchange.BalanceRequest{
				AccountId: u1,
			})
			b2, err2 := s.Balance(context.Background(), &spot_exchange.BalanceRequest{
				AccountId: u1,
			})

			// Then
			assert.NoError(t, err)
			assert.NoError(t, err2)
			assert.EqualValues(t, b1.Assets, test.account.Assets)
			assert.EqualValues(t, b2.Assets, test.account2.Assets)
		})
	}
}
