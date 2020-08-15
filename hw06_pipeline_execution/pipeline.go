package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type (
	Stage func(in In) (out Out)
)

func makeStream(fn func(In) Out, in In, done In) Out {
	valueStream := make(Bi)
	go func() {
		defer close(valueStream)
		for item := range in {
			select {
			case <-done:
				return
			default:
			}
			select {
			case <-done:
				return
			case valueStream <- item:
			}
		}
	}()

	return fn(valueStream)
}

func worker(stage Stage, in In, done In) Out {
	return makeStream(stage, in, done)
}

func reduce(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	return reduce(worker(stages[0], in, done), done, stages[1:]...)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	fn := func(item In) Out {
		return item
	}

	return makeStream(fn, reduce(in, done, stages...), done)
}
