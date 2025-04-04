package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	windowPrevState []int
	windowCurrState []int
	numbers         []int
	avg             float64
}

var state []int

func main() {
	windowSize := 10
	e := echo.New()

	e.GET("/numbers/:numberId", avgCalc)

	e.Logger.Fatal(e.Start(":8080"))
}

func avgCalc(c echo.Context) error {
	var res Response
	res.windowPrevState = state

	id := c.Param("numberId")
	var url string

	switch id {
	case "p":
		url = "http://20.244.56.144/evaluation-service/primes"
	case "f":
		url = "http://20.244.56.144/evaluation-service/fibo"
	case "e":
		url = "http://20.244.56.144/evaluation-service/even"
	case "r":
		url = "http://20.244.56.144/evaluation-service/rand"

	}

	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzQzNzQyODM0LCJpYXQiOjE3NDM3NDI1MzQsImlzcyI6IkFmZm9yZG1lZCIsImp0aSI6IjYwZGNjNzZhLWQ4OWYtNDQyYS1iYjBlLWNmN2VkNWMyODEyMCIsInN1YiI6ImUyMmNzZXUxMjIwQGJlbm5ldHQuZWR1LmluIn0sImVtYWlsIjoiZTIyY3NldTEyMjBAYmVubmV0dC5lZHUuaW4iLCJuYW1lIjoiYWJoaXJhaiBwYXR3YSIsInJvbGxObyI6ImUyMmNzZXUxMjIwIiwiYWNjZXNzQ29kZSI6InJ0Q0haSiIsImNsaWVudElEIjoiNjBkY2M3NmEtZDg5Zi00NDJhLWJiMGUtY2Y3ZWQ1YzI4MTIwIiwiY2xpZW50U2VjcmV0IjoiS3VEQ0F4SnZiVFRoeUZWRyJ9.qWR2aWponD5V2hS2DJBpCWfpp85VLOmsoHN_tBkP95A"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+authToken)
	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	var numbers []int
	_ = json.Unmarshal(body, &numbers)
	res.numbers = numbers

}
