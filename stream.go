package progress

import (
	"context"
	"time"
)

func Stream(ctx context.Context, progressor Progressor, interval time.Duration) <-chan Progress {
	results := make(chan Progress)

	go func() {
		defer close(results)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(interval):
				progress := progressor.Progress()
				results <- progress
				if progress.Complete() {
					return
				}
			}
		}
	}()
	return results
}
