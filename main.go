package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

const TotalNumbers = 10000000
const NumWorkers = 8

// collatzSteps обчислює кількість кроків, необхідних для досягнення 1 відповідно до гіпотези Колаца.
func collatzSteps(n int64) int64 {
	if n <= 0 {
		return 0
	}

	steps := int64(0)
	current := n

	for current != 1 {
		if current%2 == 0 {
			current /= 2
		} else {
			current = 3*current + 1
		}
		steps++
	}
	return steps
}

// job представляє задачу для обробки (число, для якого треба порахувати кроки)
type job struct {
	number int64
}

type result struct {
	steps int64
}

func workerPool(jobs <-chan job, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := range jobs {
		steps := collatzSteps(j.number)
		results <- result{steps: steps}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Паралельне обчислення градини Колаца")
	fmt.Printf("Кількість чисел: %d\n", TotalNumbers)
	fmt.Printf("Кількість робітників (Go-рутин): %d\n\n", NumWorkers)

	jobs := make(chan job, TotalNumbers)
	results := make(chan result, TotalNumbers)
	var wg sync.WaitGroup

	startTime := time.Now()

	for w := 1; w <= NumWorkers; w++ {
		wg.Add(1)
		go workerPool(jobs, results, &wg)
	}
	fmt.Printf("Запущено %d робітників...\n", NumWorkers)

	go func() {
		for i := int64(1); i <= TotalNumbers; i++ {
			jobs <- job{number: i}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	totalSteps := int64(0)
	count := 0
	for r := range results {
		totalSteps += r.steps
		count++
	}

	elapsedTime := time.Since(startTime)

	if count != TotalNumbers {
		fmt.Printf("\nПомилка: Оброблені %d чисел замість очікуваних %d.\n", count, TotalNumbers)
		os.Exit(1)
	}

	averageSteps := float64(totalSteps) / float64(TotalNumbers)

	fmt.Println("\nОбчислення завершено!!")
	fmt.Printf("Загальний час виконання: %s\n", elapsedTime)
	fmt.Printf("Загальна кількість кроків: %d\n", totalSteps)
	fmt.Printf("Середня кількість кроків до 1: %.4f\n", averageSteps)
}
