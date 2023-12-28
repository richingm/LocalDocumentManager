package application

import "richingm/LocalDocumentManager/internal/domain"

type MountService struct {
}

type MountInfoDto struct {
	MountID        int
	ParentID       int
	MajorMinor     string
	Root           string
	MountPoint     string
	MountOptions   []string
	OptionalFields string
	FSType         string
	MountSource    string
	SuperOptions   []string
}

func NewMountService() *MountService {
	return &MountService{}
}

func (m *MountService) GetMountInfo(mountFilePath string) ([]MountInfoDto, error) {
	mountBiz := domain.NewMountBiz()
	mountInfoDos, err := mountBiz.GetMountInfo(mountFilePath)
	if err != nil {
		return nil, err
	}
	res := make([]MountInfoDto, 0)
	for _, mountInfoDo := range mountInfoDos {
		res = append(res, MountInfoDto{
			MountID:        mountInfoDo.MountID,
			ParentID:       mountInfoDo.ParentID,
			MajorMinor:     mountInfoDo.MajorMinor,
			Root:           mountInfoDo.Root,
			MountPoint:     mountInfoDo.MountPoint,
			MountOptions:   mountInfoDo.MountOptions,
			OptionalFields: mountInfoDo.OptionalFields,
			FSType:         mountInfoDo.FSType,
			MountSource:    mountInfoDo.MountSource,
			SuperOptions:   mountInfoDo.SuperOptions,
		})
	}
	return res, nil
}

func (m *MountService) GetDockerPath(mountFilePath string) (map[string]string, error) {
	mountBiz := domain.NewMountBiz()
	return mountBiz.GetDockerPath(mountFilePath)
}
