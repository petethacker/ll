package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type file_list struct {
	name    string
	size    string
	modTime string
	isDir   bool
	symlink bool
}

func get_cwd() string {
	dir, _ := os.Getwd()
	return dir
}

func symlink_check(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		log.Fatal(err)
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

func main() {
	var largest_size int = 4
	var total_size int64 = 0
	var size string
	var total_files int = 0
	var total_dirs int = 0
	var total_links int = 0
	var bonus_spacing int = 2
	var working_path_target string = ""

	// check the args for a working path. we need to move this to some args processing once we start adding more features.
	working_path := "."
	if len(os.Args) > 1 {
		working_path = os.Args[1]
	}
	// if no path is specified we get the current working directory so that we can print the path rather than just a "."
	if working_path == "." {
		working_path = get_cwd()
	}
	// when on windows searching the path c: would raise an error. so we now add a "/" onto the path.
	if working_path[len(working_path)-1:] == ":" {
		working_path = working_path + "/"
	}
	// bit of stuff for consistancy.
	working_path = strings.Replace(working_path, "\\", "/", -1)

	if symlink_check(working_path) == true {
		link_path, err := os.Readlink(working_path)
		if err != nil {
			fmt.Println(working_path + " is a symlink but we failed to get the source path.")
			os.Exit(1)
		}
		working_path_target = link_path
	}

	files, err := ioutil.ReadDir(working_path)
	if err != nil {
		fmt.Println(working_path + " is not a valid directory.")
		os.Exit(1)
	}

	storage := map[file_list]bool{}
	for _, f := range files {
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
			if symlink_check(path.Join(working_path, f.Name())) == true {
				s.symlink = true
				total_links += 1
				bonus_spacing = 7
			}
		} else {
			total_files += 1
		}

		storage[*s] = true
	}

	fmt.Println("")
	for i := range storage {
		if i.isDir == true {
			size = "<DIR>" + strings.Repeat(" ", largest_size+bonus_spacing)
		} else if i.symlink == true {
			size = "<JUNCTION>" + strings.Repeat(" ", largest_size+bonus_spacing-5)
			link, err := os.Readlink(path.Join(working_path, i.name))
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
	path_string := strings.Repeat(" ", 14) + "Path\t" + working_path
	if len(working_path_target) > 0 {
		path_string = path_string + " => " + working_path_target
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
}
