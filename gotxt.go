package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// printDirectoryTree prints the directory structure to the output writer
func printDirectoryTree(dir string, output *bufio.Writer, exceptions, fileTypes []string, level int) error {
	if level == 0 {
		_, err := fmt.Fprintln(output, filepath.Base(dir))
		if err != nil {
			return err
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Name() == "sum.txt" {
			continue // Skip the output file itself
		}
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() && !contains(exceptions, entry.Name()) {
			// Print directory name with appropriate indentation
			_, err := fmt.Fprintf(output, "%s├── %s/\n", strings.Repeat("│   ", level), entry.Name())
			if err != nil {
				return err
			}
			// Recursively print subdirectory contents
			err = printDirectoryTree(path, output, exceptions, fileTypes, level+1)
			if err != nil {
				return err
			}
		} else if !entry.IsDir() && (contains(fileTypes, "all") || contains(fileTypes, filepath.Ext(entry.Name())[1:])) {
			// Print file name with appropriate indentation
			_, err := fmt.Fprintf(output, "%s├── %s\n", strings.Repeat("│   ", level), entry.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// collectFileContents collects the contents of files in the specified directory and writes them to the output file
func collectFileContents(dir, outputFile string, exceptions, fileTypes []string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Print directory tree
	err = printDirectoryTree(dir, writer, exceptions, fileTypes, 0)
	if err != nil {
		return err
	}

	// Write separator
	_, err = fmt.Fprintf(writer, "\n%s\n\n", strings.Repeat("=", 50))
	if err != nil {
		return err
	}

	// Walk through directory and collect file contents
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && contains(exceptions, info.Name()) {
			return filepath.SkipDir
		}
		if !info.IsDir() && info.Name() != "sum.txt" &&
			!contains(exceptions, info.Name()) &&
			(contains(fileTypes, "all") || contains(fileTypes, filepath.Ext(info.Name())[1:])) {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(writer, "File: %s\n\n", relPath)
			if err != nil {
				return err
			}
			err = writeFileContents(path, writer)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(writer, "\n%s\n\n", strings.Repeat("-", 50))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// writeFileContents writes the contents of a file to the output writer
func writeFileContents(path string, writer *bufio.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(writer, file)
	return err
}

// contains checks if a string is present in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: gotxt <exceptions> <file_types>")
		fmt.Println("  <exceptions>   : Exception files/directories (comma-separated)")
		fmt.Println("  <file_types>   : File types to include (comma-separated) or 'all' to include all file types")
		os.Exit(1)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}

	outputFile := "sum.txt"
	exceptions := append(strings.Split(os.Args[1], ","), "sum.txt")
	fileTypes := strings.Split(os.Args[2], ",")

	// Clean up user input
	for i := range exceptions {
		exceptions[i] = strings.TrimSpace(exceptions[i])
	}
	for i := range fileTypes {
		fileTypes[i] = strings.TrimSpace(fileTypes[i])
	}

	err = collectFileContents(dir, outputFile, exceptions, fileTypes)
	if err != nil {
		fmt.Println("Error collecting file contents:", err)
		os.Exit(1)
	}

	fmt.Printf("File contents collected and saved to %s\n", outputFile)
}