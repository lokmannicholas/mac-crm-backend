package models

type Company struct {
	Code      string
	Name      *MultiLangText `gorm:"embedded;embeddedPrefix:name_"`
	ShortName *MultiLangText `gorm:"embedded;embeddedPrefix:short_name_"`
	Address   *MultiLangText `gorm:"embedded;embeddedPrefix:address_"`
	Email     string
	Phone     string
	TandC     map[string][]byte `gorm:"-"`
	Logo      *string
}
