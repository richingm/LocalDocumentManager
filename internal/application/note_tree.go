package application

import (
	"errors"
	"richingm/LocalDocumentManager/configs"
)

type NoteTreeService struct {
}

func NewNoteTreeService() *NoteTreeService {
	return &NoteTreeService{}
}

func (b *NoteTreeService) GetTree(config configs.Config) []configs.NoteGroup {
	return config.Notes
}

func (b *NoteTreeService) GetNote(config configs.Config, noteKey string) (configs.NoteChild, error) {
	for _, item := range config.Notes {
		if len(item.Children) > 0 {
			for _, val := range item.Children {
				if noteKey == val.NoteKey {
					return val, nil
				}
			}
		}
	}
	return configs.NoteChild{}, errors.New("不存在")
}
