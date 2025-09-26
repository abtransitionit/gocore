# Intro

This package defines utility to run task concurrently 


# The wait function
```go
// example to test the wait using a closure because wait want a "func() (bool error)"
func CheckVpsIsReady(start time.Time) func() (bool, error) {
	return func() (bool, error) {
		elapsed := time.Since(start)
		fmt.Printf("checking resource after %v...\n", elapsed.Truncate(time.Second))
		if elapsed >= 10*time.Second {
			return true, nil
		}
		return false, nil
	}
}
```

```go
// define a specific timeout - example when rebuilding the VM
waitTimeout := 3 * time.Second
waitCtx, cancel := context.WithTimeout(ctx, waitTimeout)
defer cancel() // always defer cancel to release resources

// create the wait for status function 
start := time.Now()
healthFunc := ovh.CheckVpsIsReady(start)

// launch the waiting
fmt.Println("Waiting for resource to be ready...")
err := syncx.WaitForReady(waitCtx, logger, 2*time.Second, healthFunc)
if err != nil {
  fmt.Printf("Resource not ready: %v\n", err)
  return
}
fmt.Println("Resource is ready! ✅")
```
