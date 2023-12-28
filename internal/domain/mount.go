package domain

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type MountBiz struct {
}

func NewMountBiz() *MountBiz {
	return &MountBiz{}
}

type MountInfoDo struct {
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

func (m *MountBiz) GetDockerPath(mountFilePath string) (map[string]string, error) {
	mountInfo, err := m.GetMountInfo(mountFilePath)
	if err != nil {
		return nil, err
	}
	res := make(map[string]string, 0)
	for _, do := range mountInfo {
		res[do.Root] = do.MountPoint
	}
	return res, nil
}

func (m *MountBiz) GetMountInfo(filePath string) ([]MountInfoDo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	var mountInfos []MountInfoDo
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 10 {
			mountInfo := MountInfoDo{
				MountID:        parseInt(fields[0]),
				ParentID:       parseInt(fields[1]),
				MajorMinor:     fields[2],
				Root:           fields[3],
				MountPoint:     fields[4],
				MountOptions:   strings.Split(fields[5], ","),
				OptionalFields: fields[6],
				FSType:         fields[7],
				MountSource:    fields[8],
				SuperOptions:   strings.Split(fields[9], ","),
			}
			if strings.HasPrefix(mountInfo.MountPoint, "/app") {
				mountInfos = append(mountInfos, mountInfo)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return mountInfos, nil
}

func parseInt(s string) int {
	var i int
	_, _ = fmt.Sscanf(s, "%d", &i)
	return i
}
