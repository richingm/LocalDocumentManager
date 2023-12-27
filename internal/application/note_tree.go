package application

type NoteTreeService struct {
}

func NewNoteTreeService() *NoteTreeService {
	return &NoteTreeService{}
}

func (b *NoteTreeService) GetNote(menuDtos []MenuDto, noteKey string) *MenuDto {
	for _, menuDto := range menuDtos {
		note := loopGetNote(menuDto, noteKey)
		if note != nil {
			return note
		}
	}
	return nil
}

func loopGetNote(menuDto MenuDto, noteKey string) *MenuDto {
	if menuDto.MenuKey == noteKey {
		return &MenuDto{
			MenuKey:  menuDto.MenuKey,
			MenuName: menuDto.MenuName,
			DirName:  menuDto.DirName,
			DirPath:  menuDto.DirPath,
		}
	}
	if len(menuDto.Children) > 0 {
		for _, val := range menuDto.Children {
			note := loopGetNote(val, noteKey)
			if note != nil {
				return note
			}
		}
	}
	return nil
}
