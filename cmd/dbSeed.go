/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/jerpsp/go-fiber-beginner/config"
	"github.com/jerpsp/go-fiber-beginner/db/seeds"
	"github.com/jerpsp/go-fiber-beginner/pkg/database"
	"github.com/spf13/cobra"
)

// dbSeedCmd represents the dbSeed command
var dbSeedCmd = &cobra.Command{
	Use:   "dbSeed",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.InitConfig()
		db := database.NewGormDB(cfg.PostgresDB)
		seeds.CreateSeedData(cfg, db)

		fmt.Println("RUN dbSeed Completed")
	},
}

func init() {
	rootCmd.AddCommand(dbSeedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbSeedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbSeedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
