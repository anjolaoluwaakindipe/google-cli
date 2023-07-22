/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	appclient "github.com/anjolaoluwaakindipe/testcli/internal/pkg"
	"github.com/spf13/cobra"
	"google.golang.org/api/drive/v3"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "This command gets a gopher by name",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		driveService, err := appclient.NewDriveClient()

		if err != nil {
			log.Fatalln(err)
		}
		ctx := context.Background()

		var fileList *drive.FileList

		err = driveService.Files.List().Corpora("allDrives").Q("'root' in parents").SupportsAllDrives(true).Spaces("drive").IncludeItemsFromAllDrives(true).IncludePermissionsForView("published").PageSize(100).Pages(ctx, func(fl *drive.FileList) error {
			fileList = fl
			return nil
		})

		if err != nil {
			log.Fatalf("Could not get file list: %v", err)
		}
		fmt.Println(fileList.NextPageToken)
		for _, file := range fileList.Files {
			fmt.Println(file.Name)
		}
		// var gopherName = "dr-who"

		// if len(args) >= 1 && args[0] != "" {
		// 	gopherName = args[0]
		// }

		// URL := "https://github.com/scraly/gophers/raw/main/" + gopherName + ".png"

		// fmt.Println("Try to get '" + gopherName + "' Gopher...")

		// // Get the data
		// response, err := http.Get(URL)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// defer response.Body.Close()

		// if response.StatusCode == 200 {
		// 	// Create the file
		// 	out, err := os.Create(gopherName + ".png")
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	defer out.Close()

		// 	// Writer the body to file
		// 	_, err = io.Copy(out, response.Body)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}

		// 	fmt.Println("Perfect! Just saved in " + out.Name() + "!")
		// } else {
		// 	fmt.Println("Error: " + gopherName + " not exists! :-(")
		// }

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
