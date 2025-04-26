package internal

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"
	"time"
)

func RunScan(args []string, cfg *Config) error {
	start := time.Now()
	dir := "."
	var printLock sync.Mutex

	// take the first argument passed in the cmd after "scan"
	if len(args) > 0 {
		dir = args[0]
	}

	// create worker
	filePaths := make(chan FileTask)
	var wg sync.WaitGroup

	for i := 0; i < cfg.NumWorkers; i++ {
		wg.Add(1)
		// create goroutine
		go func(id int) {
			defer wg.Done()
			for task := range filePaths {
				processFile(task, cfg, id, &printLock)
			}
		}(i)
	}

	//Walk each file in the dir
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error walking to %s: %v\n", path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("error getting info for %s: %v", path, err)
		}

		filePaths <- FileTask{Path: path, Info: info} //Send path and info to the channel
		return nil
	})

	if err != nil {
		return fmt.Errorf("scan failed: %v", err)
	}

	close(filePaths) // signal workers to stop
	wg.Wait()        // block main from exiting too early, wait for all workers to finish

	fmt.Printf("Scan completed in %s\n", time.Since(start))
	return nil
}

func processFile(task FileTask, cfg *Config, workerID int, printLock *sync.Mutex) {
	path := task.Path
	info := task.Info
	mode := info.Mode().Perm()
	safePerm, isSensitive := SensitiveFiles[info.Name()]

	if isSensitive {
		if mode != safePerm {
			WithLock(printLock, func() {
				fmt.Printf("Worker %d: %s Sensitive file %s has bad permissions: %04o (expected %04o)\n", workerID, cfg.InsecureTag, path, mode, safePerm)
			})

			if err := FixPermissions(path, safePerm, cfg.FixMode); err != nil {
				WithLock(printLock, func() {
					fmt.Printf("%s Failed to fix permissions for %s: %v\n", cfg.InsecureTag, path, err)
				})
			}
		} else if !cfg.InsecureOnly {
			WithLock(printLock, func() {
				fmt.Printf("Worker %d: %s Sensitive file %s is secure - %04o\n", workerID, cfg.SecureTag, path, mode)
			})
		}
		return
	}

	// check if it is world writable
	// 0o002 octal -> 000 000 010 binary -> ------w-
	if mode&0o002 != 0 {
		WithLock(printLock, func() {
			fmt.Printf("Worker %d: %s Insecure permissions on %s - %04o\n", workerID, cfg.InsecureTag, path, mode)
		})

		if err := FixPermissions(path, 0644, cfg.FixMode); err != nil {
			WithLock(printLock, func() {
				fmt.Printf("%s Failed to fix permissions for %s: %v\n", cfg.InsecureTag, path, err)
			})
		}
	} else if !cfg.InsecureOnly {
		WithLock(printLock, func() {
			fmt.Printf("Worker %d: %s %s - %04o\n", workerID, cfg.SecureTag, path, mode)
		})
	}
}
