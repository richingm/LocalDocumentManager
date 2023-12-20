package application

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"os"
	"path/filepath"
	"richingm/LocalDocumentManager/internal/domain"
	"strings"
)

const (
	NoteImageRelativePath = "/static/note/images"
)

type NodeService struct {
}

func NewNodeService() *NodeService {
	return &NodeService{}
}

type NodeDto struct {
	ID       string    `json:"id"`
	Topic    string    `json:"topic"`
	Children []NodeDto `json:"children"`
	Expanded bool      `json:"expanded"` // 子节点默认不打开
}

func (n *NodeService) GetContentByNodeId(dir string, nodeId string) {

}

func convertNodeDoToNodeDto(fieldsDo domain.FileDo, level int) NodeDto {
	expanded := true
	if level <= 0 {
		expanded = false
	}
	nodeDto := NodeDto{
		ID:       generateID(fieldsDo.Path),
		Topic:    fieldsDo.Name,
		Expanded: expanded,
	}

	level = level - 1

	for _, child := range fieldsDo.Children {
		childDto := convertNodeDoToNodeDto(child, level)
		nodeDto.Children = append(nodeDto.Children, childDto)
	}
	return nodeDto
}

func (n *NodeService) GetMind(dir string, noteName string, level int, fileSuffix string) (NodeDto, error) {
	nodeBiz := domain.NewFileBiz()
	var res NodeDto
	fieldDo, err := nodeBiz.GetFiles(dir, fileSuffix)
	if err != nil {
		return res, err
	}
	fieldDo.Name = noteName
	return convertNodeDoToNodeDto(fieldDo, level), nil
}

func (n *NodeService) GetContentAndTitle(dir string, noteName string, nodeId string, fileSuffix string) (string, string, error) {
	fileBiz := domain.NewFileBiz()
	fieldDo, err := fileBiz.GetFiles(dir, fileSuffix)
	if err != nil {
		return "", "", err
	}

	isFile, filePath, title := getPathAndTitleByNodeId(fieldDo, nodeId)
	if !isFile {
		return "", "", errors.New("不是文件")
	}
	if filePath == "" {
		return "", "", errors.New("数据不存在")
	}

	content, err := fileBiz.GetFileContent(filePath)
	if err != nil {
		return "", "", err
	}

	content = strings.Replace(content, "\n\n", "<br>", -1)

	content = addBlankTargetAttribute(content)

	return fmt.Sprintf("%s/%s", noteName, title), content, nil
}

func getPathAndTitleByNodeId(file domain.FileDo, targetID string) (bool, string, string) {
	if generateID(file.Path) == targetID {
		return file.Type == domain.TypeFile, file.Path, file.Name
	}

	for _, child := range file.Children {
		isFile, filePath, title := getPathAndTitleByNodeId(child, targetID)
		if len(filePath) > 0 {
			return isFile, filePath, title
		}
	}
	return false, "", ""
}

func addBlankTargetAttribute(htmlString string) string {
	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		// 处理解析错误
		return htmlString
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// 检查是否已经存在target属性
			hasTargetAttr := false
			for _, attr := range n.Attr {
				if attr.Key == "target" {
					hasTargetAttr = true
					break
				}
			}

			// 如果不存在target属性，则添加target="_blank"
			if !hasTargetAttr {
				n.Attr = append(n.Attr, html.Attribute{
					Key: "target",
					Val: "_blank",
				})
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	var sb strings.Builder
	html.Render(&sb, doc)
	return sb.String()
}

func generateID(input string) string {
	data := []byte(input) // 替换为要生成MD5哈希值的数据
	return generateMD5(data)
}

func generateMD5(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func getNoteImagePath(dir string) string {
	parentDir := filepath.Dir(dir)
	grandparentDir := filepath.Dir(parentDir)
	return grandparentDir + "/" + strings.TrimLeft(NoteImageRelativePath, "/")
}

func (n *NodeService) ExtractImagePaths(fileDir string, htmlString string, destinationDir string) string {
	destinationDir = getNoteImagePath(destinationDir)

	_ = deleteFilesAndDirs(destinationDir)

	doc, err := html.Parse(strings.NewReader(htmlString))
	if err != nil {
		return htmlString
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for i, attr := range n.Attr {
				if attr.Key == "src" {
					// 移动图像文件到特定目录
					currentFile := attr.Val
					if attr.Val == filepath.Base(attr.Val) {
						currentFile = strings.TrimRight(fileDir, "/") + "/" + attr.Val
					}
					newPath := filepath.Join(destinationDir+"/", filepath.Base(currentFile))
					if err := moveFile(currentFile, newPath); err == nil {
						n.Attr[i].Val = NoteImageRelativePath + "/" + filepath.Base(currentFile)
					} else {
						panic(err)
					}
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	var sb strings.Builder
	if err := html.Render(&sb, doc); err != nil {
		return htmlString
	}

	return sb.String()
}

func deleteFilesAndDirs(dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == dirPath {
			return nil // 跳过传递的目录本身
		}

		if info.IsDir() {
			// 删除子目录
			err := os.RemoveAll(path)
			if err != nil {
				return err
			}
			fmt.Println("Deleted directory:", path)
		} else {
			// 删除文件
			err := os.Remove(path)
			if err != nil {
				return err
			}
			fmt.Println("Deleted file:", path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func moveFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = srcFile.Close()
	if err != nil {
		return err
	}

	return nil
}
