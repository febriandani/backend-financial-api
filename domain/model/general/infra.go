package general

type SectionService struct {
	// AppAccount
	AppName         string `json:"APP_NAME"`
	AppEnvirontment string `json:"APP_ENV"`
	AppURL          string `json:"APP_URL"`
	AppPort         string `json:"APP_PORT"`
	AppSecretKey    string `json:"APP_KEY"`

	// RouteAccount
	RouteMethods string `json:"ROUTES_METHODS"`
	RouteHeaders string `json:"ROUTES_HEADERS"`
	RouteOrigins string `json:"ROUTES_ORIGIN"`

	// DatabaseAccount
	// Read
	DatabaseReadUsername     string `json:"DATABASE_READ_USERNAME"`
	DatabaseReadPassword     string `json:"DATABASE_READ_PASSWORD"`
	DatabaseReadURL          string `json:"DATABASE_READ_URL"`
	DatabaseReadPort         string `json:"DATABASE_READ_PORT"`
	DatabaseReadDBName       string `json:"DATABASE_READ_NAME"`
	DatabaseReadMaxIdleConns string `json:"DATABASE_READ_MAXIDLECONNS"`
	DatabaseReadMaxOpenConns string `json:"DATABASE_READ_MAXOPENCONNS"`
	DatabaseReadMaxLifeTime  string `json:"DATABASE_READ_MAXLIFETIME"`
	DatabaseReadTimeout      string `json:"DATABASE_READ_TIMEOUT"`
	DatabaseReadSSLMode      string `json:"DATABASE_READ_SSL_MODE"`

	// Write
	DatabaseWriteUsername     string `json:"DATABASE_WRITE_USERNAME"`
	DatabaseWritePassword     string `json:"DATABASE_WRITE_PASSWORD"`
	DatabaseWriteURL          string `json:"DATABASE_WRITE_URL"`
	DatabaseWritePort         string `json:"DATABASE_WRITE_PORT"`
	DatabaseWriteDBName       string `json:"DATABASE_WRITE_NAME"`
	DatabaseWriteMaxIdleConns string `json:"DATABASE_WRITE_MAXIDLECONNS"`
	DatabaseWriteMaxOpenConns string `json:"DATABASE_WRITE_MAXOPENCONNS"`
	DatabaseWriteMaxLifeTime  string `json:"DATABASE_WRITE_MAXLIFETIME"`
	DatabaseWriteTimeout      string `json:"DATABASE_WRITE_TIMEOUT"`
	DatabaseWriteSSLMode      string `json:"DATABASE_WRITE_SSL_MODE"`

	// Authorization
	// JWT
	AuthorizationJWTIsActive             string `json:"AUTHORIZATION_JWT_IS_ACTIVE"`
	AuthorizationJWTAccessTokenSecretKey string `json:"AUTHORIZATION_JWT_ACCESS_TOKEN_SECRET_KEY"`
	AuthorizationJWTAccessTokenDuration  string `json:"AUTHORIZATION_JWT_ACCESS_TOKEN_DURATION"`

	// Key Account
	KeyAccountUser string `json:"KEY_USER"`
}

type AppService struct {
	App           AppAccount   `json:",omitempty"`
	Route         RouteAccount `json:",omitempty"`
	Database      Database     `json:",omitempty"`
	Authorization AuthAccount  `json:",omitempty"`
	KeyData       KeyAccount   `json:",omitempty"`
}

type AppAccount struct {
	Name         string `json:",omitempty"`
	Environtment string `json:",omitempty"`
	URL          string `json:",omitempty"`
	Port         string `json:",omitempty"`
	SecretKey    string `json:",omitempty"`
}

type RouteAccount struct {
	Methods []string `json:",omitempty"`
	Headers []string `json:",omitempty"`
	Origins []string `json:",omitempty"`
}

type Database struct {
	Read  DBDetail `json:",omitempty"`
	Write DBDetail `json:",omitempty"`
}

type DBDetail struct {
	Username     string `json:",omitempty"`
	Password     string `json:",omitempty"`
	URL          string `json:",omitempty"`
	Port         string `json:",omitempty"`
	DBName       string `json:",omitempty"`
	MaxIdleConns int    `json:",omitempty"`
	MaxOpenConns int    `json:",omitempty"`
	MaxLifeTime  int    `json:",omitempty"`
	Timeout      string `json:",omitempty"`
	SSLMode      string `json:",omitempty"`
}

type AuthAccount struct {
	JWT JWTCredential `json:",omitempty"`
}

type JWTCredential struct {
	IsActive              bool   `json:",omitempty"`
	AccessTokenSecretKey  string `json:",omitempty"`
	AccessTokenDuration   int    `json:",omitempty"`
	RefreshTokenSecretKey string `json:",omitempty"`
	RefreshTokenDuration  int    `json:",omitempty"`
}

type KeyAccount struct {
	User string `json:",omitempty"`
}
