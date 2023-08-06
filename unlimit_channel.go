package UnlimitChannel

type UnlimitChan struct {
	In    chan interface{}
	Out   chan interface{}
	Name  string
	queue *Queue
}

func NewUnlimitChan(initSize int, name string) *UnlimitChan {
	ch := &UnlimitChan{
		In:    make(chan interface{}, initSize),
		Out:   make(chan interface{}, initSize),
		Name:  name,
		queue: NewQueue(),
	}

	go ch.process()

	return ch
}

func (ch *UnlimitChan) process() {
	defer close(ch.Out)

	drain := func() {
		for ch.queue.Size() != 0 {
			ch.Out <- ch.queue.Top()
			ch.queue.Pop()
		}
	}

	for {
		//in is closed
		value, ok := <-ch.In
		if !ok {
			drain()
			return
		}
		//out is not full
		select {
		case ch.Out <- value:
			continue
		default:
		}

		//out is full
		ch.queue.Push(value)

		//process cache
		for ch.queue.Size() != 0 {
			select {
			case val, ok := <-ch.In:
				if !ok {
					drain()
					return
				}
				ch.queue.Push(val)
			case ch.Out <- ch.queue.Top():
				ch.queue.Pop()
			}
		}
	}
}

func (ch *UnlimitChan) Size() int {
	return len(ch.In) + ch.queue.Size() + len(ch.Out)
}
