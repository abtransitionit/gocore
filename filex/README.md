# Intro

This package defines a library (i.e. no `main()`) to manage 
- files
- folders

# The code

## `exists.go`

Example usage for a path:
```go
pathFile := "/path/to/a/file"
exists, err := ExistsPath(path)
if err != nil {
    // Unexpected error
    return err
}
if !exists {
    // Path not found (expected case)
    return nil
}
// Path exists — proceed
```

Example usage for a file:
```go
exists, err := ExistsFile(filePath)
if err != nil {
    // Unexpected error, e.g., permission denied
    return err
}
if !exists {
    // File does not exist (expected case)
    return nil
}
// File exists — proceed
...
```

Example usage for a folder:
```go
folderPath := "/tmp/myfolder"
exists, err := ExistsFolder(folderPath)
if err != nil {
    // Handle unexpected or permission errors
    fmt.Printf("Error checking folder: %v\n", err)
    return
}

// Check the boolean result
if exists {
    fmt.Println("Folder exists!")
} else {
    fmt.Println("Folder does not exist.")
}
```






## `touch.go`

# Todo
- use `go-playground/validator` to validate each time parameters/argument provoded to a function

# User action
```go
    currentUser, _ := user.Current()
    if currentUser.Username != "root" {
        fmt.Println("This program must be run as the 'root' user.")
        os.Exit(1)
    }
    // ... rest of the program
```