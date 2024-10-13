package language

const (
	PasswordInvalid         = "PasswordInvalid"
	BlockedLogin            = "BlockedLogin"
	InactiveUser            = "InactiveUser"
	PasswordChanged         = "PasswordChanged"
	ChagePasswordInvalid    = "ChangePasswordInvalid"
	PasswordRecover         = "PasswordRecover"
	PasswordRecoverySuccess = "PasswordRecoverySuccess"
	AccountClose            = "AccountClose"
	AccountCreated          = "AccountCreated"
	AccountExisting         = "AccountExisting"
	AccountNotExisting      = "AccountNotExisting"
	Welcome                 = "Welcome"
)

// All holds all language maps
var All = map[string]map[string]string{
	"1": en,
	"2": pt,
}

var en = map[string]string{
	PasswordInvalid:         "Invalid email or password",
	InactiveUser:            "This account is not active",
	PasswordChanged:         "Password has been changed",
	ChagePasswordInvalid:    "You current password is invalid",
	PasswordRecover:         "Password Recovery Request",
	PasswordRecoverySuccess: "If your account exists, your new password was sent to your email, please change it as soon as you login again.",
	AccountClose:            "We are sorry to see you leave.",
	AccountCreated:          "Your new account has been created, before logging in, please verify your account through the email you have received.",
	AccountExisting:         "Email is already in use",
	AccountNotExisting:      "Account does not exist",
	Welcome:                 "Memtravel welcomes you",
	BlockedLogin:            "Your account is currently locked",
}

var pt = map[string]string{
	PasswordInvalid:         "Email ou senha incorretos",
	InactiveUser:            "Esta conta está desativada",
	PasswordChanged:         "Senha atualizada",
	ChagePasswordInvalid:    "A sua senha atual não é valida",
	PasswordRecover:         "Pedido the recuperacão de senha",
	PasswordRecoverySuccess: "Se a sua conta existir, a sua nova senha foi enviada para o seu email, por favor mude a senha assim que iniciar sessão outra vez.",
	AccountClose:            "Estamos tristes por fechar a conta.",
	AccountCreated:          "A sua nova conta foi criada, antes the entrar, por favor verifique a sua conta usando o email que enviamos.",
	AccountExisting:         "Email ja se encontra em uso.",
	AccountNotExisting:      "Esta conta nao exist.",
	Welcome:                 "Bem-vindo a Memtravel",
	BlockedLogin:            "A sua conta está bloqueada",
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
