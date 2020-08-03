package client_bk

import "time"

type Breaker interface {
	Call(func() error, time.Duration) error
	Fail()
	Success()
	Ready()
}
