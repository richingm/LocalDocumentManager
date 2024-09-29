package application

import "richingm/LocalDocumentManager/internal/domain"

type BackupService struct {
}

func NewBackupService() *BackupService {
	return &BackupService{}
}

func (service *BackupService) Back() error {
	backupBiz := domain.NewBackupBiz()
	err := backupBiz.Save()
	if err != nil {
		return err
	}
	return nil
}
