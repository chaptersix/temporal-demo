package main

import (
	"context"
	"log"
	"net/http"
	"temporal-demo/submit_review"
	"temporal-demo/webui/view"

	"github.com/labstack/echo/v4"
	"go.temporal.io/sdk/client"
)

func main() {
	// Create a new instance of Echo
	e := echo.New()

	// Route to handle the GET request on the home page

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/signup")
	})
	e.GET("/signup", func(c echo.Context) error {
		return view.Submit().Render(c.Request().Context(), c.Response())
	})
	e.POST("/signup/submit-form", func(ec echo.Context) error {
		name := ec.FormValue("name")
		email := ec.FormValue("email")
		history := ec.FormValue("criminalHistory")
		c, err := client.Dial(client.Options{})
		if err != nil {
			log.Fatalln("Failed to create Temporal client", err)
		}
		defer c.Close()
		// Start a workflow execution
		_, err = c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
			ID:        "emailWorkflow_" + email,
			TaskQueue: submit_review.TaskQueue,
		}, submit_review.SubmitAndReview, submit_review.Applciation{
			Email:                      email,
			Name:                       name,
			HasCriminalHistory:         len(history) > 0,
			CriminalHistoryDescription: history,
		})
		if err != nil {
			log.Fatalln("Failed to start workflow", err)
		}
		return ec.String(http.StatusOK, "Form submitted successfully")
	})
	e.GET("/admin", func(c echo.Context) error {
		return view.Admin([]submit_review.Applciation{{Name: "TestName", Email: "test@mail.com"}}).Render(c.Request().Context(), c.Response())
	})

	// Start the server on port 8080
	e.Logger.Fatal(e.Start(":8080"))
}

// func eample() {
// 	c, err := client.Dial(client.Options{})
// 	if err != nil {
// 		log.Fatalln("Failed to create Temporal client", err)
// 	}
// 	defer c.Close()

// 	email := "example@email.com" // Replace with the actual email address.

// 	// Start a workflow execution
// 	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
// 		ID:        "emailWorkflow_" + email,
// 		TaskQueue: submit_review.TaskQueue,
// 	}, submit_review.SubmitAndReview, submit_review.Applciation{
// 		Email: email,
// 		Name:  "Test",
// 	})
// 	if err != nil {
// 		log.Fatalln("Failed to start workflow", err)
// 	}

// 	err = c.SignalWorkflow(context.Background(), "emailWorkflow_"+email, we.GetRunID(), submit_review.ReviewChanName, submit_review.ReviewSignal{Outcome: submit_review.OutcomeApproved, Email: email})
// 	if err != nil {
// 		log.Fatalln("Failed to start signal", err)
// 	}
// 	log.Println("Workflow started", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
// }
