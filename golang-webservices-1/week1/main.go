package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func getDirs(files []os.FileInfo) []os.FileInfo {
	dirs := make([]os.FileInfo, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}
	return dirs
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return printTree(out, path, printFiles, "")
}

func printTree(out io.Writer, path string, printFiles bool, indent string) error {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}

	fileInfoSl, err := f.Readdir(-1)
	if err != nil {
		return err
	}

	if !printFiles {
		fileInfoSl = getDirs(fileInfoSl)
	}

	sort.Slice(fileInfoSl, func(i, j int) bool {
		return fileInfoSl[i].Name() < fileInfoSl[j].Name()
	})

	var nodePrefix, newIndent string
	for idx, file := range fileInfoSl {
		if idx == len(fileInfoSl)-1 {
			nodePrefix = "└───"
			newIndent = indent + "\t"
		} else {
			nodePrefix = "├───"
			newIndent = indent + "│\t"
		}

		if file.IsDir() {
			fmt.Fprintf(out, "%s%s%s\n", indent, nodePrefix, file.Name())
			err = printTree(out, filepath.Join(path, file.Name()), printFiles, newIndent)
			if err != nil {
				return err
			}
		} else {
			var sizeStr string
			if file.Size() == 0 {
				sizeStr = "empty"
			} else {
				sizeStr = fmt.Sprintf("%db", file.Size())
			}
			fmt.Fprintf(out, "%s%s%s (%s)\n", indent, nodePrefix, file.Name(), sizeStr)
		}
	}
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
