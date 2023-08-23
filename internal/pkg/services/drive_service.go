package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	appclient "github.com/anjolaoluwaakindipe/testcli/internal/pkg"
	"github.com/anjolaoluwaakindipe/testcli/internal/pkg/utils"
	"google.golang.org/api/drive/v3"
)

type DriveService interface {
	GetDirectoryList(string) ([]FileInfo, error)
}

type GoogleDriveService struct {
	googleDrive *drive.Service
}

func (gds *GoogleDriveService) GetDirectoryList(folderId string) ([]FileInfo, error) {

	folderId = strings.TrimSpace(folderId)

	if folderId == "" {
		folderId = "root"
	}

	query := fmt.Sprintf("'%s' in parents", folderId)

	var fileInfos []FileInfo = make([]FileInfo, 0)
	ctx := context.Background()

	var fileList *drive.FileList

	err := gds.googleDrive.Files.List().Corpora("allDrives").Q(query).OrderBy("folder,name").SupportsAllDrives(true).Spaces("drive").Fields("files(id,name,size,mimeType)").IncludeItemsFromAllDrives(true).IncludePermissionsForView("published").PageSize(100).Pages(ctx, func(fl *drive.FileList) error {
		fileList = fl
		return nil
	})

	if err != nil {
		return fileInfos, err
	}

	for _, file := range fileList.Files {
		var fileName = file.Name
		var documentType utils.DocumentType = utils.File
		if strings.Contains(file.MimeType, "application/vnd.google-apps.folder") {
			fileName += "/"
			documentType = utils.Folder
		} else if strings.Contains(file.MimeType, "application/vnd.google-apps.shortcut") {
			fileName += "/"
			documentType = utils.Shortcut
		}
		fileInfos = append(fileInfos, FileInfo{Id: file.Id, MimeType: file.MimeType, Name: fileName, DocumentType: documentType, Size: int(file.Size)})
	}

	return fileInfos, nil
}

func InitDriveService() DriveService {
	drive, err := appclient.NewDriveClient()
	if err != nil {
		log.Panic("Google drive error: &v", err)
	}
	return &GoogleDriveService{googleDrive: drive}
}

type FileInfo struct {
	Id           string
	Name         string
	Size         int
	MimeType     string
	DocumentType utils.DocumentType
}
