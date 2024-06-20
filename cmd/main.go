package main

import (
	"context"
	app "github.com/WildEgor/sast-worker-docker/internal"
	"os/signal"
	"syscall"
	"time"
)

// @title			[TODO] Swagger Doc
// @version			1.0
// @description		[TODO]
// @termsOfService	/
// @contact.name	mail
// @contact.url		/
// @contact.email	TODO
// @license.name	MIT
// @license.url		http://www.apache.org/licenses/MIT.html
// @host			localhost:8887
// @BasePath		/
// @schemes			http
func main() {
	// Catch terminate signals
	ctx, done := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer done()

	app, _ := app.NewApp()
	app.Run(ctx)

	<-ctx.Done()

	// Wait before shutdown
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		cancel()
	}()

	app.Shutdown(ctx)
}
