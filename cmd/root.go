package cmd

import (
	"github.com/fabianonunes/syslog-perf/perf"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "syslog-perf",
	Run: func(cmd *cobra.Command, args []string) {
		number, _ := cmd.Flags().GetInt("number")
		qps, _ := cmd.Flags().GetInt("qps")
		timeout, _ := cmd.Flags().GetInt("timeout")
		workers, _ := cmd.Flags().GetInt("workers")
		host, _ := cmd.Flags().GetString("host")
		tag, _ := cmd.Flags().GetString("tag")
		messageSize, _ := cmd.Flags().GetInt("message-size")
		perf.Run(number, workers, qps, timeout, host, tag, messageSize)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntP("number", "n", 0, "Número de requisições")
	rootCmd.Flags().IntP("qps", "q", 1000, "Queries por segundo")
	rootCmd.Flags().IntP("timeout", "t", 600, "Timeout em segundos")
	rootCmd.Flags().IntP("workers", "w", 1, "Workers")
	rootCmd.Flags().StringP("host", "a", "localhost:514", "Endereço do servidor host:port")
	rootCmd.Flags().String("tag", "perf", "Tag das mensagens syslog")
	rootCmd.Flags().IntP("message-size", "s", 1024, "Tamanho da mensagem")

	_ = rootCmd.MarkFlagRequired("number")
}
