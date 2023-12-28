package main

import (
	"log"
	"temporal-demo/submit_review"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, submit_review.TaskQueue, worker.Options{})
	a := &submit_review.Activities{}
	w.RegisterWorkflow(submit_review.SubmitAndReview)
	w.RegisterActivity(a)
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
