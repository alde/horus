package cmd

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configure the horus CLI",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		cfg := config.Config{}
		fmt.Printf("Configuring the Horus CLI\n\n")
		fmt.Println("Horus settings")
		fmt.Println("--------------")
		fmt.Print("Host (ex: https://horus.local:7654/): ")
		var host string
		var err error
		if host, err = reader.ReadString('\n'); err != nil {
			log.Fatal("error reading host")
		}
		cfg.Horus.Host = strings.Trim(host, "\n")
		var project, keyname, keyring, location string
		fmt.Println("GoogleCloudKMS settings")
		fmt.Println("-----------------------")
		fmt.Print("GCP Project ID: ")
		if project, err = reader.ReadString('\n'); err != nil {
			log.Fatal("error reading GCP Project ID")
		}
		fmt.Print("KMS keyname: ")
		if keyname, err = reader.ReadString('\n'); err != nil {
			log.Fatal("error reading KMS keyname")
		}
		fmt.Print("KMS keyring: ")
		if keyring, err = reader.ReadString('\n'); err != nil {
			log.Fatal("error reading KMS keyring")
		}
		fmt.Print("KMS key location: ")
		if location, err = reader.ReadString('\n'); err != nil {
			log.Fatal("error reading KMS key location")
		}
		cfg.GoogleKMS = struct {
			Project  string
			KeyName  string
			KeyRing  string
			Location string
		}{
			strings.Trim(project, "\n"),
			strings.Trim(keyname, "\n"),
			strings.Trim(keyring, "\n"),
			strings.Trim(location, "\n"),
		}
		cfg.WriteFile(configFile)
	},
}

// Execute the rootCmd
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "configuration file")

	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(configureCmd)

	rootCmd.Execute()
}
