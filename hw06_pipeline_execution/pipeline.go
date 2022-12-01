package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := doneStage(in, done)
	for _, s := range stages {
		out = doneStage(s(out), done)
	}
	return out
}

func doneStage(in In, done In) Out {
	out := make(Bi)

	go func(in In, done In) {
		defer func(in In) {
			close(out)
			for range in {
			}
		}(in)

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
