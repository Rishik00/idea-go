package llm

// Elaborate is introduced in phase 2. For phase 1 it is a no-op async stub.
func Elaborate(raw string) <-chan string {
	ch := make(chan string, 1)
	go func() {
		ch <- ""
		close(ch)
	}()
	return ch
}
