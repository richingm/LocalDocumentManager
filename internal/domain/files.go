package domain

import (
	"github.com/russross/blackfriday/v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileBiz struct {
}

const (
	TypeFile = "file"
	TypeDir  = "dir"
)

type FileDo struct {
	Name     string
	Type     string
	Children []FileDo
	Path     string
}

func NewFileBiz() *FileBiz {
	return &FileBiz{}
}

func (n *FileBiz) GetFiles(dir string, fileSuffix string) (FileDo, error) {
	Files, err := loopFiles(dir, fileSuffix)
	if err != nil {
		return FileDo{}, err
	}
	return Files, nil
}

func (n *FileBiz) GetFileContent(path string) (string, error) {
	bytes, err := readFile(path)
	if err != nil {
		return "", err
	}
	html := blackfriday.Run(bytes)
	return string(html), nil
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func loopFiles(path string, fileSuffix string) (FileDo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileDo{}, err
	}

	Files := FileDo{
		Name: info.Name(),
		Type: getType(info),
		Path: path,
	}

	if info.IsDir() {
		files, err := filepath.Glob(filepath.Join(path, "*"))
		if err != nil {
			return FileDo{}, err
		}
		for _, file := range files {
			if filepath.Base(file) == "a.assets" {
				continue
			}
			childFiles, err := loopFiles(file, fileSuffix)
			if err != nil {
				return FileDo{}, err
			}
			if strings.HasSuffix(file, fileSuffix) || childFiles.Type == TypeDir {
				Files.Children = append(Files.Children, childFiles)
			}
		}
	}

	return Files, nil
}

func getType(info os.FileInfo) string {
	if info.IsDir() {
		return TypeDir
	}
	return TypeFile
}
