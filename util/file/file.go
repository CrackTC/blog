package file

import (
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"time"
)

type Tree struct {
	Name     string
	Children []*Tree
}

func FindFile(root, name string, ignoredPaths []string) string {
	var res string
	err := filepath.WalkDir(root, func(path string, ent fs.DirEntry, err error) error {
		base := ent.Name()
		for _, ignoredPath := range ignoredPaths {
			if base == ignoredPath {
				if ent.IsDir() {
					return fs.SkipDir
				} else {
					return nil
				}
			}
		}
		if !ent.IsDir() && base == name {
			res, _ = filepath.Rel(root, path)
			return fs.SkipAll
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return res
}

type Time struct {
	Name    string
	Path    string
	ModTime string
}

func GetFileTimesRecursive(root string, ignoredPaths []string) []Time {
	var res []Time
	var cstZone = time.FixedZone("CST", 8*3600)
	err := filepath.WalkDir(root, func(path string, ent fs.DirEntry, err error) error {
		base := ent.Name()
		for _, ignoredPath := range ignoredPaths {
			if base == ignoredPath {
				if ent.IsDir() {
					return fs.SkipDir
				} else {
					return nil
				}
			}
		}
		if !ent.IsDir() {
			if info, err := ent.Info(); err != nil {
				return err
			} else {
				t := info.ModTime().In(cstZone).Format("2006-01-02 15:04:05")
				res = append(res, Time{base, path, t})
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ModTime > res[j].ModTime
	})
	return res
}
