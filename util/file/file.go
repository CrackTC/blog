package file

import (
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"time"
)

type FileTree struct {
	Name     string
	Children []*FileTree
}

func FindFile(root, name string, ignoredPaths []string) string {
	var res string
	filepath.WalkDir(root, func(path string, ent fs.DirEntry, err error) error {
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
	return res
}

type FileTime struct {
	Name    string
	Path    string
	ModTime string
}

func GetFileTimesRecursive(root string, ignoredPaths []string) []FileTime {
	var res []FileTime
	log.Println(root, ignoredPaths)
	var cstZone = time.FixedZone("CST", 8*3600)
	filepath.WalkDir(root, func(path string, ent fs.DirEntry, err error) error {
		log.Println(path, ent, err)
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
				time := info.ModTime().In(cstZone).Format("2006-01-02 15:04:05")
				res = append(res, FileTime{base, path, time})
			}
		}
		return nil
	})

	sort.Slice(res, func(i, j int) bool {
		return res[i].ModTime > res[j].ModTime
	})
	return res
}
