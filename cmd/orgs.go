package cmd

import "github.com/spf13/cobra"

var orgsCmd = &cobra.Command{
	Use: "orgs",
	Short: "Manage Organizations",
	Long:  `Commands for listing and viewing organizations.`,
}

func init(){
	rootCmd.AddCommand(orgsCmd)
}