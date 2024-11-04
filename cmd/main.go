package main

import (
	"bufio"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	wpool "worker-pool/workers"
)

//go:embed help.txt
var help string

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	workerCount := flag.Int("count", 3, "Initial count of workers")
	separate := flag.String("sep", ",", "Separate of jobs")

	flag.Parse()

	wp := wpool.NewWorkerPool(*workerCount)

	workerControl()

	go wp.Run()

	for scanner.Scan() {
		var data = scanner.Text()

		str := strings.Split(data, " ")

		switch str[0] {
		case "/add":
			if len(str) <= 1 {
				fmt.Println("input zero count to add")
				break
			}
			count, err := strconv.Atoi(str[1])
			if err != nil {
				fmt.Println("invalid count of workers to add")
				break
			}
			err = wp.ChangeWorkersLimit(count)
			if errors.Is(err, wpool.ErrInvalidDelta) {
				fmt.Printf("invalid delta %d\n", count)
				break
			} else {
				fmt.Printf("Limit has set to %d\n", wp.GetWorkersLimit())
			}
		case "/remove":
			if len(str) <= 1 {
				fmt.Println("input zero count to remove")
				break
			}
			count, err := strconv.Atoi(str[1])
			if err != nil {
				fmt.Println("invalid count of workers to delete")
				break
			}
			err = wp.ChangeWorkersLimit(-count)
			if errors.Is(err, wpool.ErrInvalidDelta) {
				fmt.Printf("invalid delta %d\n", count)
				break
			} else {
				fmt.Printf("Limit has set to %d\n", wp.GetWorkersLimit())
			}
		case "/limit":
			fmt.Println("Count of workers", wp.GetWorkersLimit())
		case "/help":
			workerControl()
		case "/addjob":
			if len(str) <= 1 {
				fmt.Println("input zero jobs")
				break
			}

			for _, job := range str[1:] {
				for _, payload := range strings.Split(job, *separate) {
					wp.AddPayloadItem(payload)
				}
			}
		}
	}
}

func workerControl() {
	fmt.Println(help)
}
