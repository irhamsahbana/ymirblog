// Package cmd is the command surface of ymirblog cli tool provided by kubuskotak.
// # This manifest was generated by ymir. DO NOT EDIT.
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"entgo.io/ent/dialect"
"database/sql"	

    	"github.com/fatih/color"
    	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/persist/ymirblog"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/adapters"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/infrastructure"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/ports/rest"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/shared"
	"gitlab.playcourt.id/dedenurr12/ymirblog/pkg/version"
)

type rootOptions struct {
	Path     string
	Filename string
}

// NewRootCmd creates the root command.
func NewRootCmd() *cobra.Command {
	root := &rootOptions{}
	cmds := &cobra.Command{
		Use:   "ymirblog",
		Short: "lorem ipsum abracadabra, make some a noise.",
		Long: GenerateTemplate(description,
			map[string]any{
				"BuildTime":  version.GetVersion().BuildDate,
				"Version":    version.GetVersion().VersionNumber(),
				"CommitHash": version.GetVersion().Revision,
			}),
		RunE: root.runServer,
	}
	cmds.PersistentFlags().StringVarP(
		&root.Filename, "config", "c", "config.yaml", "config file name")
	cmds.PersistentFlags().StringVarP(
		&root.Path, "config-path", "d", "./", "config dir path")

	// subcommands
	cmds.AddCommand(newVersionCmd(), newMigrateCmd())

	return cmds
}

func (r *rootOptions) runServer(_ *cobra.Command, _ []string) error {
	infrastructure.Configuration(
		infrastructure.WithPath(r.Path),
		infrastructure.WithFilename(r.Filename),
	).Initialize()
	infrastructure.InitializeLogger()
	info := color.New(color.BgBlack, color.FgRed).SprintFunc()
	fmt.Printf("%s\n", info(fmt.Sprintf(logo,
	    version.GetVersion().VersionNumber(),
		infrastructure.Envs.Ports.HTTP,
	)))
	log.Info().Str("Stage", infrastructure.Envs.App.Environment).Msg("server running...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// open-telemetry
	var (
		cleanupTracer infrastructure.TracerReturnFunc
		cleanupMetric infrastructure.MetricReturnFunc
	)
	if infrastructure.Envs.Telemetry.CollectorEnable {
		traceExp, traceErr := infrastructure.TraceExporter(ctx,
			infrastructure.Envs.Telemetry.CollectorDebug,
			infrastructure.Envs.Telemetry.CollectorGrpcAddr,
			infrastructure.Envs.Server.Timeout)
		if traceErr != nil {
			log.Error().Err(traceErr).Msg("exporter trace is failed")
			return traceErr
		} // tracing exporter
		cleanupTracer = infrastructure.InitTracer(traceExp)
		// initial service for tracing
		var span trace.Span
		tp := otel.Tracer(infrastructure.Envs.App.ServiceName)
		ctx, span = tp.Start(ctx, "root-running")
		defer span.End()

		metricExp, metricErr := infrastructure.MetricExporter(ctx,
			infrastructure.Envs.Telemetry.CollectorDebug,
			infrastructure.Envs.Telemetry.CollectorGrpcAddr,
			infrastructure.Envs.Server.Timeout)
		if metricErr != nil {
			log.Error().Err(metricErr).Msg("exporter metric is failed")
			return metricErr
		} // tracing exporter
		cleanupMetric = infrastructure.InitMetric(metricExp)
	}
	/**
	* Initialize Main
	*/
d := infrastructure.Envs.YmirBlogMySQL
    adaptor := &adapters.Adapter{}
	adaptor.Sync(
		adapters.WithYmirBlogMySQL(&adapters.YmirBlogMySQL{
			NetworkDB: adapters.NetworkDB{
				Database: d.Database,
				User:d.User,
				Password: d.Password,
				Host: d.Host,
				Port: d.Port,
			},
		}),
) // adapters init
_ = ymirblog.Driver(
		ymirblog.WithDriver(adaptor.YmirBlogMySQL, dialect.MySQL),
		ymirblog.WithTxIsolationLevel(sql.LevelSerializable),
	)

    var errCh chan error
	/**
	* Initialize HTTP
	*/
	h := rest.NewServer(
		rest.WithPort(strconv.Itoa(infrastructure.Envs.Ports.HTTP)),
	)
	h.Handler(rest.Routes().Register(
		func(c chi.Router) http.Handler {
		    // http register handler
			return c
		},
	))
	if err := h.ListenAndServe(); err != nil {
		return err
	}
	errCh = h.Error()
	// end http
	stopCh := shared.SetupSignalHandler()
	return shared.Graceful(stopCh, errCh, func(ctx context.Context) { // graceful shutdown
		log.Info().Dur("timeout", infrastructure.Envs.Server.Timeout).Msg("Shutting down HTTP/HTTPS server")
		// open-telemetry
		if infrastructure.Envs.Telemetry.CollectorEnable {
            if err := cleanupTracer(context.Background()); err != nil {
                log.Error().Err(err).Msg("tracer provider server is failed shutdown")
            }
            if err := cleanupMetric(context.Background()); err != nil {
                log.Error().Err(err).Msg("metric provider server is failed shutdown")
            }
		}
		// rest
		if err := h.Quite(context.Background()); err != nil {
			log.Error().Err(err).Msg("http server is failed shutdown")
		}
		h.Stop()
		// adapters
		if err := adaptor.UnSync(); err != nil {
			log.Error().Err(err).Msg("there is failed on UnSync adapter")
		}
	})
}

// Execute is the execute command for root command.
func Execute() error {
	return NewRootCmd().Execute()
}
