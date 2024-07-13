package language

const (
	EmptyOldNewPassword = "EmptyOldNewPassword"
	PasswordInvalid     = "PasswordInvalid"
	EmptyEmailPassword  = "EmptyEmailPassword"
	InactiveUser        = "InactiveUser"
	NewPasswordInvalid  = "NewPasswordInvalid"
	PasswordChanged     = "PasswordChanged"
)

var All = map[string]map[string]string{
	"1": en,
	"2": pt,
}

var en = map[string]string{
	EmptyOldNewPassword: "Your current or new password cannot be empty",
	PasswordInvalid:     "Your current password is not correct",
	EmptyEmailPassword:  "Your email or password cannot be empty",
	InactiveUser:        "This user account has been closed",
	NewPasswordInvalid:  "The given password is not valid",
	PasswordChanged:     "Password has been changed",
}

var pt = map[string]string{
	EmptyOldNewPassword: "A sua senha atual ou nova não podem estar vazias",
	PasswordInvalid:     "A sua senha atual nao está correta",
	EmptyEmailPassword:  "O seu email ou sennha nao podem estar vazios",
	InactiveUser:        "Esta conta está atualmente encerrada",
	NewPasswordInvalid:  "A sua nova senha não é valida",
	PasswordChanged:     "Senha atualizada",
}

func GetTranslation(languageID string, translationKey string) string {
	return All[languageID][translationKey]
}
