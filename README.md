# GoTxt

GoTxt is a command-line tool written in Go that generates a text file containing a directory tree structure and the contents of files within that directory. It's designed to be a faster alternative to similar tools written in other languages.

## Features

- Generates a directory tree structure
- Collects contents of files in the directory
- Allows specifying exceptions (files/directories to exclude)
- Supports filtering by file types
- Outputs results to a file named `sum.txt`

## Installation

1. Ensure you have Go installed on your system. If not, download and install it from [golang.org](https://golang.org/).

2. Clone this repository.

3. Add the directory containing `gotxt.go` and `gotxt.bat` to PATH:

## Usage

Run the tool from any directory using the following command:

```
gotxt <exceptions> <file_types>
```

- `<exceptions>`: Comma-separated list of files/directories to exclude
- `<file_types>`: Comma-separated list of file extensions to include, or 'all' for all file types

Example:
```
gotxt node_modules,.git,vendor txt,md,go
```

This command will:
- Exclude the `node_modules`, `.git`, and `vendor` directories
- Include only `.txt`, `.md`, and `.go` files
- Generate a `sum.txt` file in the current directory with the results

## Output

The tool generates a file named `sum.txt` in the current directory, containing:
1. A tree structure of the directory
2. A separator line
3. The contents of each included file, preceded by its relative path

## License

This project is licensed under the MIT License.