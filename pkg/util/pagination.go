package util

import (
	"math"

	"dmglab.com/mac-crm/pkg/config"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"gorm.io/gorm"
)

type Pagination struct {
	TotalCount int64 `json:"totalCount"`
	PageCount  int   `json:"pageCount"`
	Limit      int   `json:"limit,omitempty;query:limit"`
	Page       int   `json:"page,omitempty;query:page"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		limit, ok := config.GetConfig().Setting[_const.ROW_LIMIT]
		if ok {
			p.Limit = StrMustToInt(limit)
		} else {
			p.Limit = 100
		}

	}
	return p.Limit
}
func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}
func PaginationScope(v interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {

	var totalRows int64
	db.Model(v).Count(&totalRows)
	pagination.TotalCount = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.PageCount = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
