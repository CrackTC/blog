package file

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"sora.zip/blog/config"
)

type FileTree struct {
	Name     string
	Children []*FileTree
}

func GetFileTree(dir string) (*FileTree, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Println("[ERROR] Failed to read directory:", dir)
		return nil, errors.New("Failed to read directory")
	}

	root := &FileTree{Name: dir}
	root.Children = make([]*FileTree, 0, len(entries))

	for _, entity := range entries {
		var shouldIgnore bool
		for _, ignore := range config.Get().IgnoredPaths {
			if entity.Name() == ignore {
				shouldIgnore = true
				break
			}
		}
		if shouldIgnore {
			continue
		}

		if entity.IsDir() {
			subtree, err := GetFileTree(filepath.Join(dir, entity.Name()))
			if err != nil {
				return nil, err
			}
			root.Children = append(root.Children, subtree)
		} else {
			root.Children = append(root.Children, &FileTree{Name: filepath.Join(dir, entity.Name())})
		}
	}

	return root, nil
}
