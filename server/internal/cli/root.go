/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/ajugalushkin/goph-keeper/internal/app/client/auth"
)

const version = "1.0"

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(root *cobra.Command) {
	//err := rootCmd.Execute()
	//if err != nil {
	//	os.Exit(1)
	//}
	cobra.CheckErr(root.Execute())
}

func NewRootCmd() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use: `goph-keeper [-s SERVER_ADDRESS] [-t TIMEOUT] [-r RETRIES_COUNT]
	CONFIG: Config file in yaml format`,
		Short:   "GophKeeper cli client",
		Long:    `GophKeeper cli client allows keep and return secrets in/from server.`,
		Version: version,
	}

	//home, err := os.UserHomeDir()
	//cobra.CheckErr(err)
	var (
		srvAddr      string        // for persistent flag
		timeout      time.Duration // for persistent flag
		retriesCount int           // for persistent flag
	)

	rootCmd.PersistentFlags().StringVarP(&srvAddr, "server", "s", "localhost:5000", "ip/dns:port")
	rootCmd.PersistentFlags().DurationP("timeout", "t", timeout, "timeout for client connections")
	rootCmd.PersistentFlags().IntP("retries_count", "r", retriesCount, "number of times to retry client connections")
	//rootCmd.PersistentFlags().StringVarP(&tcache, "cache", "c", home+"/.keeppas.token", "token cache")
	//cobra.OnInitialize(func() {
	//	client.config.LogLevel = config.LoggerConfig(dbg)
	//	client.config.ServerAddr = srvAddr
	//	client.config.TokenCache = tcache
	//	client.logger = loggerConfig(client.config.LogLevel)
	//	client.transport = newGRPCConnection
	//})
	//
	//rootCmd.SetVersionTemplate(version + " Build at " + BuildTime + "\n")
	//
	//kvCmd := newKVCmd()
	//kvCmd.AddCommand(newKVCmdAdd(&client))
	//kvCmd.AddCommand(newKVCmdCP(&client))
	//kvCmd.AddCommand(newKVCmdGet(&client))
	//kvCmd.AddCommand(newKVCmdRm(&client))
	//kvCmd.AddCommand(newKVCmdRename(&client))
	//kvCmd.AddCommand(newKVCmdUpdate(&client))
	//kvCmd.AddCommand(newKVCmdList(&client))
	//

	client, err := auth.New(context.Background(), srvAddr, timeout, retriesCount)
	if err != nil {
		return nil
	}

	rootCmd.AddCommand(newRegisterCmd(client))
	rootCmd.AddCommand(newLoginCmd(client))
	//rootCmd.AddCommand(kvCmd)

	return rootCmd
}
