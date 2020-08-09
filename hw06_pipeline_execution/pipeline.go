package hw06_pipeline_execution //nolint:golint,stylecheck

import (
	"sync"
)

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type Stage func(in In) (out Out)
type Data struct {
	mx sync.Mutex
	values map[int]I
}

func makeChannel(item I) <-chan I {
	ch := make(Bi)
	go func() {
		ch <- item
	}()
	return ch
}

func fillResult(d map[int]I, l int) Bi {
	result := make(Bi, l)
	for i := 0; i < l; i++ {
		// filter out empty values
		// they appear, when channel get closed
		if d[i] != nil {
			result <- d[i]
		}
	}
	close(result)
	return result
}

// Using closed channel to broadcast signal everywhere.
func checkDone(done In) Bi {
	chDone := make(Bi)
	go func(){
		for {
			select {
			case <-done:
				close(chDone)
				return
			default:
			}
		}
	}()
	return chDone
}

func execStages(in In, done In, stages ...Stage) Out {
	valueStream := make(Bi)
	go func() {
		defer close(valueStream)
		switch len(stages) {
		case 0:
			return
		case 1:
			valueStream <- <- stages[0](in)
			return
		default:
			for {
				select {
				case <-done:
					return
				case valueStream <- <- execStages(stages[0](in), done, stages[1:]...):
				}
			}
		}
	}()
	return valueStream
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var wg sync.WaitGroup
	var d Data
	d.values = make(map[int]I)
	var i int
	chDone := checkDone(done)
	for item := range in {
		wg.Add(1)
		go func(item I, i int){
			defer wg.Done()
			for {
				select {
				case <-chDone:
					return
				case value := <- execStages(makeChannel(item), chDone, stages...):
					d.mx.Lock()
					d.values[i] = value
					d.mx.Unlock()
					return
				}
			}
		}(item, i)
		i++
	}
	wg.Wait()
	return fillResult(d.values, i)
}