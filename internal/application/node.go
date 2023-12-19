package application

import (
	"crypto/sha256"
	"encoding/base64"
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
	// 计算SHA-256哈希值
	hash := sha256.Sum256([]byte(input))
	// 将哈希值转换为Base64编码
	return base64.StdEncoding.EncodeToString(hash[:])
}
