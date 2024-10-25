package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/term"

	"github.com/tartale/kmttg-plus/go/dist"
	"github.com/tartale/kmttg-plus/go/pkg/beacon"
	"github.com/tartale/kmttg-plus/go/pkg/client"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/jobs"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
	"github.com/tartale/kmttg-plus/go/pkg/resolvers"
	"github.com/tartale/kmttg-plus/go/pkg/server"
	"github.com/tartale/kmttg-plus/go/pkg/tivos"
)

const crlf = "\r\n"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kmttg",
	Short: "Port of KMTTG to golang",
	Run: func(cmd *cobra.Command, args []string) {
		startBeaconListener(cmd.Context())
		startLoader(cmd.Context())
		startJobWorkers(cmd.Context())
		runWebServer(cmd.Context())
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
	cobra.OnInitialize(func() { config.InitConfig(cfgFile) })
	config.Values.Validate()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kmttg.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "terminal",
		Short: "Starts a terminal for sending RPC commands to a Tivo",
		Run: func(cmd *cobra.Command, args []string) {
			startBeaconListener(cmd.Context())
			runTerminal()
		},
	})
}

func startBeaconListener(ctx context.Context) {
	go beacon.Listen(ctx)
}

func startLoader(ctx context.Context) {
	go tivos.RunBackgroundLoader(ctx)
}

func startJobWorkers(ctx context.Context) {
	go jobs.RunWorkerPool(ctx)
}

func runWebServer(_ context.Context) {
	router := mux.NewRouter()

	addCORSMiddleware(router)
	addGraphQLRoutes(router)
	addWebUIRoutes(router)

	err := http.ListenAndServe(":"+config.Values.Port, router)
	logz.Logger.Fatal("error while running kmttg server", zap.Errors("error", []error{err}))
}

func addCORSMiddleware(router *mux.Router) {
	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	corz := cors.New(cors.Options{
		AllowCredentials: true,
		Debug:            true,
	})
	corz.Log = logz.LoggerX
	router.Use(corz.Handler)
}

func addGraphQLRoutes(router *mux.Router) {
	gqlExecutableSchema := server.NewExecutableSchema(server.Config{Resolvers: &resolvers.Resolver{}})
	gqlServer := gqlhandler.NewDefaultServer(gqlExecutableSchema)

	router.Handle("/api/playground", playground.Handler("GraphQL playground", "/api/query"))
	router.Handle("/api/query", gqlServer)

	logz.Logger.Info("POST to http://localhost:" + config.Values.Port + "/api/query for GraphQL queries")
	logz.Logger.Info("connect to http://localhost:" + config.Values.Port + "/api/playground for GraphQL playground")
}

func addWebUIRoutes(router *mux.Router) {
	var webUIServer http.Handler
	if config.Values.WebUIDir != "" {
		webUIServer = http.FileServer(http.Dir(config.Values.WebUIDir))
	} else {
		webUIFiles, err := fs.Sub(dist.Filesystem, "webui")
		if err != nil {
			panic(err)
		}
		webUIServer = http.FileServer(http.FS(webUIFiles))
	}

	router.PathPrefix("/").Handler(http.StripPrefix("/", webUIServer))

	logz.Logger.Info("connect to http://localhost:" + config.Values.Port + " for the KMTTG web UI")
}

func runTerminal() {
	logz.Logger = logz.NopLogger.Logger

	fmt.Println("detecting Tivos on the network")
	for {
		if len(tivos.List(context.Background())) > 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// TODO: allow selection of a Tivo
	tvo := tivos.List(context.Background())[0]
	tivoClient, err := client.NewRpcClient(tvo)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	defer tivoClient.Close()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	terminal := term.NewTerminal(os.Stdin, tvo.Name+")> ")

	go io.Copy(os.Stdout, tivoClient)

	for {
		line, err := terminal.ReadLine()
		if err != nil {
			fmt.Printf("error: %v", err)
			break
		}
		if strings.ToLower(line) == "exit" {
			break
		}
		tivoClient.Write([]byte(line + crlf))
	}
}
