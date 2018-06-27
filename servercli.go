package celerity

import (
	"fmt"
	"os"

	"github.com/5Sigma/vox"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// HandleCLI - Use the built in CLI handling for the server.
func HandleCLI(server *Server) error {

	var cfgFile string

	var rootCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the server",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			banner := `
	███████╗███████╗██╗     ███████╗██████╗ ██╗████████╗██╗   ██╗
	██╔════╝██╔════╝██║     ██╔════╝██╔══██╗██║╚══██╔══╝╚██╗ ██╔╝
	██║     █████╗  ██║     █████╗  ██████╔╝██║   ██║    ╚████╔╝ 
	██║     ██╔══╝  ██║     ██╔══╝  ██╔══██╗██║   ██║     ╚██╔╝  
 	███████╗███████╗███████╗███████╗██║  ██║██║   ██║      ██║   
 	╚═════╝╚══════╝╚══════╝╚══════╝╚══╝  ╚═╝╚═╝   ╚═╝      ╚═╝   
					Celerity v1.0
			`
			hostString := fmt.Sprintf("%s:%s",
				viper.GetString("host"),
				viper.GetString("port"),
			)
			vox.Println(banner)
			vox.Println("Listening on ", hostString)
			server.Start(hostString)
		},
	}

	runCmd.PersistentFlags().Int("port", 5000, "Web server listening port")
	viper.BindPFlag("port", runCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("port", "5000")

	runCmd.PersistentFlags().String("host", "0.0.0.0", "Web server listening port")
	viper.BindPFlag("host", runCmd.PersistentFlags().Lookup("host"))
	viper.SetDefault("host", "0.0.0.0")

	runCmd.PersistentFlags().String("env", "debug", "Select the environment setup. Can be 'dev' or 'prod'")
	viper.BindPFlag("env", runCmd.PersistentFlags().Lookup("env"))
	viper.SetDefault("env", "dev")

	rootCmd.AddCommand(runCmd)

	var routesCmd = &cobra.Command{
		Use:   "routes",
		Short: "List routes",
		Long:  `Prints a list of all registered routes in the application.`,
		Run: func(cmd *cobra.Command, args []string) {
			vox.Println("All server routes:\n")
			printScope(server.Router.Root)
			vox.Println("\n")
		},
	}
	rootCmd.AddCommand(routesCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")

	// SET UP FLAGS AND DEFAULTS

	cobra.OnInitialize(func() {
		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			viper.AddConfigPath(".")
			viper.SetConfigName(".config")
		}

		viper.AutomaticEnv() // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}

	})

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

func printScope(s *Scope) {
	vox.Println(vox.Yellow, "[SCOPE]", vox.ResetColor, " ", s.Path)

	for _, ss := range s.Scopes {
		printScope(ss)
	}

	for _, r := range s.Routes {
		vox.Println("\t", r.Method, "\t", r.Path)
	}
}
