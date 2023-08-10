package file

import (
	"github.com/CanftIn/gothafoss/pkg/im/config"
	"github.com/CanftIn/gothafoss/pkg/log"
)

// File 文件操作
type File struct {
	ctx *config.Context
	log.Log
	service IService
}

// New New
func New(ctx *config.Context) *File {
	return &File{
		ctx:     ctx,
		Log:     log.NewTLog("File"),
		service: NewService(ctx),
	}
}
