package internal

import "os"

type FileTask struct {
	Path string
	Info os.FileInfo
}
