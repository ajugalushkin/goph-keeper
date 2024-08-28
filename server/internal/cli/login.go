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

func newLoginCmd(clnt *auth.Client) *cobra.Command {
	param := models.LoginParam{}

	registerCmd := &cobra.Command{
		Use:   "login -u username -p password",
		Short: "Login new user on KeepPas server",
		Long:  `Login new user on KeepPas server and get authenticated token.`,

		Run: func(cmd *cobra.Command, args []string) {
			runLogin(clnt, param, cmd)
		},
	}
	registerCmd.Flags().StringVarP(&param.User, "username", "u", "", "login of user")
	registerCmd.Flags().StringVarP(&param.Password, "password", "p", "", "password of user")

	return registerCmd
}

func runLogin(client *auth.Client, param models.LoginParam, cmd *cobra.Command) {
	_, err := client.Login(context.Background(), param.User, param.Password)
	if err != nil {
		return
	}
}
