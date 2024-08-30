package language

const (
	EmptyOldNewPassword     = "EmptyOldNewPassword"
	PasswordInvalid         = "PasswordInvalid"
	EmptyEmailPassword      = "EmptyEmailPassword"
	EmptyEmail              = "EmptyEmail"
	InactiveUser            = "InactiveUser"
	NewPasswordInvalid      = "NewPasswordInvalid"
	PasswordChanged         = "PasswordChanged"
	ChagePasswordInvalid    = "ChangePasswordInvalid"
	PasswordRecover         = "PasswordRecover"
	PasswordRecoverySuccess = "PasswordRecoverySuccess"
	AccountClose            = "AccountClose"
)

var All = map[string]map[string]string{
	"1": en,
	"2": pt,
}

var en = map[string]string{
	EmptyOldNewPassword:     "Your current or new password cannot be empty",
	PasswordInvalid:         "Invalid email or password",
	EmptyEmailPassword:      "Your email or password cannot be empty",
	InactiveUser:            "This user account has been closed",
	NewPasswordInvalid:      "The new password is not valid",
	PasswordChanged:         "Password has been changed",
	ChagePasswordInvalid:    "You current password is invalid",
	PasswordRecover:         "Password Recovery Request",
	EmptyEmail:              "Your email cannot be empty",
	PasswordRecoverySuccess: "Your new password was sent to your email, please change it as soon as you login again.",
	AccountClose:            "We are sorry to see you leave.",
}

var pt = map[string]string{
	EmptyOldNewPassword:     "A sua senha atual ou nova não podem estar vazias",
	PasswordInvalid:         "Email ou senha incorretos",
	EmptyEmailPassword:      "O seu email ou senha não podem estar vazios",
	InactiveUser:            "Esta conta está atualmente encerrada",
	NewPasswordInvalid:      "A sua nova senha não é valida",
	PasswordChanged:         "Senha atualizada",
	ChagePasswordInvalid:    "A sua senha atual não é valida",
	PasswordRecover:         "Pedido the recuperacão de senha",
	EmptyEmail:              "O seu email não pode estar vazio",
	PasswordRecoverySuccess: "A sua nova senha foi enviada para o seu email, por favor mude a senha assim que iniciar sessão outra vez.",
	AccountClose:            "Estamos tristes por fechar a conta.",
}

func GetTranslation(languageID string, translationKey string) string {
	return All[languageID][translationKey]
}

func SupportedLanguage(languageID string) bool {
	_, supported := All[languageID]
	return supported
}
