package db

const (
	// Account
	GetUserLogin         = "SELECT userid, email, password, active FROM users WHERE email=$1"
	UpdateUserPassword   = "UPDATE users SET password=$1 WHERE email=$2"
	GetPasswordDetails   = "SELECT userid, password FROM Users WHERE userid=$1"
	UpdatePassword       = "UPDATE users SET password=$1 WHERE userid=$2"
	UpdateUserStatus     = "UPDATE users SET active=$1 WHERE userid=$2"
	AddNewUser           = "INSERT INTO users (email, password, fullname, dob, country) VALUES ($1, $2, $3, $4, $5)"
	AddActivationCode    = "INSERT INTO activation (code, email) VALUES ($1, $2)"
	GetActivationCode    = "SELECT code, email FROM activation WHERE code=$1"
	RemoveActivationCode = "DELETE FROM activation WHERE code=$1 AND email=$2"
	ActivateUser         = "UPDATE users SET active=true WHERE email=$1"

	// Friends
	AddFriendRequest             = "INSERT INTO friendsrequest (requesterid, requestedid) VALUES ($1, $2)"
	CheckIfUserHasSpecificFriend = "SELECT 1 FROM friends WHERE (userone=$1 AND usertwo=$2) OR (userone=$2 AND usertwo=$1)"
	RemoveFromFriendsRequest     = "DELETE FROM friendsrequest WHERE requesterid=$1 AND requestedid=$2"
	DeclineFriendRequest         = "DELETE FROM friendsrequest WHERE requesterid=$1 AND requestedid=$2"
	RemoveFriendRequest          = "DELETE FROM friendsrequest WHERE requesterid=$1 AND requestedid=$2"
	AddNewFriend                 = "INSERT INTO friends (userone, usertwo) VALUES ($1, $2)"
	RemoveFriend                 = "DELETE FROM friends WHERE (userone=$1 AND usertwo=$2) OR (userone=$2 AND usertwo=$1)"
	GetAllFriends                = "SELECT u.userid, u.fullname, u.profilepic FROM friends f JOIN users u ON f.userone = u.userid OR f.usertwo = u.userid WHERE (f.userone = $1 OR f.usertwo = $1) AND u.userid != $1"

	RemoveTrip     = "DELETE FROM trips WHERE id=$1 AND userid=$2"
	GetUserByEmail = "SELECT email FROM users WHERE email=$1"
)
