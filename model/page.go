package model

import (
	"gorm.io/gorm"
	db "haibara/database"
	"haibara/util"
	"math"
)

type Page struct {
	PageNum  int `form:"pageNum" json:"pageNum"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type PageResult struct {
	List any `json:"list"`
	Page
	Total int `json:"total"`
	Pages int `json:"pages"`
}

func Paginate[V any](page Page, condition V) PageResult {
	page = setPage(page)
	pageResult := PageResult{
		Page: page,
		List: []struct{}{},
	}
	var count int64
	db.GORM.Model(condition).Count(&count)
	pageResult.Total = int(count)
	pages := int(math.Ceil(float64(count) / float64(pageResult.PageSize)))
	pageResult.Pages = pages
	if count == 0 || pageResult.PageNum > pages {
		return pageResult
	}
	var data []V
	db.GORM.Scopes(paginate(page)).Where(&condition).Find(&data)
	pageResult.List = data
	return pageResult
}

func setPage(page Page) Page {
	pageNum := page.PageNum
	pageSize := page.PageSize
	pageNum = util.ConditionalExpression(pageNum <= 0, 1, pageNum)
	pageSize = util.ConditionalExpression(pageSize <= 0, 10, util.ConditionalExpression(pageSize > 100, 100, pageSize))
	page.PageNum = pageNum
	page.PageSize = pageSize
	return page
}

func paginate(page Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page.PageNum - 1) * page.PageSize
		return db.Offset(offset).Limit(page.PageSize)
	}
}
