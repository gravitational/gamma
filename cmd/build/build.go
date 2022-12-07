package build

import (
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	"github.com/gravitational/gamma/internal/logger"
	"github.com/gravitational/gamma/internal/utils"
	"github.com/gravitational/gamma/internal/workspace"
)

var outputDirectory string
var workingDirectory string

var Command = &cobra.Command{
	Use:   "build",
	Short: "Builds all the actions in the monorepo",
	Long:  `Builds all the actions in the monorepo and puts them into the specified output directory, separated by repo.`,
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

		var hasError bool

		for _, action := range actions {
			logger.Infof("action %s has changes, building", action.Name())

			buildStarted := time.Now()

			if err := action.Build(); err != nil {
				hasError = true
				logger.Errorf("error building action %s: %v", action.Name(), err)

				continue
			}

			buildTook := time.Since(buildStarted)

			logger.Successf("successfully built action %s in %.2fs", action.Name(), buildTook.Seconds())
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
}
