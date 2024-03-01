package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetAbsFilePath(folderPath, fileName string) (string, error) {
	absFolderPath, err := filepath.Abs(strings.TrimSpace(folderPath))
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	return filepath.Join(absFolderPath, strings.TrimSpace(fileName)), nil
}

func ReadFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func ParseBinaryFile(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var byteArray []string
	for _, b := range content {
		byteArray = append(byteArray, fmt.Sprintf("0x%02X", b))
	}

	return byteArray, nil
}

func SaveToFile(folderPath string, fileName string, content string) error {
	filePath, err := GetAbsFilePath(folderPath, fileName)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
			return err
		}
	}

	return os.WriteFile(filePath, []byte(content), 0o644)
}

func ToCArray(filePath string) (string, error) {
	byteArray, err := ParseBinaryFile(filePath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{ %s }", strings.Join(byteArray, ", ")), nil
}

func CheckPath(cmd string) (string, error) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to find %s: %v", cmd, err)
	}

	return path, nil
}

func Compile64(loaderPath string, make string) error {
	cmd := exec.Command(make, "x64")
	cmd.Dir = loaderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Compile32(loaderPath string, make string) error {
	cmd := exec.Command(make, "x86")
	cmd.Dir = loaderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CompileLoader(ldrPath string, arch string, mingw_make string) error {
	switch arch {
	case "x64":
		return Compile64(ldrPath, mingw_make)
	case "x86":
		return Compile32(ldrPath, mingw_make)
	default:
		return fmt.Errorf("unsupported architecture: %s", arch)
	}
}
