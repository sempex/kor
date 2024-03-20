package kor

import (
	"github.com/spf13/cobra"
	"github.com/yonahd/kor/pkg/kor"
)

var clusterRoleBindingCmd = &cobra.Command{
	Use:     "clusterrolebindings",
	Aliases: []string{"clusterrolebindings"},
	Short:   "Gets unused clusterrolebindings",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		clientset := kor.GetKubeClient(kubeconfig)

		kor.RetrieveUsedClusterRoleBindings(*clientset, filterOptions)
	},
}

func init() {
	rootCmd.AddCommand(clusterRoleBindingCmd)
}
