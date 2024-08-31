package db

const (
	GetUserLogin       = "SELECT id, email, password, active FROM users WHERE email=$1"
	UpdateUserPassword = "UPDATE users SET password=$1 WHERE email=$2"
	GetPasswordDetails = "SELECT id, password FROM Users WHERE id=$1"
	UpdatePassword     = "UPDATE users SET password=$1 WHERE id=$2"
	UpdateUserStatus   = "UPDATE users SET active=$1 WHERE id=$2"
	RemoveTrip         = "DELETE FROM trips WHERE id=$1 AND userid=$2"
	GetUserByEmail     = "SELECT email FROM users WHERE email=$1"
)
