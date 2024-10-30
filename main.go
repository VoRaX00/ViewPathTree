package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err)
	}
}

const (
	simDir = "├───"
	endSim = "└───"
)

func dirTree(out *os.File, path string, files bool) interface{} {
	_, err := out.WriteString(simDir + "project\n")
	if err != nil {
		panic(err)
	}

	err = printDir(out, path, 1)
	if err != nil {
		panic(err)
	}
	return nil
}

func printDir(out *os.File, path string, cntSpaces int) error {
	f, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, file := range f {
		if file.Name()[0] == '.' {
			continue
		}
		spaces := strings.Repeat("│"+"   ", cntSpaces)

		if i == len(f)-1 {
			_, err = out.WriteString(spaces + endSim + file.Name())
		} else {
			_, err = out.WriteString(spaces + simDir + file.Name())
		}

		if file.IsDir() {
			out.WriteString("\n")
			err = printDir(out, path+"/"+file.Name(), cntSpaces+1)
			if err != nil {
				return err
			}
		} else {
			info, err := file.Info()
			if err != nil {
				return err
			}

			size := strconv.Itoa(int(info.Size())) + "b"
			if size == "0b" {
				size = "empty"
			}
			out.WriteString(fmt.Sprintf(" (%s)\n", size))
		}
	}
	return nil
}
