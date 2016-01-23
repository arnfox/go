package main

type listener struct {
	subscribers map[string][]chan interface{}
}

// adds new listener to certain event (by name).
// E.q. add("trades", func(){return make(chan interface{})})
func (l *listener) add(lname string, handler func() chan interface{}) {
	go func() {
		for evt := range handler() {
			for _, s := range l.subscribers[lname] {
				s <- evt
			}
		}
		for _, s := range l.subscribers[lname] {
			close(s)
		}
	}()

}

// event subscriber. It will return channel to check upon new messages emitted
func (l *listener) subscribe(lname string) chan interface{} {
	ch := make(chan interface{})
	l.subscribers[lname] = append(l.subscribers[lname], ch)
	return ch
}

func main() {
	l := listener{map[string][]chan interface{}{}}

	evt := l.subscribe("taah")

	l.add("taah", func() chan interface{} {
		ch := make(chan interface{})
		go func() {
			for i := 0; i < 5; i++ {
				ch <- ("showing " + strconv.Itoa(i))
				<-time.After(500 * time.Millisecond)
			}
			close(ch)
		}()
		return ch
	})

	go func() {
		for {
			if m, ok := <-evt; ok {
				fmt.Println(m)
				<-time.After(700 * time.Millisecond)
			}
		}
	}()

	<-time.After(2 * time.Second)
}
