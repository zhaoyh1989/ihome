package models

import (
	"github.com/beego/beego/v2/core/logs"
	fdfs "github.com/tedcy/fdfs_client"
	"path"
)

func UploadByFileName(file string) {
	client, err := fdfs.NewClientWithConfig("conf/fdfs.conf")
	defer client.Destory()
	if err != nil {
		logs.Error("fdfs.NewClientWithConfig failed err=", err)
		return
	}
	fileId, err := client.UploadByFilename(file)
	if err != nil {
		logs.Error(" client.UploadByFilename(file) failed err=", err)
		return
	}
	logs.Info(fileId)
	//if err := client.DownloadToFile(fileId, "tempFile", 0, 0); err != nil {
	//	logs.Error(" client.UploadByFilename(file) failed err=", err)
	//	return
	//}
	//if buffer, err := client.DownloadToBuffer(fileId, 0, 19); err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Println(string(buffer))
	//}
	//if err := client.DeleteFile(fileId); err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
}

func UploadByBuffer(fb []byte, filename string) (string, error) {
	client, err := fdfs.NewClientWithConfig("conf/fdfs.conf")
	defer client.Destory()
	if err != nil {
		logs.Error("fdfs.NewClientWithConfig failed err=", err)
		return "", err
	}
	logs.Debug("文件名：", filename)
	//split := strings.Split(filename, ".")
	//filesuffix := split[len(split)-1]
	filesuffix := path.Ext(filename)
	logs.Debug("上传文件后缀", filesuffix)
	fileurl, err := client.UploadByBuffer(fb, filesuffix[1:])
	if err != nil {
		logs.Error(" client.UploadByFilename(file) failed err=", err)
		return "", err
	}
	logs.Info(fileurl)
	return fileurl, nil
}
