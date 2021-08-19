package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func stageWorker(out Bi, done In, stage Stage, v interface{}) {
	in := make(Bi)
	r := stage(in)
	select {
	case in <- v:
		close(in)
	case <-done:
		break
	}
	select {
	case sr := <-r:
		out <- sr
	case <-done:
		break
	}
}

func stagePipe(in In, out Bi, done In, stage Stage) {
L:
	for {
		select {
		case v, ok := <-in:
			if !ok {
				break L
			}
			stageWorker(out, done, stage, v)
		case <-done:
			break L
		}
	}
	close(out)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return in
	}
	pineIn := make(Bi)
	sin := pineIn
	var out Bi
	for _, stage := range stages {
		out = make(Bi, 1000)
		go stagePipe(sin, out, done, stage)
		sin = out
	}
	go func() {
	L:
		for {
			select {
			case v, ok := <-in:
				if !ok {
					break L
				}
				select {
				case pineIn <- v:
				case <-done:
					break L
				}
			case <-done:
				break L
			}
		}
		close(pineIn)
	}()

	return out
}
