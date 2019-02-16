package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

//var recursive = flag.Bool("s", false, "Recurse subdirectories")
var help = flag.Bool("h", false, "Help")
var excludeDirectories = flag.Bool("xd", false, "Exclude Directories")
var excludeFiles = flag.Bool("xf", false, "Exclude Files")
var excludeSymlinks = flag.Bool("xs", false, "Exclude Symlinks")
var textSearch = flag.String("f", "", "Text Search")

type FileList struct {
	name    string
	size    string
	modTime string
	isDir   bool
	symlink bool
}

func GetCwd() string {
	dir, _ := os.Getwd()
	return dir
}

func SymlinkCheck(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		//log.Fatal(err) // if we have this enabled it raises access errors on some files unless ll is run as admin.
		return false
	}
	switch mode := fi.Mode(); {
	case mode&os.ModeSymlink != 0:
		return true
	}
	return false
}

func SizeCommaed(oldSize int64) string {
	size := strconv.Itoa(int(oldSize))
	if len(size) < 4 {
		return size
	}
	var commaed string
	index := 1
	for i := len(size) - 1; i >= 0; i-- {
		if index == 4 {
			commaed = "," + commaed
			index = 1
		}
		commaed = string(size[i]) + commaed
		index++
	}
	return commaed
}

func HelpOutput() {
	fmt.Println("Lists all files in a directory")
	fmt.Println()
	fmt.Println("ll <commands> <path to list>")
	fmt.Println()
	fmt.Println("  <Comamnds>")
	//fmt.Println("  -s : Include subdirectories")
	//fmt.Println("  -sp : Pause after each directory")
	fmt.Println("  -xf : Exclude files")
	fmt.Println("  -xd : Exclude Directories")
	fmt.Println("  -xs : Exclude Symlinks")
	fmt.Println("  -h  : Help menu")
}

func pause() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func ListPath(workingPath string) int64 {
	var largestSize int = 4
	var totalSize int64 = 0
	var size string
	var totalFiles int = 0
	var totalDirs int = 0
	var totalLinks int = 0
	var bonusSpacing int = 2
	var workingPathTarget string = ""

	if SymlinkCheck(workingPath) == true {
		linkPath, err := os.Readlink(workingPath)
		if err != nil {
			fmt.Println(workingPath + " is a symlink but we failed to get the source path.")
			os.Exit(1)
		}
		workingPathTarget = linkPath
	}

	files, err := ioutil.ReadDir(workingPath)
	if err != nil {
		fmt.Println(workingPath + " is not a valid directory.")
		os.Exit(1)
	}

	storage := map[FileList]bool{}
	for _, f := range files {
		if strings.Contains(strings.ToLower(f.Name()), strings.ToLower(*textSearch)) {
			if f.IsDir() {
				if *excludeDirectories == true {
					continue
				}
			} else {
				if SymlinkCheck(path.Join(workingPath, f.Name())) == true && *excludeSymlinks == true {
					continue
				} else if *excludeFiles == true {
					continue
				}
			}
			s := new(FileList)
			s.name = f.Name()
			s.size = SizeCommaed(f.Size())
			totalSize += f.Size()
			if len(s.size) > largestSize {
				largestSize = len(s.size)
			}
			s.modTime = f.ModTime().String()[:19]
			s.isDir = f.IsDir()
			if s.isDir == true {
				totalDirs += 1
			} else if f.Size() == 0 {
				if SymlinkCheck(path.Join(workingPath, f.Name())) == true {
					s.symlink = true
					totalLinks += 1
					bonusSpacing = 7
				}
			} else {
				totalFiles += 1
			}
			storage[*s] = true
		}
	}

	fmt.Println("")
	for i := range storage {
		if i.isDir == true {
			size = "<DIR>" + strings.Repeat(" ", largestSize+bonusSpacing)
		} else if i.symlink == true {
			size = "<JUNCTION>" + strings.Repeat(" ", largestSize+bonusSpacing-5)
			link, err := os.Readlink(path.Join(workingPath, i.name))
			if err != nil {
				continue
			}
			i.name = i.name + " => " + link
		} else {
			size = strings.Repeat(" ", largestSize-len(i.size)+len("<DIR>")+bonusSpacing) + i.size
		}
		fmt.Println(i.modTime + "  " + size + "  " + i.name)
	}
	fmt.Println("")
	pathString := strings.Repeat(" ", 14) + "Path\t" + workingPath
	if len(workingPathTarget) > 0 {
		pathString = pathString + " => " + workingPathTarget
	}
	fmt.Println(pathString)
	if totalDirs > 0 {
		fmt.Println(strings.Repeat(" ", 14) + strconv.Itoa(totalDirs) + " Dir(s)")
	}
	if totalLinks > 0 {
		fmt.Println(strings.Repeat(" ", 14) + strconv.Itoa(totalLinks) + " Symlink(s)")
	}
	if totalFiles > 0 {
		fmt.Println(strings.Repeat(" ", 14) + strconv.Itoa(totalFiles) + " File(s)\t\t" + SizeCommaed(totalSize) + " bytes")
	}
	return totalSize
}

func GetDirectories(workingPath string) []string {
	directories := []string{}
	// if *recursive == true {
	//	filepath.Walk(workingPath, func(path string, f os.FileInfo, err error) error {
	//		if f.IsDir() == true {
	//			directories = append(directories, path.Join(workingPath, f.Name()))
	//		}
	//		return nil
	//	})
	//} else {
	//	directories = append(directories, workingPath)
	//}
	directories = append(directories, workingPath) // temp thing until we get the recursion sorted out.
	return directories
}

func DoesPathExist(workingPath string) bool {
	_, err := os.Stat(workingPath)
	if err == nil {
		return true
	} else {
		print("The path does not exist...\n")
		os.Exit(0)
	}
	return false
}

func RemoveArgs() []string {
	var newArgs []string
	last := "empty" // need a better solution for this.
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-") == false {
			if strings.HasPrefix(last, "-f") == false {
				newArgs = append(newArgs, arg)
			}
		}
		last = arg
	}
	return newArgs
}

func main() {
	flag.Parse()

	if *help == true {
		HelpOutput()
		os.Exit(0)
	}

	os.Args = RemoveArgs()

	//var all_totalSize int64 = 0

	// check the args for a working path. we need to move this to some args processing once we start adding more features.
	workingPath := "."
	if len(os.Args) > 1 {
		workingPath = os.Args[len(os.Args)-1]
	}
	// if no path is specified we get the current working directory so that we can print the path rather than just a "."
	if workingPath == "." {
		workingPath = GetCwd()
	}
	// when on windows searching the path c: would raise an error. so we now add a "/" onto the path.
	if workingPath[len(workingPath)-1:] == ":" {
		workingPath = workingPath + "/"
	}
	// bit of stuff for consistancy.
	workingPath = strings.Replace(workingPath, "\\", "/", -1)

	// finally we list the path
	if DoesPathExist(workingPath) == true {
		ListPath(workingPath)
	}

	// fmt.Println(all_totalSize)

	//directories := GetDirectories(workingPath)
	//for _, directory := range directories {
	//	ListPath(directory)
	//}

	//if *recursive == true { // doesnt do anything yet
	//	fmt.Println(all_totalSize)
	//}
}
