package db

const (
	GetUserLogin         = "SELECT id, email, password, active FROM users WHERE email=$1"
	UpdateUserPassword   = "UPDATE users SET password=$1 WHERE email=$2"
	GetPasswordDetails   = "SELECT id, password FROM Users WHERE id=$1"
	UpdatePassword       = "UPDATE users SET password=$1 WHERE id=$2"
	UpdateUserStatus     = "UPDATE users SET active=$1 WHERE id=$2"
	RemoveTrip           = "DELETE FROM trips WHERE id=$1 AND userid=$2"
	GetUserByEmail       = "SELECT email FROM users WHERE email=$1"
	AddNewUser           = "INSERT INTO users (email, password, fullname, dob, country) VALUES ($1, $2, $3, $4, $5)"
	AddActivationCode    = "INSERT INTO activation (code, email) VALUES ($1, $2)"
	GetActivationCode    = "SELECT code, email FROM activation WHERE code=$1"
	RemoveActivationCode = "DELETE FROM activation WHERE code=$1 AND email=$2"
	ActivateUser         = "UPDATE users SET active=true WHERE email=$1"
)
