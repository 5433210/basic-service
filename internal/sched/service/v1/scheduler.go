package servicev1

import (
	"context"
	"time"

	"github.com/google/uuid"
	cron "github.com/robfig/cron/v3"

	"wailik.com/internal/pkg/log"
	apiv1 "wailik.com/internal/sched/api/v1"
	"wailik.com/internal/sched/executors"
	"wailik.com/internal/sched/store"
)

type CronJob struct {
	scheduler *Scheduler
	jobId     string
	config    map[string]interface{}
	data      interface{}
	executor  executors.Executor
	onetime   bool
}

func NewCronJob(scheduler *Scheduler, jobId string, onetime bool, executor executors.Executor, cofnig map[string]interface{}, data interface{}) *CronJob {
	return &CronJob{
		scheduler: scheduler,
		jobId:     jobId,
		executor:  executor,
		config:    cofnig,
		data:      data,
		onetime:   onetime,
	}
}

func (c *CronJob) Run() {
	log.Debugf("job %+v running...", c)
	startAt := new(string)
	finishAt := new(string)
	id := uuid.NewString()
	err := c.executor.Execute(c.config, c.data, func() {
		*startAt = time.Now().String()
		err := c.scheduler.createExecutionInStore(apiv1.Execution{
			JobId:      c.jobId,
			Id:         id,
			Success:    nil,
			Output:     nil,
			StartedAt:  startAt,
			FinishedAt: nil,
		})
		if err != nil {
			log.Fatalf("createExecution err:%+v", err)
		}
	}, func(success bool, output string, ts time.Time) {
		*finishAt = ts.String()
		err := c.scheduler.updateExecutionInStore(apiv1.Execution{
			JobId:      c.jobId,
			Id:         id,
			Success:    &success,
			Output:     &output,
			StartedAt:  startAt,
			FinishedAt: finishAt,
		})
		if err != nil {
			log.Fatalf("updateExecution err:%+v", err)
		}
	})
	if err != nil {
		log.Fatalf("job execute error:%+v", err)
	}

	if c.onetime {
		err = c.scheduler.Remove(c.jobId)
		if err != nil {
			log.Fatalf("job remove error:%+v", err)
		}

		err = c.scheduler.removeJobInStore(c.jobId)
		if err != nil {
			log.Fatalf("job remove error:%+v", err)
		}
	}
}

type Scheduler struct {
	ids   map[string]cron.EntryID
	cron  *cron.Cron
	store *store.Store
}

func NewScheduler(store *store.Store) *Scheduler {
	srvc := &Scheduler{
		cron:  cron.New(cron.WithSeconds()),
		ids:   map[string]cron.EntryID{},
		store: store,
	}

	return srvc
}

func (c *Scheduler) Has(jobId string) bool {
	return c.ids[jobId] > 0
}

func (c *Scheduler) Load() error {
	log.Debug("scheduler loading...")
	client := c.store.Obtain()
	txn, err := client.Txn(true)
	if err != nil {
		return err
	}

	jobStore := &store.JobStore{}
	jobs, err := jobStore.RetrieveAll(txn)
	if err != nil {
		_ = txn.Rollback()

		return err
	}
	_ = txn.Commit(context.Background())
	for _, job := range *jobs {
		log.Debugf("load job(%+v)...", job.Id)
		if err = c.Add(job); err != nil {
			return err
		}
	}
	log.Debug("scheduler loaded")

	return nil
}

func (c *Scheduler) Start() {
	log.Debug("scheduler starting...")
	c.cron.Start()
}

func (c *Scheduler) Add(job apiv1.Job) error {
	var id cron.EntryID
	var err error

	if job.Disabled {
		log.Debugf("ignored disabled job:%+v", job.Id)

		return nil
	}

	if c.Has(job.Id) {
		err = c.Remove(job.Id) // remove first
		if err != nil {
			return err
		}
	}

	schedule := job.Schedule
	// timezone := job.Timezone
	executor := job.Executor
	executorInstance := executors.Get(executor.Name)
	if executorInstance == nil {
		return nil
	}
	cronjob := NewCronJob(c, job.Id, job.Onetime, executorInstance, executor.Config, job.Data)
	id, err = c.cron.AddJob(schedule, cronjob)
	if err != nil {
		return err
	}
	log.Debugf("job %+v added...", cronjob)
	c.ids[job.Id] = id

	return nil
}

func (c *Scheduler) Remove(jobId string) error {
	log.Debug("remove job")
	id := c.ids[jobId]
	c.cron.Remove(id)

	return nil
}

func (c *Scheduler) Stop() context.Context {
	return c.cron.Stop()
}

func (c *Scheduler) createExecutionInStore(execution apiv1.Execution) error {
	client := c.store.Obtain()
	txn, err := client.Txn(true)
	if err != nil {
		return err
	}

	executionStore := &store.ExecutionStore{}
	err = executionStore.Create(txn, execution)
	if err != nil {
		_ = txn.Rollback()

		return err
	}
	_ = txn.Commit(context.Background())

	return nil
}

func (c *Scheduler) updateExecutionInStore(execution apiv1.Execution) error {
	client := c.store.Obtain()
	txn, err := client.Txn(true)
	if err != nil {
		return err
	}

	executionStore := &store.ExecutionStore{}
	err = executionStore.Update(txn, execution)
	if err != nil {
		_ = txn.Rollback()

		return err
	}
	_ = txn.Commit(context.Background())

	return nil
}

func (c *Scheduler) removeJobInStore(jobId string) error {
	client := c.store.Obtain()
	txn, err := client.Txn(true)
	if err != nil {
		return err
	}

	jobStore := &store.JobStore{}
	err = jobStore.Delete(txn, jobId)
	if err != nil {
		_ = txn.Rollback()

		return err
	}
	_ = txn.Commit(context.Background())

	return nil
}
