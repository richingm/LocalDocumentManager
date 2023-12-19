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

type NoteTreeJson struct {
}

func (b *NoteTreeService) GetTree(config configs.Config) []configs.NoteGroup {
	return config.Notes
}

func (b *NoteTreeService) GetDirAndName(config configs.Config, noteKey string) (string, string, error) {
	for _, item := range config.Notes {
		if len(item.Children) > 0 {
			for _, val := range item.Children {
				if noteKey == val.NoteKey {
					return val.NoteName, val.Dir, nil
				}
			}
		}
	}
	return "", "", errors.New("不存在")
}
