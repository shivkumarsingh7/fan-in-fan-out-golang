package main

import (
	"fmt"
	"sync"
)
func main() {
	randomNumbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	// generate the common channel with inputs
	inputChan := generatePipeline(randomNumbers)

	// Fan-out to 3 Go-routine
	c1 :=processData(inputChan)
	c2 :=processData(inputChan)
	c3 :=processData(inputChan)

	// Fan-in the resulting squared numbers
	c := fanIn(c1, c2, c3)


	// Do the summation
	for i := 0; i < len(randomNumbers); i++ {
		fmt.Println("Processed <- ", <-c, i)
	}
}

func generatePipeline(numbers []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range numbers {
			fmt.Println("Started <- ",n)
			out <- n
		}
		close(out)
	}()
	return out
}

func processData(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			//fmt.Println("processData ",n)
			out <- n
		}
		close(out)
	}()
	return out
}


func fanIn(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	c := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(cc <-chan int) {
		for n := range cc {
			c <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, cc := range cs {
		go output(cc)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(c)
	}()
	return c
}
