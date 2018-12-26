package main

import (
	"context"
	"log"

	"github.com/alindeman/gprint"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deleteJobsAll bool

var deleteJobsCmd = &cobra.Command{
	Use:   "delete-jobs [--all] [ID...]",
	Short: "Delete jobs",
	RunE: func(_ *cobra.Command, args []string) error {
		ctx := context.Background()

		oauthConfig, err := readOAuthConfig()
		if err != nil {
			return err
		}

		token, err := readOrFetchToken(ctx, oauthConfig)
		if err != nil {
			return err
		}

		oauthClient := oauthConfig.Client(ctx, token)

		client := &gprint.Client{
			OAuthClient: oauthClient,
		}

		if deleteJobsAll {
			for {
				jobs, err := client.Jobs()
				if err != nil {
					return err
				} else if len(jobs) == 0 {
					break
				}

				for _, job := range jobs {
					log.Printf("Deleting job %s", job.ID)
					if err := client.DeleteJob(job.ID); err != nil {
						// TODO: Retry w/ backoff
						return err
					}
				}
			}
		} else {
			// TODO
			return errors.New("not implemented")
		}

		return nil
	},
}

func init() {
	deleteJobsCmd.PersistentFlags().BoolVar(&deleteJobsAll, "all", false, "Delete all jobs")
	rootCmd.AddCommand(deleteJobsCmd)
}
