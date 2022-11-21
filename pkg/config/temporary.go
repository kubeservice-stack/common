package config

import (
	"fmt"
)

// Temporary represents a logging configuration
type Temporary struct {
	MaxBufferSize int64  `toml:"max_buffer_size" json:"max_buffer_size" env:"TEMPORARY_MAX_BUFFER_SIZE"` // 最大使用内存空间, 超过时则转化成文件
	FileDir       string `toml:"file_dir" json:"file_dir" env:"TEMPORARY_FILE_DIR"`                      // 临时文件目录
	FilePattern   string `toml:"file_pattern" json:"file_pattern" env "TEMPORARY_FILE_PATTERN"`          // 临时文件名格式
	MaxUploadSize int64  `toml:"max_upload_size" json:"max_upload_size" env "TEMPORARY_MAX_UPLOAD_SIZE"` // 最大上传文件大小
}

func (l Temporary) TOML() string {
	return fmt.Sprintf(`
[temporary]
  ## 最大使用内存空间, 超过时则转化成文件, 默认是 5242880 byte = 5MB
  max_buffer_size = %d
  ## 上传文件 临时文件目录, 默认 /tmp
  file_dir = "%s"
  ## 上传文件临时文件名格式, 默认前缀 uploadd-*
  file_pattern = "%s"
  ## 上传文件临时文件名格式, 默认前缀 104857600 byte = 100MB
  max_upload_size = %d`,
		l.MaxBufferSize,
		l.FileDir,
		l.FilePattern,
		l.MaxUploadSize)
}

func (l Temporary) DefaultConfig() Temporary {
	l = Temporary{
		MaxBufferSize: 5242880,
		FileDir:       "/tmp",
		FilePattern:   "uploadd-*",
		MaxUploadSize: 104857600,
	}
	return l
}
