package api

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
)

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
