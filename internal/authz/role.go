package authz

import (
	"encoding/json"
	"os"
)

type Resources []string

type Actions map[string]Resources

type Roles map[string]Actions

func LoadRoles() (Roles, error) {
	var roles Roles

	f, err := os.Open("roles.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
