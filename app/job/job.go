package job

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

type Job struct {
	tick   time.Duration
	w      Worker
	cancel context.CancelFunc
}

type Worker interface {
	Init() error
	Work(ctx context.Context) error
	Stop()
	Name() string
}

func NewJob(tick time.Duration, w Worker) *Job {
	return &Job{
		tick: tick,
		w:    w,
	}
}

func (j *Job) Run(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("Job complite with error: %v", err)
		}
	}()
	ctx, j.cancel = context.WithCancel(ctx)
	log.WithField("job", j.w.Name()).Debug("Job start")

	err = j.w.Init()
	if err != nil {
		return err
	}

	tick := time.Duration(0)
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-time.After(tick):
			log.WithField("job", j.w.Name()).Debug("Job tick")
			tick = j.tick

			span, spanCtx := opentracing.StartSpanFromContext(ctx, j.w.Name())
			defer span.Finish()

			err := j.w.Work(spanCtx)
			if err != nil {
				log.WithField("job", j.w.Name()).Error("Job error: %v", err)
			}
		}
	}
}

func (j *Job) Stop() {
	log.WithField("job", j.w.Name()).Debug("Job stop")
	j.cancel()
	j.w.Stop()
}

// RunJobs ...
func RunJobs(ctx context.Context, jobs ...*Job) {
	for _, j := range jobs {
		go j.Run(ctx)
	}
}

// RunJobs ...
func StopJobs(jobs ...*Job) {
	for _, j := range jobs {
		j.Stop()
	}
}