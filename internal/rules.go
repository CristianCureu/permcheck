package internal

import (
	"os"
	"path/filepath"
)

type PermissionRule struct {
	Name  string
	Match func(info os.FileInfo) (bool, os.FileMode)
}

var Rules = []PermissionRule{
	{
		Name: "Sensitive file by name",
		Match: func(info os.FileInfo) (bool, os.FileMode) {
			if perm, ok := SensitiveFiles[info.Name()]; ok {
				return true, perm
			}
			return false, 0
		},
	},
	{
		Name: "Sensitive file by extension",
		Match: func(info os.FileInfo) (bool, os.FileMode) {
			if perm, ok := SensitiveExtensions[filepath.Ext(info.Name())]; ok {
				return true, perm
			}
			return false, 0
		},
	},
}
