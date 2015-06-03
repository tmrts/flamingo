package env

import "os/exec"

type User struct {
	Name string
	Password
	//Groups []Group

	//HomeDir
	//ExpireDate
	//InactiveDate
	//IsSystemAccount
}

func NewUser(usr User) error {
	args := FlattenArguments(map[string]string{
		"user-name": usr.Name,
		"password":  usr.Hash,
	})

	err := exec.Command("useradd", args).Run()
	if err != nil {
		return err
	}

	return nil
}
