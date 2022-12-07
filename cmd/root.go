package cmd

import (
	"github.com/spf13/cobra"

	"github.com/gravitational/gamma/cmd/build"
	"github.com/gravitational/gamma/cmd/deploy"
	"github.com/gravitational/gamma/internal/color"
)

var rootCmd = &cobra.Command{
	Use:   "gamma",
	Short: "Gamma builds a monorepo of Github actions into individual repos",
}

var gammaLogo = "\x1B[38;2;236;147;168m#\x1B[39m\x1B[38;2;226;142;179m#\x1B[39m\x1B[38;2;216;138;191m#\x1B[39m \x1B[38;2;206;133;202mG\x1B[39m\x1B[38;2;197;129;214ma\x1B[39m\x1B[38;2;187;124;225mm\x1B[39m\x1B[38;2;177;120;237mm\x1B[39m\x1B[38;2;167;115;248ma\x1B[39m \x1B[38;2;167;115;248mb\x1B[39m\x1B[38;2;160;109;244my\x1B[39m \x1B[38;2;153;104;240mT\x1B[39m\x1B[38;2;146;98;236me\x1B[39m\x1B[38;2;138;93;232ml\x1B[39m\x1B[38;2;131;87;228me\x1B[39m\x1B[38;2;124;82;225mp\x1B[39m\x1B[38;2;117;76;221mo\x1B[39m\x1B[38;2;110;70;217mr\x1B[39m\x1B[38;2;103;65;213mt\x1B[39m \x1B[38;2;95;59;209m#\x1B[39m\x1B[38;2;88;54;205m#\x1B[39m\x1B[38;2;81;48;201m#\x1B[39m"

func Execute() error {
	return rootCmd.Execute()
}
func logo() string {
	return gammaLogo
}

func init() {
	cobra.AddTemplateFunc("emoji", emoji)
	cobra.AddTemplateFunc("colorize", colorize)
	cobra.AddTemplateFunc("green", color.Green)
	cobra.AddTemplateFunc("logo", logo)

	rootCmd.AddCommand(deploy.Command)
	rootCmd.AddCommand(deploy.Command)

	rootCmd.SetHelpTemplate(`{{ logo }}

{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)

	rootCmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{green .CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{emoji .Name}}  {{colorize .Name (rpad .Name .NamePadding) }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)
}

func colorize(s, name string) string {
	switch s {
	case build.Command.Name():
		return color.Magenta(name)
	case deploy.Command.Name():
		return color.Teal(name)
	case "help":
		return color.Purple(name)
	case "completion":
		return color.Yellow(name)
	}

	return s
}

func emoji(s string) string {
	switch s {
	case build.Command.Name():
		return "🔧"
	case deploy.Command.Name():
		return "🚀"
	case "help":
		return "❓"
	case "completion":
		return "⚡️"
	}

	return s
}
