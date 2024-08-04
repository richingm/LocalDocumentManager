package domain

import (
	"golang.org/x/net/context"
	"richingm/LocalDocumentManager/internal/entity"
	"richingm/LocalDocumentManager/internal/repo"
)

type ArticleDo struct {
	ID         int
	CreatedAt  []uint8
	Pid        int
	CategoryID int
	Title      string
	Children   []ArticleDo
}

type ArticleWithContentDo struct {
	ID         int
	CreatedAt  []uint8
	Pid        int
	CategoryID int
	Title      string
	Content    string
}

type ArticleBiz struct {
	articleRepo        *repo.ArticleRepo
	articleContentRepo *repo.ArticleContentRepo
}

func NewArticleBiz(ctx context.Context, articleRepo *repo.ArticleRepo, articleContentRepo *repo.ArticleContentRepo) *ArticleBiz {
	return &ArticleBiz{
		articleRepo:        articleRepo,
		articleContentRepo: articleContentRepo,
	}
}

func (b *ArticleBiz) List(ctx context.Context, categoryId int) ([]ArticleDo, error) {
	list, err := b.articleRepo.List(ctx, entity.ArticleParam{
		CategoryID: categoryId,
	})
	if err != nil {
		return nil, err
	}
	dos := convertArticlePosPtrToArticleDosPtr(list)
	dos = buildTree(dos, 0)
	return dos, nil
}

func (b *ArticleBiz) Create(ctx context.Context, categoryId, pid int, title string, content string) (*ArticleWithContentDo, error) {
	articlePo, err := b.articleRepo.Create(ctx, &entity.ArticlePo{
		Pid:        pid,
		CategoryID: categoryId,
		Title:      title,
	})
	if err != nil {
		return nil, err
	}

	err = b.articleContentRepo.Create(ctx, &entity.ArticleContentPo{
		ArticleID: articlePo.ID,
		Content:   content,
	})

	if err != nil {
		return nil, err
	}

	do := &ArticleWithContentDo{
		ID:         articlePo.ID,
		CreatedAt:  articlePo.CreatedAt,
		Pid:        articlePo.Pid,
		CategoryID: articlePo.CategoryID,
		Title:      articlePo.Title,
		Content:    content,
	}

	return do, nil
}

func (b *ArticleBiz) Update(ctx context.Context, id int, title string, content string) error {
	fields := make(map[string]interface{})
	fields["title"] = title
	err := b.articleRepo.Update(ctx, id, fields)
	if err != nil {
		return err
	}

	contentFields := make(map[string]interface{})
	contentFields["content"] = content
	err = b.articleContentRepo.Update(ctx, id, contentFields)
	if err != nil {
		return err
	}
	return nil
}

func (b *ArticleBiz) Delete(ctx context.Context, id int) error {
	err := b.articleContentRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = b.articleRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (b *ArticleBiz) Get(ctx context.Context, id int) (*ArticleWithContentDo, error) {
	articlePo, err := b.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	contentPo, err := b.articleContentRepo.GetByID(ctx, articlePo.ID)
	if err != nil {
		return nil, err
	}

	do := &ArticleWithContentDo{
		ID:         articlePo.ID,
		CreatedAt:  articlePo.CreatedAt,
		Pid:        articlePo.Pid,
		CategoryID: articlePo.CategoryID,
		Title:      articlePo.Title,
		Content:    contentPo.Content,
	}

	return do, nil
}

func convertArticlePosPtrToArticleDosPtr(articlePosPtrs []*entity.ArticlePo) []ArticleDo {
	var articleDosPtrs []ArticleDo
	for _, articlePoPtr := range articlePosPtrs {
		articleDo := ArticleDo{
			ID:         articlePoPtr.ID,
			CreatedAt:  articlePoPtr.CreatedAt,
			Pid:        articlePoPtr.Pid,
			CategoryID: articlePoPtr.CategoryID,
			Title:      articlePoPtr.Title,
			Children:   make([]ArticleDo, 0),
		}
		articleDosPtrs = append(articleDosPtrs, articleDo)
	}
	return articleDosPtrs
}

func buildTree(items []ArticleDo, pid int) []ArticleDo {
	var result []ArticleDo
	for _, item := range items {
		if item.Pid == pid {
			children := buildTree(items, item.ID)
			item.Children = children
			result = append(result, item)
		}
	}
	return result
}

func GroupArticleDosPtrByPid(articleDosPtrs []*ArticleDo) []*ArticleDo {
	pidMap := make(map[int][]*ArticleDo)

	// 将 ArticleDo 按 Pid 分组
	for _, articleDoPtr := range articleDosPtrs {
		pidMap[articleDoPtr.Pid] = append(pidMap[articleDoPtr.Pid], articleDoPtr)
	}

	// 构建嵌套的 ArticleDo 结构
	var result []*ArticleDo
	for _, articleDosPtrs := range pidMap {
		if articleDosPtrs[0].Pid == 0 {
			// 根节点
			result = append(result, articleDosPtrs...)
		} else {
			// 子节点
			for i := range articleDosPtrs {
				parentPid := articleDosPtrs[i].Pid
				parentIndex := findParentIndex(result, parentPid)
				if parentIndex != -1 {
					result[parentIndex].Children = append(result[parentIndex].Children, *articleDosPtrs[i])
				}
			}
		}
	}

	return result
}

func findParentIndex(articleDosPtrs []*ArticleDo, pid int) int {
	for i, articleDoPtr := range articleDosPtrs {
		if articleDoPtr.ID == pid {
			return i
		}
	}
	return -1
}
