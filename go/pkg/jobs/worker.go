package jobs

import (
	"context"
	"sync"
)

func Worker(ctx context.Context, wg *sync.WaitGroup, queue <-chan *Pipeline) {

	defer wg.Done()
	for {
		select {
		case pipeline, ok := <-queue:
			if !ok {
				return
			}
			err := pipeline.Run(ctx)
			if err != nil {
				return
			}
			return
		case <-ctx.Done():
			return
		}
	}
}
