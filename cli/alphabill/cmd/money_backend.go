package cmd

import (
	"context"
	"os"
	"path/filepath"

	indexer "github.com/alphabill-org/alphabill/pkg/wallet/backend/money"
	"github.com/spf13/cobra"
)

const (
	moneyBackendHomeDir = "money-backend"

	serverAddrCmdName  = "server-addr"
	dbFileCmdName      = "db"
	listBillsPageLimit = "list-bills-page-limit"
)

type moneyBackendConfig struct {
	Base               *baseConfiguration
	AlphabillUrl       string
	ServerAddr         string
	DbFile             string
	LogLevel           string
	LogFile            string
	ListBillsPageLimit int
}

func (c *moneyBackendConfig) GetDbFile() (string, error) {
	var dbFile string
	if c.DbFile != "" {
		dbFile = c.DbFile
	} else {
		dbFile = filepath.Join(c.Base.HomeDir, moneyBackendHomeDir, indexer.BoltBillStoreFileName)
	}
	err := os.MkdirAll(filepath.Dir(dbFile), 0700) // -rwx------
	if err != nil {
		return "", err
	}
	return dbFile, nil
}

// newMoneyBackendCmd creates a new cobra command for the money-backend component.
func newMoneyBackendCmd(baseConfig *baseConfiguration) *cobra.Command {
	config := &moneyBackendConfig{Base: baseConfig}
	var walletCmd = &cobra.Command{
		Use:   "money-backend",
		Short: "starts money backend service",
		Long:  "starts money backend service, indexes all transactions by owner predicates, starts http server",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// initialize config so that baseConfig.HomeDir gets configured
			err := initializeConfig(cmd, baseConfig)
			if err != nil {
				return err
			}
			// init logger
			return initWalletLogger(&walletConfig{LogLevel: config.LogLevel, LogFile: config.LogFile})
		},
	}
	walletCmd.PersistentFlags().StringVar(&config.LogFile, logFileCmdName, "", "log file path (default output to stderr)")
	walletCmd.PersistentFlags().StringVar(&config.LogLevel, logLevelCmdName, "INFO", "logging level (DEBUG, INFO, NOTICE, WARNING, ERROR)")
	walletCmd.AddCommand(startMoneyBackendCmd(config))
	return walletCmd
}

func startMoneyBackendCmd(config *moneyBackendConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use: "start",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execMoneyBackendStartCmd(cmd.Context(), config)
		},
	}
	cmd.Flags().StringVarP(&config.AlphabillUrl, alphabillNodeURLCmdName, "u", defaultAlphabillNodeURL, "alphabill node url")
	cmd.Flags().StringVarP(&config.ServerAddr, serverAddrCmdName, "s", "localhost:9654", "server address")
	cmd.Flags().StringVarP(&config.DbFile, dbFileCmdName, "f", "", "path to the database file (default: $AB_HOME/"+moneyBackendHomeDir+"/"+indexer.BoltBillStoreFileName+")")
	cmd.Flags().IntVarP(&config.ListBillsPageLimit, listBillsPageLimit, "l", 100, "GET /list-bills request default/max limit size")
	return cmd
}

func execMoneyBackendStartCmd(ctx context.Context, config *moneyBackendConfig) error {
	dbFile, err := config.GetDbFile()
	if err != nil {
		return err
	}
	return indexer.CreateAndRun(ctx, &indexer.Config{
		ABMoneySystemIdentifier: defaultABMoneySystemIdentifier,
		AlphabillUrl:            config.AlphabillUrl,
		ServerAddr:              config.ServerAddr,
		DbFile:                  dbFile,
		ListBillsPageLimit:      config.ListBillsPageLimit,
	})
}