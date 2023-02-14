package goja

import (
	"context"
	"strings"

	"golang.org/x/time/rate"
)

// SetRateLimiter sets the rate limiter
func (r *Runtime) SetRateLimiter(limiter *rate.Limiter) {
	r.limiter = limiter
	if limiter == nil {
		return
	}

	r.fillBucket()
}

// NOTE: we should try to avoid making expensive operations within this
// function since it gets called millions of times per second.
func (r *Runtime) waitOneTick() {
	r.ticks++
	if r.limiter == nil {
		return
	}

	if r.limiterTicksLeft > 0 {
		r.limiterTicksLeft--
		return
	}
	r.fillBucket()

	ctx := r.vm.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	if waitErr := r.limiter.WaitN(ctx, r.limiterTicksLeft); waitErr != nil {
		if r.vm.ctx == nil {
			panic(waitErr)
		}
		if ctxErr := r.vm.ctx.Err(); ctxErr != nil {
			panic(ctxErr)
		}
		if strings.Contains(waitErr.Error(), "would exceed") {
			panic(context.DeadlineExceeded)
		}
		panic(waitErr)
	}
}

const burstDivisor = 5

func (r *Runtime) fillBucket() {
	r.limiterTicksLeft = r.limiter.Burst() / burstDivisor
}
