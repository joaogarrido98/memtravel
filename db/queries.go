package db

const (
	GetUserLogin             = "SELECT id, email, password, active FROM users WHERE email=$1"
	UpdateUserPassword       = "UPDATE users SET password=$1 WHERE email=$2"
	GetPasswordDetails       = "SELECT id, password FROM Users WHERE id=$1"
	UpdatePassword           = "UPDATE users SET password=$1 WHERE id=$2"
	UpdateUserStatus         = "UPDATE users SET active=$1 WHERE id=$2"
	RemoveTrip               = "DELETE FROM trips WHERE id=$1 AND userid=$2"
	GetUserByEmail           = "SELECT email FROM users WHERE email=$1"
	AddNewUser               = "INSERT INTO users (email, password, fullname, dob, country) VALUES ($1, $2, $3, $4, $5)"
	AddActivationCode        = "INSERT INTO activation (code, email) VALUES ($1, $2)"
	GetActivationCode        = "SELECT code, email FROM activation WHERE code=$1"
	RemoveActivationCode     = "DELETE FROM activation WHERE code=$1 AND email=$2"
	ActivateUser             = "UPDATE users SET active=true WHERE email=$1"
	AddFriendRequest         = "INSERT INTO friendsrequest (requesterid, requestedid) VALUES ($1, $2)"
	GetUsersSpecificFriend   = "SELECT 1 FROM friends WHERE (userone=$1 AND usertwo=$2) OR (userone=$2 AND usertwo=$1)"
	RemoveFromFriendsRequest = "DELETE FROM friendsrequest WHERE requesterid=$1 AND requestedid=$2"
	DeclineFriendRequest     = "DELETE FROM friendsrequest WHERE requesterid=$1 AND requestedid=$2"
	RemoveFriendRequest      = "DELETE FROM friendsrequest WHERE requesterid=$1 AND requestedid=$2"
	AddNewFriend             = "INSERT INTO friends (userone, usertwo) VALUES ($1, $2)"
	RemoveFriend             = "DELETE FROM friends where userone=$1, $2)"
)
