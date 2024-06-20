package utils

import (
	"log/slog"
	"sync"
)

// ChannelMerger merge results from multiple channels
func ChannelMerger[T any](inputs ...<-chan T) <-chan T {
	outCh := make(chan T)

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, c := range inputs {
		go func(c <-chan T) {
			defer wg.Done()
			for n := range c {
				slog.Info("read ch data", slog.Any("value", n))
				outCh <- n
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}
