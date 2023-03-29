package executor

import "context"

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		tmp := make(chan any)
		go func(in In, out chan any) {
			defer close(out)
			for {
				select {
				case <-ctx.Done():
					return
				case v, open := <-in:
					if !open {
						return
					}
					tmp <- v
				}
			}
		}(out, tmp)

		out = stage(tmp)
	}
	return out
}
