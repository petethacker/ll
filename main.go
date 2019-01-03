package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type file_list struct {
	name    string
	size    string
	modTime string
	isDir   bool
}

func get_cwd() string {
	dir, _ := os.Getwd()
	return dir
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
	var largest_size int = 7
	var total_size int64 = 0
	var size string
	var total_files int = 0
	var total_dirs int = 0

	working_path := "."
	if len(os.Args) > 1 {
		working_path = os.Args[1]
	}

	if working_path == "." {
		working_path = get_cwd()
	}

	files, err := ioutil.ReadDir(working_path)
	if err != nil {
		log.Fatal(err)
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
		} else {
			total_files += 1
		}
		storage[*s] = true
	}

	fmt.Println("")
	for i := range storage {
		if i.isDir == true {
			size = "<DIR>" + strings.Repeat(" ", largest_size)
		} else {
			size = strings.Repeat(" ", largest_size-len(i.size)+len("<DIR>")) + i.size
		}
		fmt.Println(i.modTime + "  " + size + "  " + i.name)
	}
	fmt.Println("")
	fmt.Println(strings.Repeat(" ", 14) + "Path - " + working_path)
	fmt.Println(strings.Repeat(" ", 14) + strconv.Itoa(total_dirs) + " Dir(s)")
	fmt.Println(strings.Repeat(" ", 14) + strconv.Itoa(total_files) + " File(s)\t\t" + size_commaed(total_size) + " bytes")
}
