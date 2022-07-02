package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	input, output, mutex := make(chan int64), make(chan user), new(sync.Mutex)

	for i := int64(0); i < pool; i++ {
		go worker(1, input, output, mutex)
	}

	go func(n int64, input chan<- int64) {
		for i := int64(0); i < n; i++ {
			input <- i
		}
		close(input)
	}(n, input)

	for i := int64(0); i < n; i++ {
		res = append(res, <-output)
	}
	return
}

func worker(id int, input <-chan int64, output chan<- user, mutex *sync.Mutex) {
	for id := range input {
		user := getOne(id)
		mutex.Lock()
		output <- user
		mutex.Unlock()
	}
}
