//go:build !nok8s

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type kubernetesUpgradeFlags struct {
	shared.UpgradeFlags `mapstructure:",squash"`
	Helm                cmd_utils.HelmFlags
	//Ssl                 cmd_utils.SslCertFlags
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {

	upgradeCmd := &cobra.Command{
		Use:   "kubernetes",
		Short: "upgrade a local server on kubernetes",
		Long: `Upgrade a local server on kubernetes
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags kubernetesUpgradeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, upgradeKubernetes)
		},
	}

	shared.AddUpgradeFlags(upgradeCmd)
	cmd_utils.AddHelmInstallFlag(upgradeCmd)

	return upgradeCmd
}
