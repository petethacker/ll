package main

func HelpOutput() {
	printWarning("Lists all files in a directory")
	print("")
	printWarning("Usage : ll <commands> <path to list files / directories>")
	print("")
	printWarning("  <Comamnds>")
	print("  -xf : Exclude files")
	print("  -xd : Exclude Directories")
	print("  -xs : Exclude Symlinks")
	print("")
	print("  -f <string>   : Search string to find files")
	print("                 bla   = searches for strings that contain 'bla'")
	print("                 *bla* = same as above")
	print("                 bla*  = searches for strings starting with 'bla'")
	print("                 *bla  = searches for strings ending with 'bla'")
	print("  -fh <size>kb  : Highlight files above -fs <size>kb")
	print("                 kb, mb, gb, tb, pb accepted.")
	print("  -fso <size>kb : Show only files above -fso <size>")
	print("                 kb, mb, gb, tb, pb accepted.")
	print("                 will not highlight, exclude directories and exclude symlinks")
	print("")
	print("  -h  : Help menu")
}