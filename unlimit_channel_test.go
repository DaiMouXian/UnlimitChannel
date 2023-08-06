package UnlimitChannel

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestMakeNewUnlimitChan(t *testing.T) {
	ch := NewUnlimitChan(100, "test")

	for i := 1; i < 200; i++ {
		ch.In <- int64(i)
	}

	var count int64
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch.Out {
			fmt.Println(v.(int64))
			count += v.(int64)
		}
		fmt.Println("read completed")
	}()

	for i := 200; i <= 1000; i++ {
		ch.In <- int64(i)
	}
	close(ch.In)

	wg.Wait()

	if count != 500500 {
		t.Fatalf("expected 500500 but got %d", count)
	}
}

func TestLen_DataRace(t *testing.T) {
	ch := NewUnlimitChan(1, "test_data_race")
	stop := make(chan bool)
	for i := 0; i < 100; i++ { // may tweak the number of iterations
		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					ch.In <- 42
					<-ch.Out
				}
			}
		}()
	}

	for i := 0; i < 10000; i++ { // may tweak the number of iterations
		ch.Size()
	}
	close(stop)
}

func TestLen(t *testing.T) {
	ch := NewUnlimitChan(100, "test_len")

	for i := 1; i < 200; i++ {
		ch.In <- int64(i)
	}

	// wait ch processing in normal case
	time.Sleep(time.Second)
	assert.Equal(t, 0, len(ch.In))
	assert.Equal(t, 100, len(ch.Out))
	assert.Equal(t, 199, ch.Size())
	assert.Equal(t, 99, ch.queue.Size())

	for i := 0; i < 50; i++ {
		<-ch.Out
	}

	time.Sleep(time.Second)
	assert.Equal(t, 0, len(ch.In))
	assert.Equal(t, 100, len(ch.Out))
	assert.Equal(t, 149, ch.Size())
	assert.Equal(t, 49, ch.queue.Size())

	for i := 0; i < 149; i++ {
		<-ch.Out
	}

	time.Sleep(time.Second)
	assert.Equal(t, 0, len(ch.In))
	assert.Equal(t, 0, len(ch.Out))
	assert.Equal(t, 0, ch.Size())
	assert.Equal(t, 0, ch.queue.Size())
}
