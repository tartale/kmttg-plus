package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/tartale/kmttg-plus/go/pkg/beacon"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/resolvers"
	"github.com/tartale/kmttg-plus/go/pkg/server"
)

const port = "8080"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kmttg",
	Short: "Port of KMTTG to golang",
	Run: func(cmd *cobra.Command, args []string) {
		startBeaconListener()
		runGraphQLServer()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kmttg.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".kmttg" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kmttg")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func startBeaconListener() {
	go beacon.Listen(context.Background())
}

func runGraphQLServer() {
	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler
	router.Use(corsHandler)

	gqlExecutableSchema := server.NewExecutableSchema(server.Config{Resolvers: &resolvers.Resolver{}})
	gqlServer := gqlhandler.NewDefaultServer(gqlExecutableSchema)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", gqlServer)

	logz.Logger.Info("connect to http://localhost:" + port + "/ for GraphQL playground")
	err := http.ListenAndServe(":"+port, router)

	logz.Logger.Fatal("error while running kmttg server", zap.Errors("error", []error{err}))
}
