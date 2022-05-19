package managers

type FieldOptionParam struct {
	Name *MultiLangText `json:"name,omitempty"`
}

type FieldOptionUpdate struct {
	Name   *MultiLangText `json:"name,omitempty"`
	ID     string         `json:"id,omitempty"`
	Action string         `json:"action,omitempty" validate:"required" enums:"CREATE,UPDATE,DELETE"`
}
