package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
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

func dirTree(out io.Writer, path string, files bool) interface{} {
	err := printDir(out, path, "", files)
	if err != nil {
		panic(err)
	}
	return nil
}

//func getDirFiles(out *os.File, prefix, pwd string, printFiles bool) {
//	files, err := ioutil.ReadDir(pwd)
//	if err != nil {
//		fmt.Println("Panic happend:", err)
//	}
//	if !printFiles {
//		printOnlyDir := []os.FileInfo{}
//		for _, file := range files {
//			if file.IsDir() {
//				printOnlyDir = append(printOnlyDir, file)
//			}
//		}
//		files = printOnlyDir
//	}
//	length := len(files)
//	for i, file := range files {
//		if file.Name()[0] == '.' {
//			continue
//		} else if file.IsDir() {
//			var prefixNew string
//			if length > i+1 {
//				fmt.Fprintf(out, prefix+"├───%s\n", file.Name())
//				prefixNew = prefix + "│\t"
//			} else {
//				fmt.Fprintf(out, prefix+"└───%s\n", file.Name())
//				prefixNew = prefix + "\t"
//			}
//			getDirFiles(out, prefixNew, filepath.Join(pwd, file.Name()), printFiles)
//		} else if printFiles {
//			if file.Size() > 0 {
//				if length > i+1 {
//					fmt.Fprintf(out, prefix+"├───%s (%vb)\n", file.Name(), file.Size())
//				} else {
//					fmt.Fprintf(out, prefix+"└───%s (%vb)\n", file.Name(), file.Size())
//				}
//			} else {
//				if length > i+1 {
//					fmt.Fprintf(out, prefix+"├───%s (empty)\n", file.Name())
//				} else {
//					fmt.Fprintf(out, prefix+"└───%s (empty)\n", file.Name())
//				}
//			}
//		}
//	}
//}

func printDir(out io.Writer, path, prefix string, printFiles bool) error {
	f, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if !printFiles {
		onlyDir := []os.DirEntry{}
		for _, dir := range f {
			if dir.IsDir() {
				onlyDir = append(onlyDir, dir)
			}
		}
		f = onlyDir
	}

	length := len(f)
	for i, file := range f {
		if file.Name()[0] == '.' {
			continue
		}

		if file.IsDir() {
			var newPrefix string
			if i == length-1 {
				_, err = fmt.Fprintf(out, prefix+endSim+"%s\n", file.Name())
				newPrefix = prefix + "\t"
				if err != nil {
					return err
				}
			} else {
				_, err = fmt.Fprintf(out, prefix+simDir+"%s\n", file.Name())
				newPrefix = prefix + "│\t"
				if err != nil {
					return err
				}
			}

			err = printDir(out, path+"/"+file.Name(), newPrefix, printFiles)
			if err != nil {
				return err
			}
		} else {
			info, err := file.Info()
			if err != nil {
				return err
			}

			size := strconv.Itoa(int(info.Size())) + "b"
			if info.Size() == 0 {
				size = "empty"
			}

			if i != length-1 {
				_, err = fmt.Fprintf(out, "%s%s (%s)\n", prefix+simDir, file.Name(), size)
			} else {
				_, err = fmt.Fprintf(out, "%s%s (%s)\n", prefix+endSim, file.Name(), size)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
