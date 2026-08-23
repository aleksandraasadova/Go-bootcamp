package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	timeoutFlag int
	helpFlag    bool
)

func init() {
	flag.IntVar(&timeoutFlag, "t", 15, "timeout")
	flag.IntVar(&timeoutFlag, "timeout", 15, "timeout")
	flag.BoolVar(&helpFlag, "h", false, "show help")
	flag.BoolVar(&helpFlag, "help", false, "show help")
}

type getResponse struct {
	url     string
	status  int
	headers http.Header
	body    string
	err     error
}

func main() {
	flag.Parse()
	args := flag.Args()

	if helpFlag {
		fmt.Println("The usage of hedgedcurl util is: ./hedgedcurl [OPTIONS] URL1 URL2 ...")
		os.Exit(0)
	}

	if len(args) == 0 {
		fmt.Println("There were no URLs provided for hedgedcurl util")
		os.Exit(1)
	}

	result := make(chan *getResponse)

	for _, url := range args {
		go func(targetURL string) {
			client := &http.Client{}
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutFlag)*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
			if err != nil {
				result <- &getResponse{url: targetURL, err: err}
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				result <- &getResponse{url: targetURL, err: err}
				return
			}
			defer resp.Body.Close()
			readBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				result <- &getResponse{url: targetURL, err: err}
				return
			}
			result <- &getResponse{url: targetURL, status: resp.StatusCode, headers: resp.Header, body: string(readBytes)}
		}(url)
	}

	for range args {
		first := <-result
		if first.err != nil {
			if os.IsTimeout(first.err) {
				fmt.Printf("reached timeout: %s\n", first.url)
				os.Exit(228)
			}
			continue
		}
		fmt.Println(first.status)
		fmt.Println(first.headers)
		fmt.Println(first.body)
		os.Exit(0)
	}
	fmt.Println("Couldn't get any of URLs")
	os.Exit(1)
}
