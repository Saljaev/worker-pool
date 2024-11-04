package workers

import (
	"errors"
	"github.com/google/uuid"
)

type WorkerPool struct {
	limiter chan struct{}
	payload chan interface{}
}

func NewWorkerPool(limit int) *WorkerPool {
	return &WorkerPool{
		limiter: make(chan struct{}, limit),
		payload: make(chan interface{}),
	}
}

// ChangeWorkersLimit change count of workers by delta
func (wp *WorkerPool) ChangeWorkersLimit(delta int) error {
	limit := cap(wp.limiter) + delta
	if limit < 0 {
		return ErrInvalidDelta
	}

	wp.limiter = make(chan struct{}, limit)
	return nil
}

func (wp *WorkerPool) AddPayloadItem(payloadItem interface{}) {
	wp.payload <- payloadItem
}

func (wp *WorkerPool) GetWorkersLimit() int {
	return cap(wp.limiter)
}

// Run start WorkerPool
func (wp *WorkerPool) Run() {
	for {
		payloadItem := <-wp.payload

		if cap(wp.limiter) > 0 {
			wp.limiter <- struct{}{}

			// Используем функцию для абстракции, чтобы сам работник
			// ничего не знал о канале
			go do(uuid.New(), payloadItem, func() {
				<-wp.limiter
			})
		}

	}
}

var ErrInvalidDelta = errors.New("invalid delta")
