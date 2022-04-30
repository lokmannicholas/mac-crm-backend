package _const

type Permission string

const PERMISSION_CONSUMABLE_PRODUCT Permission = "CONSUMABLE_PRODUCT"
const PERMISSION_SINGLE_PRODUCT Permission = "SINGLE_PRODUCT"
const PERMISSION_SETTING Permission = "SETTING"
const PERMISSION_ACCOUNT Permission = "ACCOUNT"
const PERMISSION_APP Permission = "APP"
const PERMISSION_BRANCH Permission = "BRANCH"
const PERMISSION_CATEGORY Permission = "CATEGORY"
const PERMISSION_COMPANY Permission = "COMPANY"
const PERMISSION_CUSTOMER Permission = "CUSTOMER"
const PERMISSION_FEATURE Permission = "FEATURE"
const PERMISSION_INVOICE Permission = "INVOICE"
const PERMISSION_RENTAL_ORDER Permission = "RENTAL_ORDER"
const PERMISSION_RENTAL_ORDER_EVENT Permission = "RENTAL_ORDER_EVENT"
const PERMISSION_STORAGE Permission = "STORAGE"
const PERMISSION_STORAGE_EVENT Permission = "STORAGE_EVENT"
const PERMISSION_STORAGE_RECORDS Permission = "STORAGE_RECORDS"
const PERMISSION_REPORT Permission = "REPORT"
const PERMISSION_PAYMENT Permission = "PAYMENT"
const PERMISSION_ROLE Permission = "ROLE"
const PERMISSION_ENCRYPTION Permission = "ENCRYPTION"

func (p Permission) Create() string {
	return string(p) + ":C"
}
func (p Permission) Read() string {
	return string(p) + ":R"
}
func (p Permission) Update() string {
	return string(p) + ":U"
}
func (p Permission) ToString() string {
	return string(p)
}
