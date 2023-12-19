package domain

import (
	"fmt"
	"os"
	"path/filepath"
)

type NodeBiz struct {
}

const (
	TypeFile = "file"
	TypeDir  = "dir"
)

type NodeDo struct {
	Name     string
	Type     string
	Children []NodeDo
	Path     string
}

func NewNodeBiz() *NodeBiz {
	return &NodeBiz{}
}

func (n *NodeBiz) GetNodes(dir string) (NodeDo, error) {
	node, err := createNode(dir)
	if err != nil {
		fmt.Printf("Error creating node: %v\n", err)
		return NodeDo{}, err
	}

	// 打印生成的节点树
	//printNode(node, 0)
	return node, nil
}

func createNode(path string) (NodeDo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return NodeDo{}, err
	}

	node := NodeDo{
		Name: info.Name(),
		Type: getType(info),
		Path: path,
	}

	if info.IsDir() {
		files, err := filepath.Glob(filepath.Join(path, "*"))
		if err != nil {
			return NodeDo{}, err
		}

		for _, file := range files {
			childNode, err := createNode(file)
			if err != nil {
				fmt.Printf("Error creating child node: %v\n", err)
				continue
			}
			node.Children = append(node.Children, childNode)
		}
	}

	return node, nil
}

func getType(info os.FileInfo) string {
	if info.IsDir() {
		return "Directory"
	}
	return "File"
}

func printNode(node NodeDo, indent int) {
	fmt.Printf("%sName: %s, Type: %s, Path: %s\n", generateIndent(indent), node.Name, node.Type, node.Path)

	for _, child := range node.Children {
		printNode(child, indent+1)
	}
}

func generateIndent(indent int) string {
	if indent <= 0 {
		return ""
	}
	return fmt.Sprintf("%"+fmt.Sprintf("%d", indent*4)+"s", "")
}
