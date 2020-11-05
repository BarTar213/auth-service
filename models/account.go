package models

type AccountInfo struct {
	ID    int    `header:"X-Account-Id" binding:"required"`
	Login string `header:"X-Account" binding:"required"`
	Role  string `header:"X-Role" binding:"required"`
}
