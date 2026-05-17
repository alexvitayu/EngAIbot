package tg_bot

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/alexvitayu/EngAIbot/internal/scheduler"
	"github.com/alexvitayu/EngAIbot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WorkerPool struct {
	bot          *tgbotapi.BotAPI
	updates      tgbotapi.UpdatesChannel
	workersCount int
	wg           *sync.WaitGroup
	service      *service.PhraseService
	Sched        *scheduler.Scheduler
	stopChan     chan struct{}
	errChan      chan error
}

func NewWorkerPool(bot *tgbotapi.BotAPI, u tgbotapi.UpdatesChannel, count int,
	wg *sync.WaitGroup, s *service.PhraseService, sched *scheduler.Scheduler) *WorkerPool {
	return &WorkerPool{
		bot:          bot,
		updates:      u,
		workersCount: count,
		wg:           wg,
		service:      s,
		Sched:        sched,
		stopChan:     make(chan struct{}),
		errChan:      make(chan error),
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 1; i <= wp.workersCount; i++ {
		wp.wg.Add(1)
		go func(idx int) {
			wp.worker(ctx, idx)
			slog.Debug(fmt.Sprintf("worker %d started", idx))
		}(i)
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.stopChan)
}
