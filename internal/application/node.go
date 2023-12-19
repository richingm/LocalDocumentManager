package application

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
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
}

func (n *NodeService) GetContentByNodeId(dir string, nodeId string) {

}

func convertNodeDoToNodeDto(nodeDo domain.NodeDo) NodeDto {
	nodeDto := NodeDto{
		ID:    generateID(nodeDo.Path),
		Topic: nodeDo.Name,
	}

	for _, child := range nodeDo.Children {
		childDto := convertNodeDoToNodeDto(child)
		nodeDto.Children = append(nodeDto.Children, childDto)
	}
	return nodeDto
}

func (n *NodeService) GetMind(dir string, noteName string) (NodeDto, error) {
	nodeBiz := domain.NewNodeBiz()
	var res NodeDto
	nodeDo, err := nodeBiz.GetNodes(dir)
	if err != nil {
		return res, err
	}

	nodeDo.Name = noteName
	nodeDo.Path = "/"

	return convertNodeDoToNodeDto(nodeDo), nil
}

func (n *NodeService) GetContent(dir string, noteName string, nodeId string) (NodeDto, error) {
	nodeBiz := domain.NewNodeBiz()
	var res NodeDto
	nodeDo, err := nodeBiz.GetNodes(dir)
	if err != nil {
		return res, err
	}

	nodeDo.Name = noteName
	nodeDo.Path = "/"

	dto := convertNodeDoToNodeDto(nodeDo)

	foundNode := findNodeByID(dto, nodeId)
	if foundNode == nil {
		return res, errors.New("不存在")
	}

	fmt.Println(foundNode)

	return res, nil
}

func findNodeByID(node NodeDto, targetID string) *NodeDto {
	if node.ID == targetID {
		return &node
	}

	for _, child := range node.Children {
		foundNode := findNodeByID(child, targetID)
		if foundNode != nil {
			return foundNode
		}
	}

	return nil
}

func generateID(input string) string {
	// 计算SHA-256哈希值
	hash := sha256.Sum256([]byte(input))
	// 将哈希值转换为Base64编码
	return base64.StdEncoding.EncodeToString(hash[:])
}
