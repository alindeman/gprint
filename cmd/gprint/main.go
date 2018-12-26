package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleCloudPrintScope = "https://www.googleapis.com/auth/cloudprint"

var credentialsFile string
var tokenFile string

var rootCmd = &cobra.Command{
	Use:   "gprint",
	Short: "Google Print CLI",
	PersistentPreRunE: func(_ *cobra.Command, args []string) error {
		if credentialsFile == "" {
			return errors.New("missing required flag: credentials-file")
		} else if tokenFile == "" {
			return errors.New("missing required flag: token-file")
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&credentialsFile, "credentials-file", "", "OAuth Client ID configuration file (downloadable from https://console.cloud.google.com/apis/credentials)")
	rootCmd.PersistentFlags().StringVar(&tokenFile, "token-file", "", "File to load or store an OAuth token")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readOAuthConfig() (*oauth2.Config, error) {
	credentialJSON, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read credentials")
	}

	cfg, err := google.ConfigFromJSON(credentialJSON, googleCloudPrintScope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse credentials")
	}

	return cfg, err
}

func readOrFetchToken(ctx context.Context, oauthConfig *oauth2.Config) (*oauth2.Token, error) {
	token, err := readTokenFromFile(tokenFile)
	if os.IsNotExist(err) {
		token, err = fetchToken(ctx, oauthConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch token")
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to load token from file")
	}

	if err := saveTokenToFile(tokenFile, token); err != nil {
		return nil, errors.Wrap(err, "failed to save token to file")
	}

	return token, err
}

func readTokenFromFile(tokenFile string) (*oauth2.Token, error) {
	token := new(oauth2.Token)

	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(token); err != nil {
		return nil, err
	}

	return token, err
}

func fetchToken(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following URL, then paste the authorization token: %v\n\n", url)
	fmt.Printf("Auth code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	return config.Exchange(ctx, authCode)
}

func saveTokenToFile(tokenFile string, token *oauth2.Token) error {
	f, err := os.OpenFile(tokenFile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}
