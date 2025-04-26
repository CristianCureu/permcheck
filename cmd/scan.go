package cmd

import (
	"runtime"

	"github.com/spf13/cobra"

	"github.com/cristiancureu/permcheck/internal"
)

func NewScanCmd(cfg *internal.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan a directory for insecure file permissions",
		Run: func(cmd *cobra.Command, args []string) {
			insecureTag := internal.DefaultInsecureTag
			secureTag := internal.DefaultSecureTag

			if cfg.IsTerminal {
				insecureTag = internal.ColorInsecureTag
				secureTag = internal.ColorSecureTag
			}

			cfg.InsecureTag = insecureTag
			cfg.SecureTag = secureTag

			if err := internal.RunScan(args, cfg); err != nil {
				cmd.PrintErrf("Scan failed: %v\n", err)
			}
		},
	}

	cmd.Flags().IntVarP(&cfg.NumWorkers, "workers", "w", runtime.NumCPU()*2, "Number of concurrent workers (default: CPU cores x 2)")
	cmd.Flags().BoolVar(&cfg.InsecureOnly, "insecure-only", false, "Show only insecure files")
	cmd.Flags().BoolVar(&cfg.FixMode, "fix", false, "Fix permissions of insecure files")

	return cmd
}
