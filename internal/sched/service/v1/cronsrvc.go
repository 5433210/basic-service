package servicev1

import (
	"context"

	cron "github.com/robfig/cron/v3"

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
	}
}

type CronSrvc struct {
	ids   map[string]cron.EntryID
	cron  *cron.Cron
	store *store.Store
}

func NewCronSrvc(s *service) *CronSrvc {
	srvc := &CronSrvc{
		cron:  cron.New(),
		ids:   map[string]cron.EntryID{},
		store: s.store,
	}

	return srvc
}

func (c *CronSrvc) Load() error {
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

func (c *CronSrvc) Run() {
	c.cron.Run()
}

func (c *CronSrvc) Add(job apiv1.Job) error {
	schedule := job.Schedule
	// timezone := job.Timezone
	executor := job.Executor
	if e := executors.Get(executor.Name); e != nil {
		id, err := c.cron.AddJob(*schedule, NewCronJob(e, executor.Config, job.Data))
		if err != nil {
			return err
		}
		c.ids[job.Id] = id
	}

	return nil
}

func (c *CronSrvc) Remove(jobId string) error {
	id := c.ids[jobId]
	c.cron.Remove(id)

	return nil
}

func (c *CronSrvc) Stop() context.Context {
	return c.cron.Stop()
}
