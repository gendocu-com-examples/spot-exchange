package internal

import (
	"context"
	"errors"
	spot_exchange "git.gendocu.com/gendocu/SpotExchange.git/sdk/go"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"math/rand"
)

func NewService() *Service {
	return &Service{accounts: map[string][]*spot_exchange.Asset{}}
}

type Service struct {
	spot_exchange.UnimplementedSpotExchangeServer
	buy             []*spot_exchange.OrderRequest
	sell            []*spot_exchange.OrderRequest
	ordersToProcess chan *spot_exchange.OrderRequest
	accounts        map[string][]*spot_exchange.Asset
}

func (s Service) ListOrders(empty *emptypb.Empty, server spot_exchange.SpotExchange_ListOrdersServer) error {
	for _, order := range s.sell {
		err := server.Send(order.Order)
		if err != nil {
			return err
		}
	}
	for _, order := range s.buy {
		err := server.Send(order.Order)
		if err != nil {
			return err
		}
	}
	// TODO subscribe new orders
	return nil
}

func (s Service) PlaceOrder(server spot_exchange.SpotExchange_PlaceOrderServer) error {
	for {
		order, err := server.Recv()
		if errors.Is(err, io.EOF) {

		}
		if err != nil {
			return err
		}
		switch order.Order.OrderType {
		case spot_exchange.OrderType_Buy:
			s.sell = append(s.sell, order)
		case spot_exchange.OrderType_Sell:
			s.buy = append(s.buy, order)
		}
	}
}

func (s Service) Balance(ctx context.Context, req *spot_exchange.BalanceRequest) (*spot_exchange.BalanceResponse, error) {
	if _, ok := s.accounts[req.AccountId]; !ok {
		s.createNewAccount(ctx, req.AccountId)
	}
	return &spot_exchange.BalanceResponse{
			AccountId: req.AccountId,
			Assets:    s.accounts[req.AccountId],
		}, nil
}

func (s Service) createNewAccount(ctx context.Context, accountId string) {
	asset := &spot_exchange.Asset{}
	switch rand.Int31n(3) {
	case 0: // USD
		asset.AssetId = "USD"
		asset.Balance = 3000 + rand.Int31n(10000) // 3,000 <= .. < 13,000
	case 1: //APPL - Apple Inc share
		asset.AssetId = "APPL"
		asset.Balance = 3 + rand.Int31n(3) // 3 <= .. < 6
	case 2: // BTC - Bitcoin in satoshi (1 satoshi = 10^-8 BTC)
		asset.AssetId = "BTC"
		asset.Balance = 1e8 + rand.Int31n(1e8) // 1 BTC <= .. < 2 BTC
	}
	s.accounts[accountId] = []*spot_exchange.Asset{asset}
}

func (s Service) Execute() {
	for order := range s.ordersToProcess {
		s.ExecuteOrder(order)
	}
}

func (s Service) ExecuteOrder(order *spot_exchange.OrderRequest) {
	sell := OrderHeap{}
	switch order.Order.OrderType {
	case spot_exchange.OrderType_Buy:
		for {
			o := sell.Top()
			if o.Volume == 0 {
				break
			}
			if o.Price > order.Order.PriceLimit {
				break
			}
			vol := order.Order.Volume
			sell.Pop()
			if vol > o.Volume {
				vol = o.Volume
			} else {
				sell.Push(OrderItem{
					Price:  o.Price,
					Volume: o.Volume - vol,
				})
			}
			s.accounts[order.AccountId]
		}
	case spot_exchange.OrderType_Sell:
	}
}

func (s Service) SetAssets(u string, t []*spot_exchange.Asset) {
	s.accounts[u] = t
}
