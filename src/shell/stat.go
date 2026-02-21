package shell

import (
	stdcontext "context"
	"os"
	"time"
)

const defaultStatTimeout = time.Second

var statTimeout = defaultStatTimeout

func statWithTimeout(path string, timeout time.Duration) (os.FileInfo, error) {
	ctx, cancel := stdcontext.WithTimeout(stdcontext.Background(), timeout)
	defer cancel()

	type result struct {
		info os.FileInfo
		err  error
	}

	ch := make(chan result, 1)

	go func() {
		info, err := os.Stat(path)
		select {
		case ch <- result{
			info: info,
			err:  err,
		}:
		case <-ctx.Done():
		}
	}()

	select {
	case res := <-ch:
		return res.info, res.err
	case <-ctx.Done():
		return nil, stdcontext.DeadlineExceeded
	}
}

func SetStatTimeout(timeout time.Duration) {
	if timeout <= 0 {
		statTimeout = defaultStatTimeout
		return
	}

	statTimeout = timeout
}
