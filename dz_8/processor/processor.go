package processor

import "sync"

func ProcessFiles(paths []string, workers int) <-chan Result {
	if workers < 1 {
		workers = 1
	}

	jobs := make(chan string)
	results := make(chan Result)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				results <- analyzeFile(path)
			}
		}()
	}

	go func() {
		for _, path := range paths {
			jobs <- path
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
