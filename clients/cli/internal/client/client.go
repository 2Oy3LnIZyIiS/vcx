package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"vcx/pkg/toolkit/printkit"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP:    &http.Client{},
	}
}

type ProgressData struct {
	Step      int    `json:"step"`
	Total     int    `json:"total"`
	Message   string `json:"message"`
	Completed bool   `json:"completed"`
}

func (c *Client) Init() error {
	resp, err := c.HTTP.Get(c.BaseURL + "/api/project/init")
	if err != nil {
		return fmt.Errorf("error calling agent: %w", err)
	}
	defer resp.Body.Close()

    return HandleBody(resp, JustPrint)
}



func HandleBody(resp *http.Response, callback func(string, bool) error ) error {
    contentType := resp.Header.Get("Content-Type")

    switch contentType {
    case "application/json":
        return HandleNonStream(resp.Body)
    case "text/event-stream":
        return HandleStream(resp.Body, callback)
    case "text/plain":
        return HandleNonStream(resp.Body)
    default:
        return HandleNonStream(resp.Body)
    }
}

func HandleStream(respBody io.ReadCloser, callback func(string, bool) error) error {
    scanner := bufio.NewScanner(respBody)
    for scanner.Scan() {
        if callback != nil {
            err := callback(scanner.Text(), false)
            if err != nil {
                return err
            }
        }
    }
    if callback != nil {
        return callback(scanner.Text(), true)
    }

	return scanner.Err()
}


func JustPrint(data string, last bool) error {
    if last {
        return nil
    }
    fmt.Println(data)
    return nil
}


func HandleNonStream(respBody io.ReadCloser) error {
	body, err := io.ReadAll(respBody)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}
    fmt.Println(string(body))

	return nil
}

func (c *Client) InitWithProgress() error {
	resp, err := c.HTTP.Get(c.BaseURL + "/api/project/init-stream")
	if err != nil {
		return fmt.Errorf("error calling agent: %w", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse SSE format: "data: {json}"
		if after, ok := strings.CutPrefix(line, "data: "); ok  {
			jsonData := after

			var progress ProgressData
			if err := json.Unmarshal([]byte(jsonData), &progress); err != nil {
				continue
			}

			if progress.Completed {
				fmt.Println("\n✅ Project initialization completed!")
				break
			}

			// Show progress
            // fmt.Print("\r\033[K")
            printkit.ClearLine()
			fmt.Printf("\r[%d/%d] %s", progress.Step, progress.Total, progress.Message)
		}
	}

	return scanner.Err()
}
