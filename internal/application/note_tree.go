package application

import "richingm/LocalDocumentManager/configs"

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
