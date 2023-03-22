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
			for v := range in {
				select {
				case <-ctx.Done():
					return
				default:
					tmp <- v
				}
			}
		}(out, tmp)

		out = stage(tmp)
	}
	return out
}
