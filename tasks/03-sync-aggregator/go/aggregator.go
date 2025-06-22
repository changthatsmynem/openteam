// Package aggregator – stub for Concurrent File Stats Processor.
package aggregator

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Result mirrors one JSON object in the final array.
type Result struct {
	Path   string `json:"path"`
	Lines  int    `json:"lines,omitempty"`
	Words  int    `json:"words,omitempty"`
	Status string `json:"status"` // "ok" or "timeout"
}

// helper struct to help wraps a Result data and its original index from the input list
// so later can sequently return the results in the same order
type resultmodel struct {
	index  int
	result Result
}

// the index is used to maintain the order of results
// the originPath is the path as it appears in the filelist
// the fullPath is the absolute path to the file so can read the file
type job struct {
	index      int
	originPath string
	fullPath   string
}

// Aggregate must read filelistPath, spin up *workers* goroutines,
// apply a per‑file timeout, and return results in **input order**.
func Aggregate(filelistPath string, workers, timeout int) ([]Result, error) {
	//open the file
	f, err := os.Open(filelistPath)
	if err != nil {
		return nil, err
	}
	//then close it later after everything is done
	defer f.Close()

	//then scan using bufio to read the file line by line
	//and will store each file paths in a slice of job structs
	//use the filepath lib to join the directory of the filelist with the line read from the filelist
	//so can get the full path
	var jobsList []job
	scanner := bufio.NewScanner(f)
	index := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			fullPath := filepath.Join(filepath.Dir(filelistPath), line)
			jobsList = append(jobsList, job{
				index:      index,
				originPath: line,
				fullPath:   fullPath,
			})
			index++
		}
	}

	//create channels for jobs and results for workers
	jobs := make(chan job)
	results := make(chan resultmodel)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//wait for job to be sent to the channel, then execute
			for job := range jobs {
				//creates a context that will auto-cancel after timeout
				//if timeout happens before processing completes, ctx.Done() will signals a timeout
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
				//pass a context.Context so we can set "timeout" status if it runs too long (exceeds timeout)
				result := processFile(ctx, job.fullPath)
				//set the path to original path to match the expected output
				result.Path = job.originPath
				//send result to the results channel with the original index, so can restore order
				results <- resultmodel{index: job.index, result: result}
				//call cancel() to release resources used by the context after processing each job is done
				cancel()
			}
		}()
	}

	//send each job into the jobs channel so workers can pick them up
	go func() {
		for _, j := range jobsList {
			jobs <- j
		}
		close(jobs)
	}()

	//collect a result for each input file
	//uses the saved index for the correct position
	finalResults := make([]Result, len(jobsList))
	for range jobsList {
		res := <-results
		finalResults[res.index] = res.result
	}

	//wait for all the workers are done (wg.Done() that defer in routine) then return

	wg.Wait()

	return finalResults, nil
}

func processFile(ctx context.Context, path string) Result {
	//open each file
	f, err := os.Open(path)
	if err != nil {
		return Result{Path: path, Status: "error"}
	}
	defer f.Close()

	//reading line by line
	scanner := bufio.NewScanner(f)

	lines := 0
	words := 0

	for scanner.Scan() {
		// check ctx before scanning next line to avoid blocking after timeout (timeout reached)
		select {
		case <-ctx.Done():
			return Result{Path: path, Status: "timeout"}
		default:
		}

		line := scanner.Text()
		//if encounter a line starting with #sleep=, sleep for that many seconds
		//and don't count the #sleep= lines and words
		if strings.HasPrefix(line, "#sleep=") {
			if n, err := strconv.Atoi(strings.TrimPrefix(line, "#sleep=")); err == nil && n > 0 {
				//use a select to wait for either the context to be done or the sleep to finish for handle the timeout case
				//if ctx is done or canceled before sleep finishes, return a timeout status
				//if sleep finishes, continue
				select {
				case <-ctx.Done():
					return Result{Path: path, Status: "timeout"}
				case <-time.After(time.Duration(n) * time.Second):
				}
			}
			continue
		}

		//count lines and words that is not a #sleep= marker
		lines++
		words += len(strings.Fields(line))
	}

	if err := scanner.Err(); err != nil {
		return Result{Path: path, Status: "error"}
	}

	return Result{
		Path:   path,
		Lines:  lines,
		Words:  words,
		Status: "ok",
	}
}
