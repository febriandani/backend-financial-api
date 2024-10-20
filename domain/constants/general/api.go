package general

const (
	HandlerErrorAuthInvalid           string = "authorization invalid"
	HandlerErrorAuthInvalidID         string = "authorization tidak valid"
	HandlerErrorRequestDataNotValid   string = "request data not valid"
	HandlerErrorRequestDataNotValidID string = "data request tidak valid"
	HandlerErrorRequestDataEmpty      string = "request data empty"
	HandlerErrorRequestDataEmptyID    string = "data request kosong"
)

type CredentialData struct {
	ID        int64  `json:"id"`
	CompanyID int64  `json:"company_id"`
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
}
