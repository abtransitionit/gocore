# Intro

This package defines a library (i.e. no `main()`) to manage 
- files
- folders

# The code

## `exists.go`
Usage for file:
```go
pathFile := "/path/to/a/file"
exists, err := ExistsFile(pathFile)
if err != nil {
    // Check if the error is specific: a permissions error
    if os.IsPermission(err) {
        fmt.Println("Error: Not enough permissions to access the file.")
    } else {
        fmt.Println("An unexpected error occurred:", err)
    }
    return
}
```

Usage for folder:
```go
pathDir := "/path/to/a/folder"	
isFolder, err := ExistsFolder(pathDir)
if err != nil {
    // Check if the error is specific: a permissions error
		if os.IsPermission(err) {
			fmt.Println("Error: Not enough permissions to access the path.")
		} else {
			fmt.Println("An unexpected error occurred:", err)
		}
		return
	}
	
	if isFolder {
		fmt.Printf("The path '%s' is an existing folder.\n", pathDir)
	} else {
		fmt.Printf("The path '%s' is not an existing folder.\n", pathDir)
	}

	// Now let's try a path that is a file
	filePath := "/path/to/my/file.txt"
	isFolder, err = ExistsFolder(filePath)

	if err != nil {
		// Error handling for the file path
		if os.IsPermission(err) {
			fmt.Println("Error: Not enough permissions to access the file.")
		} else {
			fmt.Println("An unexpected error occurred:", err)
		}
		return
	}

	if isFolder {
		fmt.Printf("The path '%s' is an existing folder.\n", filePath)
	} else {
		fmt.Printf("The path '%s' is not an existing folder.\n", filePath)
	}
}
```go





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