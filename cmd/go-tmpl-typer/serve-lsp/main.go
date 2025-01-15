package serve_lsp

import (
	"context"
	"os"

	"github.com/creachadair/jrpc2"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/walteh/go-tmpl-typer/pkg/lsp"
	"gitlab.com/tozd/go/errors"
)

type Handler struct {
	debug bool
}

func NewServeLSPCommand() *cobra.Command {
	me := &Handler{}

	cmd := &cobra.Command{
		Use:   "serve-lsp",
		Short: "start the language server",
	}

	cmd.Flags().BoolVar(&me.debug, "debug", false, "enable debug logging")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return me.Run(cmd.Context())
	}

	return cmd
}

func (me *Handler) Run(ctx context.Context) error {
	// Create a new LSP server with all the components it needs
	server := lsp.NewServer(ctx)

	opts := &jrpc2.ServerOptions{
		Logger: func(str string) {
			zerolog.Ctx(ctx).Info().Msg(str)
		},
	}

	// Start the server using stdin/stdout
	if err := server.Run(ctx, os.Stdin, os.Stdout, opts); err != nil {
		return errors.Errorf("failed to start language server: %w", err)
	}

	return nil
}

// func main() {
// 	debug := false
// 	flag.BoolVar(&debug, "debug", false, "enable debug logging")
// 	flag.Parse()

// 	logger := zerolog.New(os.Stderr).With().
// 		Str("component", "lsp-server").
// 		Bool("debug", debug).
// 		Timestamp().
// 		Logger()
// 	ctx := logger.WithContext(context.Background())

// 	if debug {
// 		zerolog.Ctx(ctx).Info().Msg("starting language server with debug logging enabled")
// 	}

// 	server := lsp.NewServer(
// 		parser.NewDefaultTemplateParser(),
// 		types.NewDefaultValidator(),
// 		ast.NewDefaultPackageAnalyzer(),
// 		diagnostic.NewDefaultGenerator(),
// 		debug,
// 	)

// 	if err := server.Start(ctx, os.Stdin, os.Stdout); err != nil {
// 		zerolog.Ctx(ctx).Error().Err(err).Msg("server error")
// 		os.Exit(1)
// 	}
// }
