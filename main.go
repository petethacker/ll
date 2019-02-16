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

type file_list struct {
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
		//log.Fatal(err)	//log.Fatal(err) // if we have this enabled it raises access errors on some files unless ll is run as admin.
		return false
	}
	switch mode := fi.Mode(); {
	case mode&os.ModeSymlink != 0:
		return true
	}
	return false
}

func size_commaed_old(old_size int64) string {
	size := strconv.Itoa(int(old_size))
	if len(size) < 4 {
		return size
	}
	var commaed string
	index := 1
	for i := len(size) - 1; i >= 0; i-- {
		commaed = string(size[i]) + commaed
		if index > 2 {
			commaed = "," + commaed
			index = 1
		} else {
			index += 1
		}
	}
	return commaed
}

func size_commaed(old_size int64) string {
	size := strconv.Itoa(int(old_size))
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
		index += 1
	}
	return commaed
}

func help_output() {
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

func list_path(workingPath string) int64 {
	var largest_size int = 4
	var total_size int64 = 0
	var size string
	var total_files int = 0
	var total_dirs int = 0
	var total_links int = 0
	var bonus_spacing int = 2
	var workingPath_target string = ""

	if SymlinkCheck(workingPath) == true {
		linkPath, err := os.Readlink(workingPath)
		if err != nil {
			fmt.Println(workingPath + " is a symlink but we failed to get the source path.")
			os.Exit(1)
		}
		workingPath_target = linkPath
	}

	files, err := ioutil.ReadDir(workingPath)
	if err != nil {
		fmt.Println(workingPath + " is not a valid directory.")
		os.Exit(1)
	}

	storage := map[file_list]bool{}
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
			s := new(file_list)
			s.name = f.Name()
			s.size = size_commaed(f.Size())
			total_size += f.Size()
			if len(s.size) > largest_size {
				largest_size = len(s.size)
			}
			s.modTime = f.ModTime().String()[:19]
			s.isDir = f.IsDir()
			if s.isDir == true {
				total_dirs += 1
			} else if f.Size() == 0 {
				if SymlinkCheck(path.Join(workingPath, f.Name())) == true {
					s.symlink = true
					total_links += 1
					bonus_spacing = 7
				}
			} else {
				total_files += 1
			}
			storage[*s] = true
		}
	}

	fmt.Println("")
	for i := range storage {
		if i.isDir == true {
			size = "<DIR>" + strings.Repeat(" ", largest_size+bonus_spacing)
		} else if i.symlink == true {
			size = "<JUNCTION>" + strings.Repeat(" ", largest_size+bonus_spacing-5)
			link, err := os.Readlink(path.Join(workingPath, i.name))
			if err != nil {
				continue
			}
			i.name = i.name + " => " + link
		} else {
			size = strings.Repeat(" ", largest_size-len(i.size)+len("<DIR>")+bonus_spacing) + i.size
		}
		fmt.Println(i.modTime + "  " + size + "  " + i.name)
	}
	fmt.Println("")
	path_string := strings.Repeat(" ", 14) + "Path\t" + workingPath
	if len(workingPath_target) > 0 {
		path_string = path_string + " => " + workingPath_target
	}
	fmt.Println(path_string)
	if total_dirs > 0 {
		fmt.Println(strings.Repeat(" ", 14) + strconv.Itoa(total_dirs) + " Dir(s)")
	}
	if total_links > 0 {
		fmt.Println(strings.Repeat(" ", 14) + strconv.Itoa(total_links) + " Symlink(s)")
	}
	if total_files > 0 {
		fmt.Println(strings.Repeat(" ", 14) + strconv.Itoa(total_files) + " File(s)\t\t" + size_commaed(total_size) + " bytes")
	}
	return total_size
}

func get_directories(workingPath string) []string {
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

func does_path_exist(workingPath string) bool {
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
		help_output()
		os.Exit(0)
	}

	os.Args = RemoveArgs()

	//var all_total_size int64 = 0

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
	if does_path_exist(workingPath) == true {
		list_path(workingPath)
	}

	// fmt.Println(all_total_size)

	//directories := get_directories(workingPath)
	//for _, directory := range directories {
	//	list_path(directory)
	//}

	//if *recursive == true { // doesnt do anything yet
	//	fmt.Println(all_total_size)
	//}
}
