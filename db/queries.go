package db

const (
	// Userdate
	UpdateUserCountry = "UPDATE users SET country=$1 WHERE userid=$2"

	// Create Account
	AddNewUser        = "INSERT INTO users (email, password, fullname, dob, country) VALUES ($1, $2, $3, $4, $5)"
	AddUserFlags      = "INSERT INTO userflags (userid) VALUES ((SELECT userid FROM users WHERE email = $1))"
	AddUserCounters   = "INSERT INTO usercounters (userid) VALUES ((SELECT userid FROM users WHERE email = $1))"
	AddActivationCode = "INSERT INTO activation (code, email) VALUES ($1, $2)"
	GetUserAccount    = "SELECT userid, email, password, active FROM users WHERE email=$1"

	// Login
	GetUserLogin       = "SELECT u.userid, u.email, u.password, u.active, uc.loginattempt, u.fullname FROM users u JOIN usercounters uc ON u.userid = uc.userid WHERE u.email=$1"
	UpdateLoginCounter = "UPDATE usercounters SET loginattempt = loginattempt + 1 WHERE userid = $1"
	ResetLoginCounter  = "UPDATE usercounters SET loginattempt = 0 WHERE userid = $1"

	// Password
	UpdateUserPassword = "UPDATE users SET password=$1 WHERE email=$2"
	GetPasswordDetails = "SELECT userid, password FROM Users WHERE userid=$1"
	EmailExists        = "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"
	UpdatePassword     = "UPDATE users SET password=$1 WHERE userid=$2"

	// User Status
	UpdateUserActiveStatus  = "UPDATE users SET active=$1 WHERE userid=$2"
	UpdateUserPrivacyStatus = "UPDATE userflags SET private=$1 WHERE userid=$2"
	GetActivationCode       = "SELECT code, email FROM activation WHERE code=$1"
	RemoveActivationCode    = "DELETE FROM activation WHERE code=$1 AND email=$2"
	ActivateUser            = "UPDATE users SET active=true WHERE email=$1"

	// Friend Requests
	AddFriendRequest         = "INSERT INTO friendsrequest (requesterid, requestedid) VALUES ($1, $2)"
	CheckIfUserHasFriend     = "SELECT 1 FROM friends WHERE (userone=$1 AND usertwo=$2) OR (userone=$2 AND usertwo=$1)"
	RemoveFromFriendsRequest = "DELETE FROM friendsrequest WHERE requesterid=$1 AND requestedid=$2"
	DeclineFriendRequest     = "DELETE FROM friendsrequest WHERE requesterid=$1 AND requestedid=$2"
	RemoveFriendRequest      = "DELETE FROM friendsrequest WHERE requesterid=$1 AND requestedid=$2"
	AddNewFriend             = "INSERT INTO friends (userone, usertwo) VALUES ($1, $2)"
	RemoveFriend             = "DELETE FROM friends WHERE (userone=$1 AND usertwo=$2) OR (userone=$2 AND usertwo=$1)"

	// Friends
	GetAllFriends = "SELECT u.userid, u.fullname, u.profilepic FROM friends f JOIN users u ON f.userone = u.userid OR f.usertwo = u.userid WHERE (f.userone = $1 OR f.usertwo = $1) AND u.userid != $1"

	// Pinned
	TripBelongsToUser = "SELECT 1 FROM trips WHERE userid=$1 AND tripid=$1"
	RemovePinned      = "DELETE FROM pinned WHERE userid=$1 AND tripid=$2"
	AddPinned         = "INSERT INTO pinned (userid, tripid) VALUES ($1, $2)"

	// Trips
	RemoveTrip = "DELETE FROM trips WHERE id=$1 AND userid=$2"

	// Countries
	GetAllCountries = "SELECT id, iso, %s FROM countries ORDER BY %s"
)
