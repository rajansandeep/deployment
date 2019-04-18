package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/coredns/deployment/kubernetes/migration"

	"github.com/spf13/cobra"
)

// defaultCmd represents the default command
var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "default returns true if the Corefile is the default for a that version of Kubernetes. If the Kubernetes version is omitted, returns true if the Corefile is the default for any version.",
	Example: `# See if the Corefile is the default in CoreDNS v1.4.0. 
corefile-tool default --k8sversion 1.4.0 --corefile /path/to/Corefile`,
	RunE: func(cmd *cobra.Command, args []string) error {
		k8sversion, _ := cmd.Flags().GetString("k8sversion")
		corefile, _ := cmd.Flags().GetString("corefile")

		isDefault, err := defaultCorefileFromPath(k8sversion, corefile)
		if err != nil {
			return fmt.Errorf("error while checking if the Corefile is the default: %v \n", err)
		}
		fmt.Println(isDefault)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(defaultCmd)

	defaultCmd.Flags().String("k8sversion", "", "The Kuberenetes version for which you are checking the default.")
	defaultCmd.Flags().String("corefile", "", "Required: The path where your Corefile is located.")
	defaultCmd.MarkFlagRequired("corefile")
}

// defaultCorefileFromPath takes the path where the Corefile is located and checks
// if the Corefile is the default for a that version of Kubernetes.
func defaultCorefileFromPath(k8sVersion, corefilePath string) (bool, error) {
	if _, err := os.Stat(corefilePath); os.IsNotExist(err) {
		return false, err
	}

	fileBytes, err := ioutil.ReadFile(corefilePath)
	if err != nil {
		return false, err
	}
	corefileStr := string(fileBytes)
	return migration.Default(k8sVersion, corefileStr), nil
}
