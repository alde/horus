package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/alde/horus/cli/config"
	"github.com/alde/horus/cli/decryptor"
	"github.com/alde/horus/cli/downloader"
)

var rootCmd = &cobra.Command{
	Use:   "horus",
	Short: "horus cli is used to interact with horus-backend",
}

var configFile string

var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "download and decrypt a secret",
	Example: `horus download github.com/alde/horus DOCKER_LOGIN`,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New(configFile)
		dl := downloader.New(&http.Client{}, cfg)
		secret := dl.Download(args[0], args[1])

		s, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			log.Fatal("failed decoding secret from base64")
		}
		dec, err := decryptor.NewGoogleCloudKMS(context.Background(), cfg, nil)
		if err != nil {
			log.Fatal("failed creating decryptor")
		}

		plaintext, err := dec.Decrypt(s)
		if err != nil {
			log.Fatal("failed decrypting secret")
		}

		fmt.Println(plaintext)
	},
}

// Execute the rootCmd
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "configuration file")

	rootCmd.AddCommand(downloadCmd)

	rootCmd.Execute()
}
