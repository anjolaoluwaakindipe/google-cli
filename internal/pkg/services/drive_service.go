package services

import (
	"context"
	"log"

	appclient "github.com/anjolaoluwaakindipe/testcli/internal/pkg"
	"google.golang.org/api/drive/v3"
)

type DriveService interface {
	GetDirectoryList(string) ([]string, error)
}

type GoogleDriveService struct {
	googleDrive *drive.Service
}

func (gds *GoogleDriveService) GetDirectoryList(string) ([]string,error) {

		var fileNames [] string = make([]string, 0)
		ctx := context.Background()

		var fileList *drive.FileList

		err := gds.googleDrive.Files.List().Corpora("allDrives").Q("'root' in parents").SupportsAllDrives(true).Spaces("drive").IncludeItemsFromAllDrives(true).IncludePermissionsForView("published").PageSize(100).Pages(ctx, func(fl *drive.FileList) error {
			fileList = fl
			return nil
		})

		if err != nil {
			return fileNames, nil
		}

		for _, file := range fileList.Files {
			fileNames = append(fileNames, file.Name + "." + file.FileExtension)
		}
		return fileNames, nil
}

func InitDriveService() DriveService {
	drive, err := appclient.NewDriveClient()
	if err != nil {
		log.Panic("Google drive error: &v", err)
	}
	return &GoogleDriveService{googleDrive: drive}
}
