package language

const (
	PasswordInvalid         = "PasswordInvalid"
	InactiveUser            = "InactiveUser"
	PasswordChanged         = "PasswordChanged"
	ChagePasswordInvalid    = "ChangePasswordInvalid"
	PasswordRecover         = "PasswordRecover"
	PasswordRecoverySuccess = "PasswordRecoverySuccess"
	AccountClose            = "AccountClose"
	AccountCreated          = "AccountCreated"
	AccountExisting         = "AccountExisting"
	Welcome                 = "Welcome"
)

// All holds all language maps
var All = map[string]map[string]string{
	"1": en,
	"2": pt,
}

var en = map[string]string{
	PasswordInvalid:         "Invalid email or password",
	InactiveUser:            "This user account has been closed",
	PasswordChanged:         "Password has been changed",
	ChagePasswordInvalid:    "You current password is invalid",
	PasswordRecover:         "Password Recovery Request",
	PasswordRecoverySuccess: "Your new password was sent to your email, please change it as soon as you login again.",
	AccountClose:            "We are sorry to see you leave.",
	AccountCreated:          "Your new account has been created, before logging in, please verify your account through the email you have received.",
	AccountExisting:         "Email is already in use",
	Welcome:                 "Memtravel welcomes you",
}

var pt = map[string]string{
	PasswordInvalid:         "Email ou senha incorretos",
	InactiveUser:            "Esta conta está atualmente encerrada",
	PasswordChanged:         "Senha atualizada",
	ChagePasswordInvalid:    "A sua senha atual não é valida",
	PasswordRecover:         "Pedido the recuperacão de senha",
	PasswordRecoverySuccess: "A sua nova senha foi enviada para o seu email, por favor mude a senha assim que iniciar sessão outra vez.",
	AccountClose:            "Estamos tristes por fechar a conta.",
	AccountCreated:          "A sua nova conta foi criada, antes the entrar, por favor verifique a sua conta usando o email que enviamos.",
	AccountExisting:         "Email ja se encontra em uso.",
	Welcome:                 "Bem-vindo a Memtravel",
}

// GetTranslation retrieves a translation for a specific language id
func GetTranslation(languageID string, translationKey string) string {
	return All[languageID][translationKey]
}

// SupportedLanguage checks if a specific language id exists
func SupportedLanguage(languageID string) bool {
	_, supported := All[languageID]
	return supported
}
