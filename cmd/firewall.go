package cmd

import (
	"os"
	"strconv"

	"github.com/absam-io/absam-cli/api"
	"github.com/absam-io/absam-cli/tui"
	"github.com/absam-io/absam-cli/utils"
	"github.com/spf13/cobra"
)

const (
	PRODUCT=0
	FIREWALL_SERVICE_ID=1
	FIREWALL_ACTION = 2
)

var (
	rule_type     string
	rule_port     string
	rule_ip       string
	rule_proto    string
	rule_comment  string
	rule_expires  string
	rule_position string
)

var (
	firewallCmd = &cobra.Command{
		Use:   "firewall [server|cloud-app id [on|off|add|edit|del|status|rules]]",
		Short: "Manage a service's firewall",
		Long:  "Manage a service's firewall",

		Run: func(cmd *cobra.Command, args []string) {
			err := manageFirewall(cmd, args)
			if err != nil {
				utils.Die(err)
			}
		},
	}
)

func init() {
	firewallCmd.Flags().StringVar(&rule_type, "type", "in", "in|out")
	firewallCmd.Flags().StringVar(&rule_port, "port", "", "Default: all ports")
	firewallCmd.Flags().StringVar(&rule_ip, "ip", "0.0.0.0", "Ip")
	firewallCmd.Flags().StringVar(&rule_proto, "proto", "", "tcp|udp|icmp")
	firewallCmd.Flags().StringVar(&rule_comment, "comment", "", "Comment")
	firewallCmd.Flags().StringVar(&rule_expires, "expires", "", "Empty or 1-24")
	firewallCmd.Flags().StringVar(&rule_position, "position", "", "(default: 0)")
}

func validateFirewallArgs(cmd *cobra.Command, args []string) {
	if len(args) == EMPTY {
		cmd.Help()
		os.Exit(utils.FAIL)
	}

	if len(args) > 3 || len(args) < 3 {
		cmd.Help()
		os.Exit(utils.FAIL)
	}

	if _, err := strconv.Atoi(args[FIREWALL_SERVICE_ID]); err != nil {
		cmd.Help()
		os.Exit(utils.FAIL)
	}

	if _, err := strconv.Atoi(args[FIREWALL_ACTION]); err == nil {
		cmd.Help()
		os.Exit(utils.FAIL)
	}
}

func parseFirewallArgs(cmd *cobra.Command, args []string) {
	validateFirewallArgs(cmd, args)

	switch args[FIREWALL_ACTION] {
	case "on":
		result := api.ChangeFirewallStatus(args[FIREWALL_SERVICE_ID], args[FIREWALL_ACTION], args[PRODUCT])
		tui.PrintFirewallMessage(result.Success)
	case "off":
		result := api.ChangeFirewallStatus(args[FIREWALL_SERVICE_ID], args[FIREWALL_ACTION], args[PRODUCT])
		tui.PrintFirewallMessage(result.Success)
	case "add":
		result := api.AddRule(
			args[FIREWALL_SERVICE_ID],
			rule_type,
			rule_port,
			rule_ip,
			rule_proto,
			rule_comment,
			rule_expires,
			args[PRODUCT],
		)
		tui.PrintFirewallMessage(result.Success)
	case "edit":
		result := api.EditRule(
			args[FIREWALL_SERVICE_ID],
			rule_type,
			rule_port,
			rule_ip,
			rule_proto,
			rule_comment,
			rule_expires,
			rule_position,
			args[PRODUCT],
		)
		tui.PrintFirewallMessage(result.Success)
	case "del":
		result := api.RemoveRule(args[FIREWALL_SERVICE_ID], rule_position, args[PRODUCT])
		tui.PrintFirewallMessage(result.Success)
	case "status":
		status := api.GetFirewallStatus(args[FIREWALL_SERVICE_ID], args[PRODUCT])
		tui.PrintFirewallStatus(status)
	case "rules":
		rules := api.GetFirewallRules(args[FIREWALL_SERVICE_ID], args[PRODUCT])

		if len(rules.Rules.In) == 0 && len(rules.Rules.Out) == 0 {
			tui.PrintFirewallMessage("No rules found.")
		}
		tui.PrintAllRules(rules)
	default:
		cmd.Help()
		os.Exit(utils.FAIL)
	}
}

func manageFirewall(cmd *cobra.Command, args []string) error {
	parseFirewallArgs(cmd, args)
	return nil
}
