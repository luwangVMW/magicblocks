package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    // Define file paths
    folder1 := "/data/folder1/myfile1.txt" // Volume mount
    folder2 := "/data/folder2/myfile2.txt" // Local directory

    // Write to folder1 (volume mount)
    err := writeFile(folder1, "This is a test file in folder1.")
    if err != nil {
        fmt.Printf("Error writing to folder1: %v\n", err)
    } else {
        fmt.Println("Successfully wrote to folder1.")
    }

    // Write to folder2 (local directory)
    err = writeFile(folder2, "This is a test file in folder2.")
    if err != nil {
        fmt.Printf("Error writing to folder2: %v\n", err)
    } else {
        fmt.Println("Successfully wrote to folder2.")
    }
    duration := 60 * time.Minute
    time.Sleep(duration)
}

// Function to write a string to a file
func writeFile(filepath string, content string) error {
    // Create or open the file
    file, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write content to the file
    _, err = file.WriteString(content)
    if err != nil {
        return err
    }

    return nil
}
