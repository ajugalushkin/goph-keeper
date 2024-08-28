/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/internal/app/client/auth"
	"github.com/ajugalushkin/goph-keeper/internal/dto/models"
)

func newRegisterCmd(clnt *auth.Client) *cobra.Command {
	param := models.LoginParam{}

	registerCmd := &cobra.Command{
		Use:   "register -u username -p password",
		Short: "Register new user on server",
		Long:  `Register new user on server`,

		Run: func(cmd *cobra.Command, args []string) {
			runRegister(clnt, param, cmd)
		},
	}
	registerCmd.Flags().StringVarP(&param.User, "username", "u", "", "login of user")
	registerCmd.Flags().StringVarP(&param.Password, "password", "p", "", "password of user")

	return registerCmd
}

func runRegister(client *auth.Client, param models.LoginParam, cmd *cobra.Command) {
	err := client.Register(context.Background(), param.User, param.Password)
	if err != nil {
		return
	}
}
