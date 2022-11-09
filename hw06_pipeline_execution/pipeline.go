package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, s := range stages {
		out = doneStage(done, s(out))
	}
	return out
}

func doneStage(done, in In) Out {
	out := make(Bi)

	go func(in, done In) {
		defer func(in, done In) {
			close(out)
			for range in {
			}
		}(in, done)

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}(in, done)
	return out
}
