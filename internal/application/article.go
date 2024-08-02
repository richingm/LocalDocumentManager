package application

import (
	"fmt"
	"golang.org/x/net/context"
	"richingm/LocalDocumentManager/internal/domain"
	"richingm/LocalDocumentManager/internal/infrastructure/mysql"
	"richingm/LocalDocumentManager/internal/repo"
)

type ArticleService struct {
}

func NewArticleService(ctx context.Context) *ArticleService {
	return &ArticleService{}
}

type NodeDto struct {
	ID       string    `json:"id"`
	Topic    string    `json:"topic"`
	Children []NodeDto `json:"children"`
	Expanded bool      `json:"expanded"` // 子节点默认不打开
}

type ArticleDto struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *ArticleService) Nodes(ctx context.Context, categoryId int) (NodeDto, error) {
	categoryBiz := domain.NewCategoryBiz(ctx, repo.NewCategoryRepo(mysql.GormDb))
	categoryDo, err := categoryBiz.Get(ctx, categoryId)
	if err != nil {
		return NodeDto{}, err
	}
	res := NodeDto{
		ID:       "0",
		Topic:    categoryDo.Name,
		Expanded: false,
	}
	articleBiz := domain.NewArticleBiz(ctx, repo.NewArticleRepo(mysql.GormDb), repo.NewArticleContentRepo(mysql.GormDb))
	articleDos, err := articleBiz.List(ctx, categoryId)
	if err != nil {
		return NodeDto{}, err
	}

	res.Children = convertArticleToNodeDto(articleDos)
	return res, nil
}

func convertArticleToNodeDto(articles []*domain.ArticleDo) []NodeDto {
	nodeMap := make(map[int]*NodeDto)

	// Create nodes
	for _, article := range articles {
		node, exists := nodeMap[article.ID]
		if !exists {
			node = &NodeDto{
				ID:       fmt.Sprintf("%d", article.ID),
				Topic:    article.Title,
				Children: make([]NodeDto, 0),
				Expanded: false,
			}
			nodeMap[article.ID] = node
		}

		if article.Pid != 0 {
			parentNode, exists := nodeMap[article.Pid]
			if !exists {
				parentNode = &NodeDto{
					ID:       fmt.Sprintf("%d", article.ID),
					Children: make([]NodeDto, 0),
					Expanded: false,
				}
				nodeMap[article.Pid] = parentNode
			}
			parentNode.Children = append(parentNode.Children, *node)
		}
	}

	// Extract nodes to slice
	nodes := make([]NodeDto, 0, len(nodeMap))
	for _, node := range nodeMap {
		nodes = append(nodes, *node)
	}

	return nodes
}

func (s *ArticleService) Get(ctx context.Context, id int) (*ArticleDto, error) {
	articleBiz := domain.NewArticleBiz(ctx, repo.NewArticleRepo(mysql.GormDb), repo.NewArticleContentRepo(mysql.GormDb))
	do, err := articleBiz.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &ArticleDto{
		ID:      do.ID,
		Title:   do.Title,
		Content: do.Content,
	}, nil
}
