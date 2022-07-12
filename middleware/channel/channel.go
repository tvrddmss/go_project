package channel

func Send(a int, ch chan<- int) {
	ch <- a
	return
}

func Resvice(ch <-chan int) (a int) {
	a = <-ch
	return a
}
