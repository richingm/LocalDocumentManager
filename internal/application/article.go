package application

import (
	"fmt"
	"golang.org/x/net/context"
	"richingm/LocalDocumentManager/internal/domain"
	"richingm/LocalDocumentManager/internal/infrastructure/mysql"
	"richingm/LocalDocumentManager/internal/repo"
	"strings"
)

type ArticleService struct {
}

func NewArticleService(ctx context.Context) *ArticleService {
	return &ArticleService{}
}

type NodeDto struct {
	ID       string     `json:"id"`
	Topic    string     `json:"topic"`
	Children []*NodeDto `json:"children"`
	Expanded bool       `json:"expanded"` // 子节点默认不打开
}

type ArticleDto struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	OrderSort int    `json:"order_sort"`
}

func (s *ArticleService) Create(ctx context.Context, cid int, pid int, title string, content string) error {
	articleBiz := domain.NewArticleBiz(ctx, repo.NewArticleRepo(mysql.GormDb), repo.NewArticleContentRepo(mysql.GormDb))
	_, err := articleBiz.Create(ctx, cid, pid, title, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArticleService) Update(ctx context.Context, id int, title string, content string, orderSort int) error {
	articleBiz := domain.NewArticleBiz(ctx, repo.NewArticleRepo(mysql.GormDb), repo.NewArticleContentRepo(mysql.GormDb))
	err := articleBiz.Update(ctx, id, title, content, orderSort)
	if err != nil {
		return err
	}
	return nil
}

func (s *ArticleService) Delete(ctx context.Context, id int) error {
	articleBiz := domain.NewArticleBiz(ctx, repo.NewArticleRepo(mysql.GormDb), repo.NewArticleContentRepo(mysql.GormDb))
	err := articleBiz.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
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

	res.Children = buildTreeArticle(articleDos)
	return res, nil
}

func buildTreeArticle(articles []domain.ArticleDo) []*NodeDto {
	res := make([]*NodeDto, 0, len(articles))
	for _, article := range articles {
		res = append(res, buildTreeArticleLoop(article))
	}
	return res
}

func buildTreeArticleLoop(do domain.ArticleDo) *NodeDto {
	if len(do.Children) == 0 {
		return &NodeDto{
			ID:       fmt.Sprintf("%d", do.ID),
			Topic:    do.Title,
			Children: nil,
			Expanded: false,
		}
	}
	children := make([]*NodeDto, 0, len(do.Children))
	for _, child := range do.Children {
		children = append(children, buildTreeArticleLoop(child))
	}
	return &NodeDto{
		ID:       fmt.Sprintf("%d", do.ID),
		Topic:    do.Title,
		Children: children,
		Expanded: false,
	}
}

func (s *ArticleService) Get(ctx context.Context, id int) (*ArticleDto, error) {
	articleBiz := domain.NewArticleBiz(ctx, repo.NewArticleRepo(mysql.GormDb), repo.NewArticleContentRepo(mysql.GormDb))
	do, err := articleBiz.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &ArticleDto{
		ID:        do.ID,
		Title:     do.Title,
		Content:   do.Content,
		OrderSort: do.OrderSort,
	}, nil
}

func (s *ArticleService) Trees(ctx context.Context, id int) (string, error) {
	articleBiz := domain.NewArticleBiz(ctx, repo.NewArticleRepo(mysql.GormDb), repo.NewArticleContentRepo(mysql.GormDb))
	articleDo, err := articleBiz.Get(ctx, id)
	if err != nil {
		return "", err
	}
	articleDos, err := articleBiz.List(ctx, articleDo.CategoryID)
	res := make([]string, 0)
	res = append(res, "<select name=\"m_pid\">")
	buildTreeArticleStr(articleDos, articleDo.Pid, &res)
	res = append(res, "</select>")
	return strings.Join(res, ""), nil
}

func buildTreeArticleStr(articles []domain.ArticleDo, selectId int, res *[]string) {
	selected := ""
	if 0 == selectId {
		selected = "selected"
	}
	option := fmt.Sprintf("<option value=\"%d\" %s >%s</option>", 0, selected, "请选择")
	*res = append(*res, option)
	for _, article := range articles {
		buildTreeArticleLoopStr(article, res, selectId, "-")
	}
}

func buildTreeArticleLoopStr(do domain.ArticleDo, res *[]string, selectId int, suffix string) {
	selected := ""
	if do.ID == selectId {
		selected = "selected"
	}
	option := fmt.Sprintf("<option value=\"%d\" %s >%s</option>", do.ID, selected, suffix+do.Title)
	*res = append(*res, option)
	if len(do.Children) == 0 {
		return
	}
	suffix = suffix + suffix
	for _, child := range do.Children {
		buildTreeArticleLoopStr(child, res, selectId, suffix)
	}
}
