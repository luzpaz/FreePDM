// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem_test

// Tests for the filesystem.

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"testing"

	"github.com/grd/FreePDM/pkg/config"
	fsm "github.com/grd/FreePDM/pkg/filesystem"
	"github.com/grd/FreePDM/pkg/util"
	"github.com/stretchr/testify/assert"
)

const testpdm = "testpdm"

var (
	fs *fsm.FileSystem

	testvaults     = path.Join(fsm.Vaults, testpdm)
	testvaultsdata = path.Join(fsm.VaultsData, testpdm)

	file1, file2, file3, file4, file5, file6 string

	freePdmDir = path.Join(os.Getenv("HOME"), "FreePDM")
)

func TestInitFileSystem(t *testing.T) {
	userName, err := user.Current()
	if err != nil {
		t.Errorf("user.Current() error: %s", err)
	}
	fs, err = fsm.NewFileSystem(testpdm, userName.Username)
	if err != nil {
		log.Fatalf("initialization failed, %v", err)
	}
}

func TestMkDir(t *testing.T) {
	if err := fs.Mkdir("Standard Parts"); err != nil {
		t.Errorf("Mkdir error message = %s", err)
	}

	fs.Mkdir("Projects")
	fs.Mkdir("test")

}

func TestImportFile(t *testing.T) {
}

func TestTheRest(t *testing.T) {

	userName, err := user.Current()
	if err != nil {
		t.Errorf("user.Current() error: %s", err)
	}

	fmt.Printf("Vault dir: %s\n", testvaults)
	fmt.Printf("Vault data dir: %s\n", testvaultsdata)
	fmt.Printf("User name: %s\n", userName.Name)

	filesDir := path.Join(freePdmDir, "ConceptOfDesign/TestFiles")
	file1 = path.Join(filesDir, "0001.FCStd")
	file2 = path.Join(filesDir, "0002.FCStd")
	file3 = path.Join(filesDir, "0003.FCStd")
	file4 = path.Join(filesDir, "0004.FCStd")
	file5 = path.Join(filesDir, "0005.FCStd")
	file6 = path.Join(filesDir, "0006.FCStd")

	listWd()

	if err = fs.Chdir("Projects"); err != nil {
		t.Errorf("Chdir error message = %s", err)
	}

	{
		f1, err := fs.ImportFile(file1)
		if err != nil {
			t.Errorf("ImportFile %s error: %s", file1, err)
		}

		{
			compareFileListLine(1, "1:0001.FCStd::Projects:")
			checkOutStatus(1, 0)
		}
		fs.CheckIn(f1, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf1-0", "Testf1-0")
		{
			checkInStatus(1, 0)
		}

		ver, _ := fs.NewVersion(f1)

		{
			checkOutStatus(1, 1)
		}
		fs.CheckIn(f1, ver, "Testf1-1", "Test1-1")
		{
			checkInStatus(1, 1)
		}

		ver, _ = fs.NewVersion(f1)
		{
			checkOutStatus(1, 2)
		}
		fs.CheckIn(f1, ver, "Testf1-2", "Test1-2")
		{
			checkInStatus(1, 2)
		}

		ver, _ = fs.NewVersion(f1)
		{
			checkOutStatus(1, 3)
		}
		fs.CheckIn(f1, ver, "Testf1-3", "Test1-3")
		{
			checkInStatus(1, 3)
		}
	}

	{
		f2, err := fs.ImportFile(file2)
		if err != nil {
			t.Errorf("ImportFile %s error: %s", file1, err)
		}
		{
			compareFileListLine(2, "2:0002.FCStd::Projects:")
			checkOutStatus(2, 0)
		}
		fs.CheckIn(f2, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf2-0", "")
		{
			checkInStatus(2, 0)
		}
	}

	{
		f3, err := fs.ImportFile(file3)
		if err != nil {
			t.Errorf("ImportFile %s error: %s", file1, err)
		}
		{
			compareFileListLine(3, "3:0003.FCStd::Projects:")
			checkOutStatus(3, 0)
		}
		fs.CheckIn(f3, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf3-0", "")
		{
			checkInStatus(3, 0)
		}

		ver, _ := fs.NewVersion(f3)
		{
			checkOutStatus(3, 1)
		}
		fs.CheckIn(f3, ver, "Testf3-1", "Test3-1")
		{
			checkInStatus(3, 1)
		}
	}

	//
	// rename, copy and move requires a bit of love...
	//
	// these functions require that all the files are checked in and will fail when they are not.
	//

	{
		if err = fs.FileRename("0001.FCStd", "0007.FCStd"); err != nil {
			t.Errorf("FileRename %s error: %s", file1, err)
		}
		{
			compareFileListLine(1, "1:0007.FCStd:0001.FCStd:Projects:")
		}
	}

	{
		if err = fs.FileCopy("0002.FCStd", "0008.FCStd"); err != nil {
			t.Errorf("FileCopy %s error: %s", file1, err)
		}
		{
			compareFileListLine(4, "4:0008.FCStd::Projects:")
		}
	}

	{
		if err = fs.FileCopy("0002.FCStd", "../test/0008a.FCStd"); err != nil {
			t.Errorf("FileCopy %s error: %s", file1, err)
		}
		{
			compareFileListLine(5, "5:0008a.FCStd::test:")
		}
	}

	{
		if err = fs.FileCopy("0002.FCStd", "../test/0010.FCStd"); err != nil {
			t.Errorf("FileCopy failed: %v", err)
		}
		{
			compareFileListLine(6, "6:0010.FCStd::test:")
		}
	}

	listWd()

	if err = fs.Mkdir("temp"); err != nil {
		t.Errorf("Mkdir failed: %v", err)
	}

	{
		if err = fs.FileMove("0002.FCStd", "temp"); err != nil {
			t.Errorf("FileMove failed: %v", err)
		}
		{
			compareFileListLine(2, "2:0002.FCStd::Projects/temp:Projects")
		}
	}

	{
		if err = fs.FileMove("0003.FCStd", ".."); err != nil {
			t.Errorf("FileMove failed: %v", err)
		}
		{
			compareFileListLine(3, "3:0003.FCStd:::Projects")
		}
	}

	listWd()

	if err = fs.Chdir("../Standard Parts"); err != nil {
		t.Errorf("Chdir failed: %v", err)
	}

	{
		f4, err := fs.ImportFile(file4)
		if err != nil {
			t.Errorf("ImportFile %s error: %s", file1, err)
		}
		fs.CheckIn(f4, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf4-0", "Testf4-0")
		{
			compareFileListLine(7, "7:0004.FCStd::Standard Parts:")
		}
	}

	{
		f5, err := fs.ImportFile(file5)
		if err != nil {
			t.Errorf("ImportFile %s error: %s", file1, err)
		}
		fs.CheckIn(f5, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf5-0", "Testf5-0")
		{
			compareFileListLine(8, "8:0005.FCStd::Standard Parts:")
		}
	}

	{
		f6, err := fs.ImportFile(file6)
		if err != nil {
			t.Errorf("ImportFile %s error: %s", file1, err)
		}
		fs.CheckIn(f6, fsm.FileVersion{Number: 0, Pretty: "0"}, "Testf6-0", "Testf6-0")
		{
			compareFileListLine(9, "9:0006.FCStd::Standard Parts:")
		}
	}

	listWd()

	if err = fs.Chdir("../Projects"); err != nil {
		t.Errorf("Chdir %s error: %s", file1, err)
	}

	listWd()

	fs.Chdir("..")

	listWd()

	{
		if err = fs.FileMove("0003.FCStd", "Projects"); err != nil {
			t.Errorf("FileMove failed: %v", err)
		}
		{
			compareFileListLine(3, "3:0003.FCStd::Projects:")
		}
	}
}

func TestListTree(t *testing.T) {
	lt, err := fs.ListTree("Projects")
	if err != nil {
		t.Errorf("listTree failed: %v", err)
	}
	list := make([]string, len(lt))
	for i, item := range lt {
		list[i] = path.Join(item.Path(), item.Name())
	}

	test := []string{
		"Projects/temp",
		"Projects/0003.FCStd",
		"Projects/0007.FCStd",
		"Projects/0008.FCStd",
		"Projects/0008a.FCStd",
		"Projects/0010.FCStd",
		"Projects/temp/0002.FCStd",
	}
	for i := range list {
		assert.Equal(t, list[i], test[i])

	}
}

// func TestFileRename(t *testing.T) {
// 	err := fs.FileRename("src.txt", "dest.txt")
// 	if err != nil {
// 		t.Errorf("FileRename failed: %v", err)
// 	}

// 	// Test cases for locking, existing destination, etc.
// 	err = fs.FileRename("lockedFile.txt", "newFile.txt")
// 	if err == nil || err.Error() != "FileRename error: File lockedFile.txt is checked out by user" {
// 		t.Errorf("Expected lock error, got: %v", err)
// 	}

// 	// Further test cases, e.g., file already exists.
// }

// func TestFileCopy(t *testing.T) {
// 	fs := fsm.FileSystem{
// 		// Initialize mock FileSystem or dependencies here
// 	}

// 	fs.Chdir("Projects")
// 	err := fs.FileCopy("0003.FCStd", "0011.FCStd")
// 	if err != nil {
// 		t.Errorf("FileCopy failed: %v", err)
// 	}

// 	// // Test locked file
// 	// err = fs.FileCopy("lockedFile.txt", "newFile.txt")
// 	// if err == nil || err.Error() != "FileCopy error: File lockedFile.txt is checked out by user" {
// 	// 	t.Errorf("Expected lock error, got: %v", err)
// 	// }

// 	// // Test for file already existing at destination
// 	// err = fs.FileCopy("src.txt", "existingFile.txt")
// 	// if err == nil || err.Error() != "file existingFile.txt already exists" {
// 	// 	t.Errorf("Expected existing file error, got: %v", err)
// 	// }
// }

func TestDirectoryRename(t *testing.T) {
	// Test with live data, including sub-dirs
	err := fs.DirectoryRename("Projects", "test/Projects3")
	if err != nil {
		t.Errorf("DirectoryMove failed: %v", err)
	}

	// Test for destination directory being a number
	err = fs.DirectoryRename("Standard Parts", "123")
	if err == nil || err.Error() != "destination directory 123 cannot be a number" {
		t.Errorf("Expected number error, got: %v", err)
	}

	// Test for source directory not existing
	err = fs.DirectoryRename("nonExistentDir", "destDir")
	if err == nil || err.Error() != "source directory nonExistentDir does not exist" {
		t.Errorf("Expected source directory error, got: %v", err)
	}

	// Test for creating new directory
	if err := fs.Mkdir("temp"); err != nil {
		t.Errorf("Expected to create dir temp got %v", err)
	}

	// Test for empty source directory
	err = fs.DirectoryRename("temp", "destDir")
	if err == nil || err.Error() != "directory is empty" {
		t.Errorf("Expected source directory to be empty, got: %v", err)
	}

	// Test with removal of the test directory
	err = os.Remove("temp")
	if err != nil {
		t.Errorf("Removal of directory temp failed: %v", err)
	}

	// Test with live data, including sub-dirs, return to Projects3
	err = fs.DirectoryRename("test/Projects3", "Projects3")
	if err != nil {
		t.Errorf("DirectoryMove failed: %v", err)
	}

	// Test with live data, including sub-dirs, rename Projects3 to Projects again
	err = fs.DirectoryRename("Projects3", "Projects")
	if err != nil {
		t.Errorf("DirectoryMove failed: %v", err)
	}

	// Test with checked-out file

	file, err := fs.GetItem("0002.FCStd")
	if err != nil {
		t.Errorf("GetItem %s error: %s", file.Name(), err)
	}
	if err = fs.CheckOut(file.IndexInt64(), fsm.FileVersion{Number: 0, Pretty: "0"}); err != nil {
		t.Errorf("Checkout %s error: %s", file.Name(), err)
	}
	checkOutStatus(2, 0)

	err = fs.DirectoryRename("Projects", "Projects4")
	if err == nil || err.Error() != "check out errors: [0002.FCStd is checked out by user]" {
		t.Errorf("Expected one file checked out, got: %v", err)
	}

	fs.CheckIn(file.IndexInt64(), fsm.FileVersion{Number: 0, Pretty: "0"}, "", "")
	checkInStatus(2, 0)
}

func TestMain(m *testing.M) {

	setup()

	m.Run()
}

// A simple "quick and dirty" remove all files from testvault
// but it also initializes the startup files again.

func setup() {

	// cleanup the testvaults

	err := os.RemoveAll(testvaults)
	util.CheckErr(err)

	err = os.Mkdir(testvaults, 0775)
	util.CheckErr(err)

	vaultUid := config.GetUid("vault")

	err = os.Chown(testvaults, os.Geteuid(), vaultUid)
	util.CheckErr(err)

	// cleanup the testvaultsdata

	err = os.RemoveAll(testvaultsdata)
	util.CheckErr(err)

	err = os.Mkdir(testvaultsdata, 0775)
	util.CheckErr(err)

	err = os.Chown(testvaultsdata, os.Geteuid(), vaultUid)
	util.CheckErr(err)

	// writing three files...

	writeIndexFileList()

	writeLockedFiles()

	writeIndexNumber()
}

func writeIndexFileList() {

	fileList := path.Join(testvaultsdata, "FileList.csv")

	f, err := os.Create(fileList)
	util.CheckErr(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Comma = ':'

	firstRecord := []string{
		"Index", "FileName", "PreviousFile", "Directory", "PreviousDir",
	}

	if err := w.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	vaultUid := config.GetUid("vault")

	err = os.Chown(fileList, os.Geteuid(), vaultUid)
	util.CheckErr(err)
}

func writeLockedFiles() {

	lockedfile := path.Join(testvaultsdata, "LockedFiles.csv")

	f, err := os.Create(lockedfile)
	util.CheckErr(err)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Comma = ':'

	firstRecord := []string{"FileNumber", "Version", "UserName"}

	if err := w.Write(firstRecord); err != nil {
		log.Fatalln("error writing record to file", err)
	}

	vaultUid := config.GetUid("vault")

	err = os.Chown(lockedfile, os.Geteuid(), vaultUid)
	util.CheckErr(err)
}

func writeIndexNumber() {

	vaultUid := config.GetUid("vault")

	idxfile := path.Join(testvaultsdata, "IndexNumber.txt")

	err := os.WriteFile(idxfile, []byte{'0'}, 0644)
	util.CheckErr(err)

	err = os.Chown(idxfile, os.Geteuid(), vaultUid)
	util.CheckErr(err)
}

func listWd() {
	wd, err := os.Getwd()
	util.CheckErr(err)
	fmt.Printf("\n\nList directory of %s\n\n", wd)
	fileInfo, err := fs.ListWD()
	if err != nil {
		log.Fatal("listwd error")
	}
	for _, info := range fileInfo {
		fmt.Println(info.FileName)
	}
	fmt.Println("")

}

func compareFileListLine(num int, line string) {
	content, err := os.ReadFile(path.Join(testvaultsdata, "FileList.csv"))
	if err != nil {
		log.Fatal(err)
	}
	s := string(content[:len(content)-1])
	slice := strings.Split(s, "\n")

	if slice[num] != line {
		log.Fatalf("error comparing lines of FileList.csv\nRead line (line %d): %s\nLine parameter:     %s", num, slice[num], line)
	}
}

func checkStatus(file int64, num int16) bool {
	content, err := os.ReadFile(path.Join(testvaultsdata, "LockedFiles.csv"))
	if err != nil {
		log.Fatal(err)
	}
	s := string(content[:len(content)-1])
	slice := strings.Split(s, "\n")
	slice = slice[1:] // first line
	if len(slice) == 0 {
		return false
	}
	for _, line := range slice {
		segments := strings.Split(line, ":")
		seg0, _ := util.Atoi64(segments[0])
		seg1, _ := util.Atoi16(segments[1])
		if seg0 == file && seg1 == num {
			return true
		}

	}

	return false
}

func checkOutStatus(file int64, num int16) {
	done := checkStatus(file, num)
	if !done {
		log.Fatalf("error file not properly checked out")
	}
}

func checkInStatus(file int64, num int16) {
	done := checkStatus(file, num)
	if done {
		log.Fatalf("error file not properly checked in")
	}
}

// func checkFileComments(file int64, dir string, num int16, descr, longDescr string){

// }
