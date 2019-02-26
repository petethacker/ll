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

	ct "github.com/daviddengcn/go-colortext"
)

// arguments
var help = flag.Bool("h", false, "Help")
var excludeDirectories = flag.Bool("xd", false, "Exclude Directories")
var excludeFiles = flag.Bool("xf", false, "Exclude Files")
var excludeSymlinks = flag.Bool("xs", false, "Exclude Symlinks")
var textSearch = flag.String("f", "", "Text Search")
var sizeCheck = flag.String("fs", "", "Highlight files larger than x")
var sizeCheckListOnly = flag.String("fso", "", "Only show files larger than x")
var sizeCheckOnly bool = false
var sizeCheckBytes int = 1125899906842620 * 1024

func print(value ...interface{}) {
	formatted_line := fmt.Sprintf(value[0].(string), value[1:len(value)]...)
	fmt.Println(formatted_line)
}

func printNormal(value ...interface{}) {
	ct.ChangeColor(ct.White, true, ct.Black, false)
	formatted_line := fmt.Sprintf(value[0].(string), value[1:len(value)]...)
	fmt.Println(formatted_line)
	ct.ResetColor()
}

func printSuccess(value ...interface{}) {
	ct.ChangeColor(ct.Green, true, ct.Black, false)
	formatted_line := fmt.Sprintf(value[0].(string), value[1:len(value)]...)
	fmt.Println(formatted_line)
	ct.ResetColor()
}

func printWarning(value ...interface{}) {
	ct.ChangeColor(ct.Yellow, true, ct.Black, false)
	formatted_line := fmt.Sprintf(value[0].(string), value[1:len(value)]...)
	fmt.Println(formatted_line)
	ct.ResetColor()
}

func printError(value ...interface{}) {
	ct.ChangeColor(ct.Red, true, ct.Black, false)
	formatted_line := fmt.Sprintf(value[0].(string), value[1:len(value)]...)
	fmt.Println(formatted_line)
	ct.ResetColor()
}

type FileList struct {
	name      string
	size      string
	modTime   string
	isDir     bool
	symlink   bool
	aboveSize bool
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

func pause() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func StringCheck(path string, searchString string) bool {
	if strings.HasSuffix(searchString, "*") && strings.HasPrefix(searchString, "*") {
		searchString = searchString[+1:]
		searchString = searchString[:len(searchString)-1]
		if searchString == "." {
			return true
		}
	} else if strings.HasPrefix(searchString, "*") {
		searchString = searchString[+1:]
		if strings.HasSuffix(path, searchString) {
			return true
		}
	} else if strings.HasSuffix(searchString, "*") {
		searchString = searchString[:len(searchString)-1]
		if strings.HasPrefix(path, searchString) {
			return true
		}
	} else if strings.Contains(path, searchString) {
		return true
	}
	return false
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
			printWarning(workingPath + " is a symlink but we failed to get the source path.")
			os.Exit(1)
		}
		workingPathTarget = linkPath
	}

	files, err := ioutil.ReadDir(workingPath)
	if err != nil {
		printWarning(workingPath + " is not a valid directory.")
		os.Exit(1)
	}

	storage := map[FileList]bool{}
	for _, f := range files {
		if StringCheck(strings.ToLower(f.Name()), strings.ToLower(*textSearch)) {
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
			// check if the file is over the specified size for highlighting
			if int(f.Size()) > sizeCheckBytes {
				s.aboveSize = true
			}
			// check if we only want to list files above x size
			if sizeCheckOnly == true && s.aboveSize == false {
				continue
			}
			// if we are only listing files above x size then we dont need to highlight
			if sizeCheckOnly == true && s.aboveSize == true {
				s.aboveSize = false
			}
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
	if len(storage) > 0 {
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
			if i.aboveSize == true {
				printWarning(i.modTime + "  " + size + "  " + i.name)
			} else {
				fmt.Println(i.modTime + "  " + size + "  " + i.name)
			}
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
	} else {
		fmt.Println("")
		printWarning("No files or folders found")
		fmt.Println("      Path:\t" + workingPath)
		if *textSearch != "" {
			fmt.Println("    Search:\t" + *textSearch)
		}
	}
	return totalSize
}

func DoesPathExist(workingPath string) bool {
	_, err := os.Stat(workingPath)
	if err == nil {
		return true
	} else {
		printWarning("The path does not exist...")
		os.Exit(0)
	}
	return false
}

func processSizeCheck() int {
	var sizeCheckString string = *sizeCheck
	var newSizeCheck int = sizeCheckBytes
	if strings.HasSuffix(strings.ToLower(sizeCheckString), "kb") {
		// 1024
		if strings.ToLower(sizeCheckString) == "kb" {
			newSizeCheck = 1024
		} else {
			tempSize, _ := strconv.Atoi(sizeCheckString[:len(sizeCheckString)-2])
			newSizeCheck = tempSize * 1024
		}
	} else if strings.HasSuffix(strings.ToLower(sizeCheckString), "mb") {
		// 1048576
		if strings.ToLower(sizeCheckString) == "mb" {
			newSizeCheck = 1048576
		} else {
			tempSize, _ := strconv.Atoi(sizeCheckString[:len(sizeCheckString)-2])
			newSizeCheck = tempSize * 1048576
		}
	} else if strings.HasSuffix(strings.ToLower(sizeCheckString), "gb") {
		// 1073741824
		if strings.ToLower(sizeCheckString) == "gb" {
			newSizeCheck = 1073741824
		} else {
			tempSize, _ := strconv.Atoi(sizeCheckString[:len(sizeCheckString)-2])
			newSizeCheck = tempSize * 1073741824
		}
	} else if strings.HasSuffix(strings.ToLower(sizeCheckString), "tb") {
		// 1099511627776
		if strings.ToLower(sizeCheckString) == "tb" {
			newSizeCheck = 1099511627776
		} else {
			tempSize, _ := strconv.Atoi(sizeCheckString[:len(sizeCheckString)-2])
			newSizeCheck = tempSize * 1099511627776
		}
	} else if strings.HasSuffix(strings.ToLower(sizeCheckString), "pb") {
		// 1125899906842620
		if strings.ToLower(sizeCheckString) == "pb" {
			newSizeCheck = 1125899906842620
		} else {
			tempSize, _ := strconv.Atoi(sizeCheckString[:len(sizeCheckString)-2])
			newSizeCheck = tempSize * 1125899906842620
		}
	} else {
		tempSize, _ := strconv.Atoi(sizeCheckString)
		newSizeCheck = tempSize * 1
	}
	return newSizeCheck
}

func HelpOutput() {
	fmt.Println("Lists all files in a directory")
	fmt.Println()
	fmt.Println("ll <commands> <path to list files / directories>")
	fmt.Println()
	printWarning("  <Comamnds>")
	fmt.Println("  -xf : Exclude files")
	fmt.Println("  -xd : Exclude Directories")
	fmt.Println("  -xs : Exclude Symlinks")
	fmt.Println()
	fmt.Println("  -f <string>   : Search string to find files")
	fmt.Println("                 bla   = searches for strings that contain 'bla'")
	fmt.Println("                 *bla* = same as above")
	fmt.Println("                 bla*  = searches for strings starting with 'bla'")
	fmt.Println("                 *bla  = searches for strings ending with 'bla'")
	fmt.Println("  -fs <size>kb  : Highlight files above -fs <size>kb")
	fmt.Println("                 kb, mb, gb, tb, pb accepted.")
	fmt.Println("  -fso <size>kb : Only list files above -fso <size>")
	fmt.Println("                 kb, mb, gb, tb, pb accepted.")
	fmt.Println("                 will not highlight, exclude directories and exclude symlinks")
	fmt.Println()
	fmt.Println("  -h  : Help menu")
}

func RemoveArgs() []string {
	var newArgs []string
	last := "empty" // need a better solution for this, its total crap.
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-") == false {
			if strings.HasPrefix(last, "-f") == false && strings.HasPrefix(last, "-s") == false {
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

	if *sizeCheckListOnly != "" {
		*sizeCheck = *sizeCheckListOnly
		sizeCheckOnly = true
	}

	if *sizeCheck != "" {
		sizeCheckBytes = processSizeCheck()
	} else {
		sizeCheckOnly = false
	}

	if sizeCheckOnly == true {
		*excludeDirectories = true
		*excludeSymlinks = true
	}

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
}
