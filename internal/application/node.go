package application

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"richingm/LocalDocumentManager/internal/domain"
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

func (n *NodeService) GetContent(dir string, noteName string, nodeId string, fileSuffix string) (string, error) {
	fileBiz := domain.NewFileBiz()
	fieldDo, err := fileBiz.GetFiles(dir, fileSuffix)
	if err != nil {
		return "", err
	}

	isFile, filePath := getPathByNodeId(fieldDo, nodeId)
	if !isFile {
		return "", errors.New("不是文件")
	}
	if filePath == "" {
		return "", errors.New("数据不存在")
	}

	content, err := fileBiz.GetFileContent(filePath)
	if err != nil {
		return "", err
	}

	content = "<html><head></head><body>" + content + "</body></html>"

	return content, nil
}

func getPathByNodeId(file domain.FileDo, targetID string) (bool, string) {
	if generateID(file.Path) == targetID {
		return file.Type == domain.TypeFile, file.Path
	}

	for _, child := range file.Children {
		isFile, filePath := getPathByNodeId(child, targetID)
		if len(filePath) > 0 {
			return isFile, filePath
		}
	}
	return false, ""
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
