package jobs

import (
	"context"
	"sync"

	"github.com/tartale/kmttg-plus/go/pkg/config"
)

func Worker(ctx context.Context, wg *sync.WaitGroup, queue <-chan *Pipeline) {

	defer wg.Done()
	for {
		select {
		case pipeline := <-queue:
			pipeline.Run(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func RunWorkerPool(ctx context.Context) {

	var wg sync.WaitGroup

	for i := 0; i < config.Values.MaxBackgroundTasks; i++ {
		wg.Add(1)
		// fan out worker goroutines
		//reading from jobs channel and
		//pushing calcs into results channel
		go Worker(ctx, &wg, pipelineQueue)
	}

	wg.Wait()
}
