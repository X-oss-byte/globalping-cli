package cmd

import (
	"fmt"

	"github.com/jsdelivr/globalping-cli/client"
	"github.com/jsdelivr/globalping-cli/model"
	"github.com/jsdelivr/globalping-cli/view"
	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:     "ping [target] from [location]",
	GroupID: "Measurements",
	Short:   "Run a ping test",
	Long: `The ping command allows sending ping requests to a target. Often used to test the network latency and stability.

Examples:
  # Ping google.com from 2 probes in New York
  ping google.com from New York --limit 2

  # Ping 1.1.1.1 from 2 probes from USA or Belgium with 10 packets in CI mode
  ping 1.1.1.1 from USA,Belgium --limit 2 --packets 10 --ci

  # Ping jsdelivr.com from a probe that is from the AWS network and is located in Montreal with latency output
  ping jsdelivr.com from aws+montreal --latency

  # Ping jsdelivr.com from a probe in ASN 123 with json output
  ping jsdelivr.com from 123 --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create context
		err := createContext(cmd.CalledAs(), args)
		if err != nil {
			return err
		}

		// Make post struct
		opts = model.PostMeasurement{
			Type:              "ping",
			Target:            ctx.Target,
			Locations:         createLocations(ctx.From),
			Limit:             ctx.Limit,
			InProgressUpdates: inProgressUpdates(ctx.CI),
			Options: &model.MeasurementOptions{
				Packets: packets,
			},
		}

		res, showHelp, err := client.PostAPI(opts)
		if err != nil {
			if showHelp {
				return err
			}
			fmt.Println(err)
			return nil
		}

		view.OutputResults(res.ID, ctx, opts)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// ping specific flags
	pingCmd.Flags().IntVar(&packets, "packets", 0, "Specifies the desired amount of ECHO_REQUEST packets to be sent (default 3)")
}
