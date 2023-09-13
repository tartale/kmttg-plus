package jobs

import (
	"context"
	"sync"
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
