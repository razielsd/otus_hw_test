package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func stageWorker(out Bi, stage Stage, v interface{}) {
	in := make(Bi)
	r := stage(in)
	in <- v
	sr := <-r
	close(in)
	out <- sr
}

func stagePipe(in In, done In, stage Stage) Bi {
	out := make(Bi)
	go func() {
	L:
		for {
			select {
			case v, ok := <-in:
				if !ok {
					break L
				}
				stageWorker(out, stage, v)
			case <-done:
				break L
			}
		}
		close(out)
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}
	result := make(Bi)
	go func() {
		sin := make(Bi)
		out := sin
		for _, stage := range stages {
			out = stagePipe(out, done, stage)
		}
		go func() {
			// в принципе можно без селекта, но если не закроют канал, тогда горутина так и останется
			// for v := range in {
			// 	sin <- v
			// }
		L:
			for {
				select {
				case v, ok := <-in:
					if !ok {
						break L
					}
					sin <- v
				case <-done:
					break L
				}
			}
			close(sin)
		}()
		for v := range out {
			result <- v
		}
		close(result)
	}()

	return result
}
