package main

import (
	"fmt"
	"sync"
	"time"
)

type ChopS struct{ sync.Mutex }

type Philo struct {
	leftCS, rightCS *ChopS
	id              int
	eats            int
	reqToEat        chan int
	startEat        chan int
	finishEat       chan int
}

var wg sync.WaitGroup

func (p Philo) eat() {

	for {
		if p.eats > 2 {
			p.finishEat <- p.id
			wg.Done()
			return
		}
		//sen request to eat
		p.reqToEat <- p.id
		eat := <-p.startEat
		//permission from host
		if eat == 1 {
			p.leftCS.Lock()
			p.rightCS.Lock()

			fmt.Printf("starting to eat %d\n", p.id)
			time.Sleep(time.Millisecond * 500) //sleep for 500ms while eating
			fmt.Printf("finishing eating %d\n", p.id)
			p.eats++
			p.finishEat <- p.id

			p.rightCS.Unlock()
			p.leftCS.Unlock()
		}
	}
}

func host(reqToEat, startEat, finishEat, quit chan int) {
	var eatingPhilos int //how many philos eating now

	for {
		select {
		case <-reqToEat:
			if eatingPhilos < 2 {
				eatingPhilos++
				startEat <- 1
			} else {
				startEat <- 0
			}
		case <-finishEat:
			eatingPhilos--
		case <-quit:
			return
		}
	}

}

func main() {
	finishEat := make(chan int)
	startEat := make(chan int)
	requestToEat := make(chan int)
	quit := make(chan int)
	//create forks
	CSticks := make([]*ChopS, 5)
	for i := 0; i < 5; i++ {
		CSticks[i] = new(ChopS)
	}
	//create philos
	philos := make([]*Philo, 5)
	for i := 0; i < 5; i++ {
		philos[i] = &Philo{CSticks[i], CSticks[(i+1)%5], i + 1, 0, requestToEat, startEat, finishEat}
	}
	//start to eat
	wg.Add(5)
	go host(requestToEat, startEat, finishEat, quit) //host goroutine
	for i := 0; i < 5; i++ {
		go philos[i].eat()
	}
	wg.Wait()

}
