package workers

import "go.uber.org/zap"

type Job func() error

func New(logger *zap.SugaredLogger) *Client {
	return &Client{
		logger: logger,
		queue:  make(chan Job),
	}
}

type Client struct {
	logger *zap.SugaredLogger
	queue  chan Job
}

func (c *Client) Enqueue(j Job) {
	c.logger.Debug("enqueuing job")
	c.queue <- j
}

func (c *Client) Start() {
	for {
		select {
		case j := <-c.queue:
			c.logger.Debug("running job")
			go func() {
				if err := j(); err != nil {
					c.logger.Errorf("job failed: %s", err)
				}
			}()
		}
	}
}
