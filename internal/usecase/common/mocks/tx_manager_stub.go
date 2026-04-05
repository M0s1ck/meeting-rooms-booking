package mocks

import "context"

type TxManagerStub struct {
	DoFunc func(ctx context.Context, fn func(context.Context) error) error
}

func (s *TxManagerStub) Do(ctx context.Context, fn func(context.Context) error) error {
	if s.DoFunc != nil {
		return s.DoFunc(ctx, fn)
	}
	return fn(ctx)
}

func NewTxMngrStubJustCall() *TxManagerStub {
	return &TxManagerStub{
		DoFunc: func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		},
	}
}
