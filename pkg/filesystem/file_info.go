// Copyright 2023 The FreePDM team. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filesystem

// This shows only the latest version of a file and is handy when being used with
// the web server. If the file is a directory then only the filename is visible.

// The FileInfo struct
type FileInfo struct {
	Dir             bool             `json:"Dir"` // Is it a directory or a file?
	FileName        string           `json:"FileName"`
	FileDescription string           `json:"FileDescription"`
	FileSecondDescr string           `json:"FileSecondDescr"`
	FileVersion     string           `json:"FileVersion"`
	FileLocked      bool             `json:"FileLocked"`      // Is the file locked out?
	FileLockedOutBy string           `json:"FileLockedOutBy"` // Who checked out the file?
	FileIcon        []byte           `json:"FileIcon"`
	FileProperties  []FileProperties `json:"FileProperties"`
}

// The File Properties
type FileProperties struct {
	Key, Value string
}

func (fi FileInfo) IsDir() bool {
	return fi.Dir
}

// Returns the directory or file name
func (fi FileInfo) Name() string {
	return fi.FileName
}

func (fi FileInfo) Description() string {
	return fi.FileDescription
}

func (fi FileInfo) SecondDescription() string {
	return fi.FileSecondDescr
}

func (fi FileInfo) Version() string {
	return fi.FileVersion
}

func (fi FileInfo) IsLocked() bool {
	return fi.FileLocked
}

func (fi FileInfo) LockedOutBy() string {
	return fi.FileLockedOutBy
}

func (fi FileInfo) Properties() []FileProperties {
	return fi.FileProperties
}