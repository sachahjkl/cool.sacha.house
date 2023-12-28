package handler

import (
	"sachahjkl/htmx_go/pkg/config"

	"gorm.io/gorm"
)

type handler struct {
	DB     *gorm.DB
	Config *config.Config
}
