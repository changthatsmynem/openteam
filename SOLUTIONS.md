### Solution Notes

### Task 01 – Run‑Length Encoder
- Language: Go
- Approach: 
    1: if it's blank "" return ""
    2: convert an input to a slice of rune to handle multi-byte chars like emoji
    3: declare count varaible - starting count from 1 because every chars is at least appeared once
    4: declare a result buffer name "result" using strings.Builder
    5: iterate over rune, starting from the second char to compare with previous one
    5.1: if the current char is same as previous one, increment the count
    5.2: if the current char is different, append the previous character and its count to result, then reset count for the 6: return the result with concatenate the last character and its count
- Why: This is the solution that only works if the string is a consecutive order for example like 
- Time spent: ~15 min
- AI tools used: Copilot

### Task 02 – Fix‑the‑Bug
- Language: Go
- Approach: 
  The problem lies on goroutines calling the NextID function concurrently, which generates a unique ID by incrementing a current number without a thread-safe so "race condition" might occur becasue multiple goroutines will read and write to the shared variable "current" at the same time, leading to duplicate IDs.
  To fix this, I use a "mutex" to ensure a thread-safe that only one goroutine can access the variable at a time
- Why: the mutex ensures that the race condition will not happen especially a case like this, an access on the global variable
- Time spent: ~10 min
- AI tools used: Copilot on writing a comment
  
### Task 03 - sync-aggregator
- Language: Go
- Approach: 
	- open the file and then close it later after everything is done then scan using bufio to read the file line by line and store each file paths in a slice of job structs which has the index to maintain the order of results, the originPath is the path as it appears in the filelist, the fullPath is the absolute path to the file so can read the file
    - use the filepath lib to join the directory of the filelist with the line read from the filelist
	so can get the full path
    - create channels for jobs and results for workers (use helper struct to help wraps a Result data and its original index from the input list so later can sequently return the results in the same order)
    - send each job into the jobs channel so workers can pick them up
    - wait for job to be sent to the channel, then execute
    - creates a context that will auto-cancel after timeout, so if timeout happens before processing completes, ctx.Done() will signals a timeout
	- pass a context.Context to helper function processFile to process each file path and so we can set "timeout" status if it runs too long (exceeds timeout)
        (inside helper function)
    	- always check ctx before scanning next line to avoid blocking after timeout (timeout reached)
    	- if encounter a line starting with #sleep=, sleep for that many seconds and don't count the #sleep= lines and words
    	- use a select to wait for either the context to be done or the sleep to finish for handle the timeout case if ctx is done or canceled before sleep finishes, return a timeout status if sleep finishes, continue
    	- then count lines and words that is not a #sleep= marker
	- send result to the results channel with the original index, so can restore order and set the path to original path to match the expected output
	- call cancel() to release resources used by the context after processing each job is done
	- collect a result for each input file and uses the saved index for the correct position and wait for all the workers are done (wg.Done() that defer in routine) then return the output
- Why: The code still did not pass the test case but i will elaborate why i use this method: i use jobs as a channel to create concurrency flow allows multiple workers to process files independently and concurrently, pulling from the same pool of tasks and as soon as a job becomes available, any free worker can pick it up and start the process
- Time spent: ~3 hours
- AI tools used: Copilot and ChatGPT

### Task 04 - sql-reasoning
- Language: Go
- Approach: 
    A: use GROUP BY to group each id and then use aggregate SUM()
    B: use CTE to group a scope of 'global' and 'thailand' then use ROW_NUMBER() to order pledges in ascending by amount then COUNT(*) total pledges in that scope, then next CTE for calculate percentile: use CAST((0.9 * total + 0.9999999) AS INT) to simulate the CEIL function and then select the amount at that percentile
- Why: on Task A it is quite obvious: aggregates pledge amounts per campaign / on Task B we need to use window function to calculate the nearest-rank rule (1-based) so in the next cte for each scope just group 
- Time spent: ~50 min
- AI tools used: ChatGPT to optimize