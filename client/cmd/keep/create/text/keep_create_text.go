package text

import (
	"context"
	"fmt"
	"github.com/ajugalushkin/goph-keeper/client/internal/config"
	"github.com/ajugalushkin/goph-keeper/client/internal/secret"
	"github.com/ajugalushkin/goph-keeper/client/internal/token_cache"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/client/internal/app"
	"github.com/ajugalushkin/goph-keeper/client/internal/logger"
	"github.com/ajugalushkin/goph-keeper/client/internal/vaulttypes"
	v1 "github.com/ajugalushkin/goph-keeper/gen/keeper/v1"
)

// NewCommand creates a cobra.Command for creating a text secret.
// The command is responsible for handling the execution of the "text" subcommand of the "keep create" command.
// It reads the required flags "name" and "data", creates a text secret, encrypts it, and sends a request to the goph-keeper server to store the secret.
// If any error occurs during the process, it logs the error using the slog package.
//
// Parameters:
// - None
//
// Return:
// - A pointer to the cobra.Command object representing the "text" subcommand.
func NewCommand() *cobra.Command {
	const op = "keep_create_text"

	cmd := &cobra.Command{
		Use:   "text",
		Short: "Create text secret",
		Run:   keepCreateTextCmdRun,
	}

	cmd.Flags().String("name", "", "Secret name")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}
	cmd.Flags().String("data", "", "Text data")
	if err := cmd.MarkFlagRequired("data"); err != nil {
		slog.Error("Error setting flag: ",
			slog.String("op", op),
			slog.String("error", err.Error()))
	}

	return cmd
}

// keepCreateTextCmdRun is responsible for handling the execution of the "text" subcommand of the "keep create" command.
// This function reads the required flags "name" and "data", creates a text secret, encrypts it, and sends a request to the goph-keeper server to store the secret.
// If any error occurs during the process, it logs the error using the slog package.
//
// Parameters:
// - cmd: A pointer to the cobra.Command object representing the "text" subcommand.
// - args: An array of strings containing any additional arguments passed to the command.
//
// Return:
// - This function does not return any value.
func keepCreateTextCmdRun(cmd *cobra.Command, args []string) {
	const op = "keep_create_text"
	log := logger.GetInstance().Log.With("op", op)

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		log.Error("Error reading secret name: ",
			slog.String("error", err.Error()))
		return
	}

	data, err := cmd.Flags().GetString("data")
	if err != nil {
		log.Error("Error reading text data: ",
			slog.String("error", err.Error()))
		return
	}

	text := vaulttypes.Text{
		Data: data,
	}

	content, err := secret.EncryptSecret(text)
	if err != nil {
		log.Error("Failed to secret secret: ",
			slog.String("error", err.Error()))
		return
	}

	token, err := token_cache.GetInstance().Load()
	if err != nil {
		return
	}

	cfg := config.GetInstance().Config.Client
	keeperClient := app.NewKeeperClient(app.GetKeeperConnection(log, cfg.Address, token))

	resp, err := keeperClient.CreateItem(context.Background(), &v1.CreateItemRequestV1{
		Name:    name,
		Content: content,
	})
	if err != nil {
		log.Error("Failed to create secret: ", slog.String("error", err.Error()))
	}

	if resp == nil {
		log.Error("Nil response received from Keeper server")
		return
	}

	fmt.Printf("Secret %s version %v created successfully\n", resp.GetName(), resp.GetVersion())
}
