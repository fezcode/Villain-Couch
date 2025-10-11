package fs

import "os"

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return info.IsDir()
}
