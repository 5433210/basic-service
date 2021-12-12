package servicev1

import (
	"context"

	cron "github.com/robfig/cron/v3"

	"wailik.com/internal/pkg/log"
	apiv1 "wailik.com/internal/sched/api/v1"
	"wailik.com/internal/sched/executors"
	"wailik.com/internal/sched/store"
)

type CronJob struct {
	config   map[string]interface{}
	data     interface{}
	executor executors.Executor
}

func NewCronJob(executor executors.Executor, cofnig map[string]interface{}, data interface{}) *CronJob {
	return &CronJob{
		executor: executor,
		config:   cofnig,
		data:     data,
	}
}

func (c *CronJob) Run() {
	err := c.executor.Execute(c.config, c.data)
	if err != nil {
		log.Fatalf("job execute error:%+v", err)
	}
}

type OnetimeJob struct {
	scheduler *Scheduler
	jobId     string
	config    map[string]interface{}
	data      interface{}
	executor  executors.Executor
}

func NewOnetimeJob(scheduler *Scheduler, jobId string, executor executors.Executor, cofnig map[string]interface{}, data interface{}) *OnetimeJob {
	return &OnetimeJob{
		scheduler: scheduler,
		jobId:     jobId,
		executor:  executor,
		config:    cofnig,
		data:      data,
	}
}

func (c *OnetimeJob) Run() {
	err := c.executor.Execute(c.config, c.data)
	if err != nil {
		log.Fatalf("job execute error:%+v", err)
	}
	err = c.scheduler.Remove(c.jobId)
	if err != nil {
		log.Fatalf("job remove error:%+v", err)
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

func (c *Scheduler) Load() error {
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
		if err = c.Add(job); err != nil {
			return err
		}
	}

	return nil
}

func (c *Scheduler) Run() {
	c.cron.Run()
}

func (c *Scheduler) Add(job apiv1.Job) error {
	var id cron.EntryID
	var err error

	if job.Disabled {
		log.Debugf("ignored disabled job:%+v", job.Id)

		return nil
	}

	schedule := job.Schedule
	// timezone := job.Timezone
	executor := job.Executor
	e := executors.Get(executor.Name)
	if e == nil {
		return nil
	}
	if job.Onetime {
		id, err = c.cron.AddJob(schedule, NewOnetimeJob(c, job.Id, e, executor.Config, job.Data))
		if err != nil {
			return err
		}
	} else {
		id, err = c.cron.AddJob(schedule, NewCronJob(e, executor.Config, job.Data))
		if err != nil {
			return err
		}
	}
	c.ids[job.Id] = id

	return nil
}

func (c *Scheduler) Remove(jobId string) error {
	id := c.ids[jobId]
	c.cron.Remove(id)

	return nil
}

func (c *Scheduler) Stop() context.Context {
	return c.cron.Stop()
}
