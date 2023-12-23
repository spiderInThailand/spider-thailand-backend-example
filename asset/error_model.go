package asset

type RootError struct {
	GeneralSystemError     ErrorCode `mapstructure:"general_system_error" json:"general_system_error"`
	InvalidLoginAccount    ErrorCode `mapstructure:"invalid_login_account" json:"invalid_login_account"`
	DecriptionError        ErrorCode `mapstructure:"decription_error" json:"decription_error"`
	GenerateRSAError       ErrorCode `mapstructure:"generate_rsa_error" json:"generate_rsa_error"`
	ErrorSpiderDB          ErrorCode `mapstructure:"error_spider_db" json:"error_spider_db"`
	ErrorTempDB            ErrorCode `mapstructure:"error_temp_db" json:"error_temp_db"`
	PasswordQualifyError   ErrorCode `mapstructure:"password_qualify_error" json:"password_qualify_error"`
	UsernameQualifyError   ErrorCode `mapstructure:"username_qualify_error" json:"username_qualify_error"`
	HashingError           ErrorCode `mapstructure:"hashing_error" json:"hashing_error"`
	PasswordMatchingError  ErrorCode `mapstructure:"password_matching_error" json:"password_matching_error"`
	UserNotLogin           ErrorCode `mapstructure:"user_not_login" json:"user_not_login"`
	InvalidImageType       ErrorCode `mapstructure:"invalid_image_type" json:"invalid_image_type"`
	SpiderNotFound         ErrorCode `mapstructure:"spider_not_found" json:"spider_not_found"`
	InsufficientUserRights ErrorCode `mapstructure:"insufficient_user_rights" json:"insufficient_user_rights"`
	DeleteSpiderFailed     ErrorCode `mapstructure:"delete_spider_failed" json:"delete_spider_failed"`
	GeographiesNotFound    ErrorCode `mapstructure:"geographies_not_found" json:"geographies_not_found"`
	RequestDataFail        ErrorCode `mapstructure:"request_data_fail" json:"request_data_fail"`
	RequestDataNotFound    ErrorCode `mapstructure:"request_data_not_found" json:"request_data_not_found"`
}

type ErrorCode struct {
	StatusCode     int    `mapstructure:"status_code" json:"status_code"`
	ErrorCode      string `mapstructure:"error_code" json:"error_code"`
	ErrorMessageTH string `mapstructure:"error_message_th" json:"error_message_th"`
	ErrorMessageEN string `mapstructure:"error_message_en" json:"error_message_en"`
}
