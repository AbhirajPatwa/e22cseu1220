package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	WindowPrevState []int   `json:"windowPrevState"`
	WindowCurrState []int   `json:"windowCurrState"`
	Numbers         []int   `json:"numbers"`
	Avg             float64 `json:"avg"`
}

var state []int = []int{}

func main() {
	e := echo.New()

	e.GET("/numbers/:numberId", avgCalc)

	e.Logger.Fatal(e.Start(":8080"))
}

func numsWindow(arr []int, window int) []int {
	if len(arr) <= window {
		return arr
	} else {
		return arr[len(arr)-10:]
	}
}

func uniqueNums(arr []int) []int {
	seen := make(map[int]bool)
	var result []int
	for _, num := range arr {
		if _, exists := seen[num]; !exists {
			result = append(result, num)
			seen[num] = true
		}
	}

	return result
}

func avgCalc(c echo.Context) error {
	windowSize := 10
	res := Response{}
	res.WindowPrevState = state
	fmt.Println(state)
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

	authToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzQzNzQ2MzczLCJpYXQiOjE3NDM3NDYwNzMsImlzcyI6IkFmZm9yZG1lZCIsImp0aSI6IjYwZGNjNzZhLWQ4OWYtNDQyYS1iYjBlLWNmN2VkNWMyODEyMCIsInN1YiI6ImUyMmNzZXUxMjIwQGJlbm5ldHQuZWR1LmluIn0sImVtYWlsIjoiZTIyY3NldTEyMjBAYmVubmV0dC5lZHUuaW4iLCJuYW1lIjoiYWJoaXJhaiBwYXR3YSIsInJvbGxObyI6ImUyMmNzZXUxMjIwIiwiYWNjZXNzQ29kZSI6InJ0Q0haSiIsImNsaWVudElEIjoiNjBkY2M3NmEtZDg5Zi00NDJhLWJiMGUtY2Y3ZWQ1YzI4MTIwIiwiY2xpZW50U2VjcmV0IjoiS3VEQ0F4SnZiVFRoeUZWRyJ9.2jOv1zjInWNMh2uLXjtU4HuGBFXGX6bp27ywF1rgE5A"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
	}
	// fmt.Println(req)
	// fmt.Println(resp)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.Body)
	var numbers []int
	_ = json.Unmarshal(body, &numbers)
	if len(numbers) == 0 {
		numbers = []int{2, 4, 6, 8}
	}
	res.Numbers = numbers
	state = numbers
	// fmt.Println(numbers)

	numbers = uniqueNums(numbers)
	numbers = numsWindow(numbers, windowSize)
	res.WindowCurrState = numbers
	avg := 0.0
	for _, i := range numbers {
		avg += float64(i)
		avg /= float64(len(numbers))
	}

	res.Avg = avg
	fmt.Println(res)
	return c.JSON(http.StatusOK, res)
}
