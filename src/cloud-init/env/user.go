package env

import "os/exec"

type User struct {
	Name            string `flag:""`
	Password        string `flag:"--password"`
	HomeDir         string `flag:"--home"`
	DefaultShell    string `flag:"--shell"`
	IsSystemAccount bool   `flag:"--system"`

	//PrimarGroup Group
	//Groups      []Group

	//ExpireDate
	//InactiveDate
}

func (usr *User) FlattenArguments() []string {
	args := []string{usr.Name}

	flattenedArgs := FlattenArguments(map[string]string{
		"home":     usr.HomeDir,
		"password": usr.Password,
		"shell":    usr.DefaultShell,
	})

	args = append(args, flattenedArgs...)

	if usr.IsSystemAccount {
		args = append(args, "--system")
	}

	return args
}

func NewUser(usr User) error {
	args := usr.FlattenArguments()

	cmd := exec.Command("useradd", args...)
	err := cmd.Run()

	return err
}
