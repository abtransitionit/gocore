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
Example usage 1
```go
// Assume the file not exists
filePath := "/path/to/my-new-file.txt"
created, err := filex.Touch(filePath)
if err != nil {
	log.Fatalf("unexpected error: %v", err)
}
fmt.Printf("File was created: %v\n", created)
// Expected output: File was created: true
```

Example usage 2
```go
// Assume a file already exists at this path
filePath := "/path/to/existing-file.txt"
created, err := filex.Touch(filePath)
if err != nil {
	log.Fatalf("unexpected error: %v", err)
}
fmt.Printf("File was created: %v\n", created)
// Expected output: File was created: false
```

Example usage 3
```go
// Assume this path is a directory
dirPath := "/path/to/my-directory"
created, err := filex.Touch(dirPath)
if err != nil {
	fmt.Printf("Error: %v\n", err)
}
fmt.Printf("File was created: %v\n", created)
// Expected output:
// Error: path is a directory, not a file at /path/to/my-directory
// File was created: false
```

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