package models

type TandC struct {
	Version int    `gorm:"type:varchar(4);primaryKey;" json:"version"`
	Lang    string `gorm:"type:varchar(4);primaryKey;" json:"lang"`
	Content []byte `gorm:"type:TEXT;" json:"content"`
}

func (TandC) TableName() string {
	return "tandcs"
}

// func (t *TandC) BeforeCreate(tx *gorm.DB) (err error) {
// 	b, err := json.Marshal(t.Content)
// 	if err != nil {
// 		return err
// 	}
// 	t.Content = string(b)
// 	return
// }
