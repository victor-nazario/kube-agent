package authz

import (
	"bytes"
	"embed"
	"encoding/json"
)

//go:embed roles.json
var embedFS embed.FS

type Resources []string

type Actions map[string]Resources

type Roles map[string]Actions

func LoadRoles() (Roles, error) {
	var roles Roles

	r, err := embedFS.ReadFile("roles.json")
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(bytes.NewReader(r)).Decode(&roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
