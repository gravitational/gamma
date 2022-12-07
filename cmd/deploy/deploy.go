package deploy

import (
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/gravitational/gamma/internal/action"
	"github.com/gravitational/gamma/internal/git"
	"github.com/gravitational/gamma/internal/logger"
	"github.com/gravitational/gamma/internal/utils"
	"github.com/gravitational/gamma/internal/workspace"
)

var outputDirectory string
var workingDirectory string
var assetPaths []string

var Command = &cobra.Command{
	Use:   "deploy",
	Short: "Builds and deploys actions",
	Long:  `Builds and deploys all the actions that have changes.`,
	Run: func(cmd *cobra.Command, args []string) {
		started := time.Now()

		if workingDirectory == "the current working directory" { // this is the default value from the flag
			wd, err := os.Getwd()
			if err != nil {
				logger.Fatalf("could not get current working directory: %v", err)
			}

			workingDirectory = wd
		}

		wd, od, err := utils.NormalizeDirectories(workingDirectory, outputDirectory)
		if err != nil {
			logger.Fatal(err)
		}

		if err := os.RemoveAll(od); err != nil {
			logger.Fatalf("could not remove output directory: %v", err)
		}

		if err := os.Mkdir(od, 0755); err != nil {
			logger.Fatalf("could not create output directory: %v", err)
		}

		repo, err := git.New(wd)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("collecting changed files")

		changed, err := repo.GetChangedFiles()
		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("files changed [%s]", strings.Join(changed, ", "))

		ws := workspace.New(wd, od)

		logger.Info("collecting actions")

		actions, err := ws.CollectActions()
		if err != nil {
			logger.Fatal(err)
		}

		if len(actions) == 0 {
			logger.Fatal("could not find any actions")
		}

		var actionNames []string
		for _, action := range actions {
			actionNames = append(actionNames, action.Name())
		}

		logger.Infof("found actions [%s]", strings.Join(actionNames, ", "))

		var actionsToBuild []action.Action

	outer:
		for _, action := range actions {
			for _, file := range changed {
				if action.Contains(file) {
					actionsToBuild = append(actionsToBuild, action)

					continue outer
				}
			}
		}

		if len(actionsToBuild) == 0 {
			logger.Warning("no actions need building, exiting")

			os.Exit(0)
		}

		var hasError bool

		for _, action := range actionsToBuild {
			logger.Infof("action %s has changes, building", action.Name())

			buildStarted := time.Now()

			if err := action.Build(); err != nil {
				hasError = true
				logger.Errorf("error building action %s: %v", action.Name(), err)

				continue
			}

			buildTook := time.Since(buildStarted)

			logger.Successf("successfully built action %s in %.2fs", action.Name(), buildTook.Seconds())

			logger.Infof("deploying action %s", action.Name())

			deployStarted := time.Now()

			if err := repo.DeployAction(action); err != nil {
				hasError = true
				logger.Errorf("error deploying action %s: %v", action.Name(), err)

				continue
			}

			deployTook := time.Since(deployStarted)

			logger.Successf("successfully deployed action %s in %.2fs", action.Name(), deployTook.Seconds())
		}

		bold := text.Colors{text.FgWhite, text.Bold}

		took := time.Since(started)

		if hasError {
			logger.Fatal(bold.Sprintf("completed with errors in %.2fs", took.Seconds()))
		}

		logger.Success(bold.Sprintf("done in %.2fs", took.Seconds()))
	},
}

func init() {
	Command.Flags().StringVarP(&outputDirectory, "output", "o", "build", "output directory")
	Command.Flags().StringVarP(&workingDirectory, "directory", "d", "the current working directory", "directory containing the monorepo of actions")
	Command.Flags().StringArrayVarP(&assetPaths, "asset", "a", []string{}, "copy over an asset to each action")
}
