package user

// This is an oversimplification for challenge purposes, just to
// provide the construct of a user with permissions

type User struct {
	UserName string   `json:"userName"`
	Roles    []string `json:"roles"`
}

type Users map[string]User

// ActiveUsers returns a list of all active users along with their current assigned roles
func ActiveUsers() Users {
	users := make(map[string]User)

	users["operant"] = User{
		UserName: "operant",
		Roles:    []string{"cluster-owner"},
	}

	users["reader"] = User{
		UserName: "reader",
		Roles:    []string{"cluster-reader"},
	}

	return users
}
