package main

import (
	"github.com/gsousadev/doolar2/internal/application"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Executa varredura de rede",
	Run: func(cmd *cobra.Command, args []string) {
		application.RunScanLoop()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
