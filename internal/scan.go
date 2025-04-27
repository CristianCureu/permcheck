package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

func RunScan(args []string, cfg *Config) error {
	start := time.Now()
	cfg.InsecureTag = FormatTag(cfg.InsecureTag, ColorRed, cfg.ColorsEnabled)
	cfg.SecureTag = FormatTag(cfg.SecureTag, ColorGreen, cfg.ColorsEnabled)
	dir := "."
	var printLock sync.Mutex

	// take the first argument passed in the cmd after "scan"
	if len(args) > 0 {
		dir = args[0]
	}

	// create worker
	filePaths := make(chan FileTask)
	var wg sync.WaitGroup

	// counters
	var filesScanned int64
	var filesSecure int64
	var filesInsecure int64

	for i := 0; i < cfg.NumWorkers; i++ {
		wg.Add(1)
		// create goroutine
		go func(id int) {
			defer wg.Done()
			for task := range filePaths {
				scanned, secure, insecure := processFile(task, cfg, id, &printLock)
				addCounters(&filesScanned, &filesSecure, &filesInsecure, scanned, secure, insecure)
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

	fmt.Printf("\nScan completed in %s\n", time.Since(start))

	fmt.Println("+----------------+--------+")
	fmt.Println("|     Result     | Count  |")
	fmt.Println("+----------------+--------+")

	if !cfg.InsecureOnly {
		fmt.Printf("| %s     | %6d |\n", FormatTag("[✓] Secure", ColorGreen, cfg.ColorsEnabled), filesSecure)
	}

	fmt.Printf("| %s   | %6d |\n", FormatTag("[!] Insecure", ColorRed, cfg.ColorsEnabled), filesInsecure)
	fmt.Printf("| Total scanned  | %6d |\n", filesScanned)
	fmt.Println("+----------------+--------+")

	return nil
}

func processFile(task FileTask, cfg *Config, workerID int, printLock *sync.Mutex) (scanned, secure, insecure int64) {
	path := task.Path
	info := task.Info
	mode := info.Mode().Perm()

	scanned = 1

	var matched bool
	var expectedPerm os.FileMode
	var reason string

	for _, rule := range Rules {
		match, expected := rule.Match(info)
		if match {
			matched = true
			expectedPerm = expected
			reason = rule.Name
			break // stop at first matching rule
		}
	}

	if !matched {
		// no rule matched
		expectedPerm = mode
		reason = "Generic file"

		if mode&0o002 != 0 { // if it is world writable
			expectedPerm = 0644
			reason = "World-writable file"
		}
	}

	if mode != expectedPerm {
		insecure = 1
		WithLock(printLock, func() {
			fmt.Printf("Worker %d: %s %s: %s has bad permissions: %04o (expected %04o)\n",
				workerID, cfg.InsecureTag, reason, path, mode, expectedPerm)
		})
		fixAndReport(path, expectedPerm, cfg, printLock)
	} else if !cfg.InsecureOnly {
		secure = 1
		WithLock(printLock, func() {
			fmt.Printf("Worker %d: %s %s - %04o\n",
				workerID, cfg.SecureTag, path, mode)
		})
	}

	return
}

func fixAndReport(path string, desiredPerm os.FileMode, cfg *Config, printLock *sync.Mutex) {
	if err := FixPermissions(path, desiredPerm, cfg.FixMode); err != nil {
		WithLock(printLock, func() {
			fmt.Printf("%s Failed to fix permissions for %s: %v\n", cfg.InsecureTag, path, err)
		})
	}
}

func addCounters(scannedTotal, secureTotal, insecureTotal *int64, scanned, secure, insecure int64) {
	atomic.AddInt64(scannedTotal, scanned)
	atomic.AddInt64(secureTotal, secure)
	atomic.AddInt64(insecureTotal, insecure)
}
