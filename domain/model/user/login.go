package user

type LoginRequest struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() map[string]string {
	if r.Email == "" {
		return map[string]string{
			"en": "Email or username cannot be empty",
			"id": "Email atau username tidak boleh kosong",
		}
	}

	// _, err := mail.ParseAddress(r.Email)
	// if err != nil {
	// 	return map[string]string{
	// 		"en": "Incorrect email format",
	// 		"id": "Format email salah",
	// 	}
	// }

	if r.Password == "" {
		return map[string]string{
			"en": "Password cannot be empty",
			"id": "Kaata sandi tidak boleh kosong",
		}
	}

	return nil
}

type LoginResponse struct {
	Token       JWTAccess `json:"token"`
	UserID      int64     `json:"user_id"`
	NamaLengkap string    `json:"nama_lengkap"`
}

type JWTAccess struct {
	AccessToken        string `json:"access"`
	AccessTokenExpired string `json:"access_expired"`
}

type CredentialData struct {
	ID       int64
	Email    string
	Fullname string
}
