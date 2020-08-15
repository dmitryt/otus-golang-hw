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

func worker(stage Stage, in In, done In) Out {
	valueStream := make(Bi)
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case value := <-in:
				if value == nil {
					return
				}
				valueStream <- value
			}
		}
	}()

	return stage(valueStream)
}

func reduce(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}

	return reduce(worker(stages[0], in, done), done, stages[1:]...)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	return reduce(in, done, stages...)
}
