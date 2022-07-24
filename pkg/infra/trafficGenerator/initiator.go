package trafficGenerator

import (
	"context"
	"github.com/wsw365904/wswlog/wlogging"

	"github.com/wsw365904/tape/pkg/infra"
	"github.com/wsw365904/tape/pkg/infra/basic"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

type Initiator struct {
	Num     int
	Burst   int
	R       float64
	Config  basic.Config
	Crypto  infra.Crypto
	Logger  *wlogging.WswLogger
	Raw     chan *basic.TracingProposal
	ErrorCh chan error
}

func (initiator *Initiator) Start() {
	limit := rate.Inf
	ctx := context.Background()
	if initiator.R > 0 {
		limit = rate.Limit(initiator.R)
	}
	limiter := rate.NewLimiter(limit, initiator.Burst)
	i := 0
	for {
		if initiator.Num > 0 {
			if i == initiator.Num {
				return
			}
			i++
		}
		prop, err := CreateProposal(
			initiator.Crypto,
			initiator.Logger,
			initiator.Config.Channel,
			initiator.Config.Chaincode,
			initiator.Config.Version,
			initiator.Config.Args...,
		)
		if err != nil {
			initiator.ErrorCh <- errors.Wrapf(err, "error creating proposal")
			return
		}

		if err = limiter.Wait(ctx); err != nil {
			initiator.ErrorCh <- errors.Wrapf(err, "error creating proposal")
			return
		}
		initiator.Raw <- prop
	}
}
