package helper

import "os"

// SubmitCodeSave 保存用户提交代码
func SubmitCodeSave(code []byte) (string, error) {
	// 创建以uuid作为名称的文件夹
	dirName := "code/" + GetUUID()
	filePath := dirName + "/main.go"
	// 创建文件夹
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	_, err = f.Write(code)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return filePath, nil
}
