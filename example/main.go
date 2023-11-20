package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/vovanec/errors"
	"github.com/vovanec/errors/loghelper"
)

type AppVersion struct {
	Major int
	Minor int
	Patch int
}

func (v AppVersion) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("major", v.Major),
		slog.Int("minor", v.Minor),
		slog.Int("patch", v.Patch),
	)
}

type Application struct {
	Name    string
	Version AppVersion
	Build   string
}

func (a Application) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", a.Name),
		slog.Any("version", a.Version),
		slog.Group("build",
			slog.String("hash", a.Build),
		),
	)
}

func doSomethingElse(ctx context.Context) error {

	slog.Info("logging in doSomethingElse",
		loghelper.Attr(ctx))

	return errors.New("error in doSomethingElse",
		slog.String("a", "a"),
	)
}

func doSomething(ctx context.Context) error {
	if err := doSomethingElse(ctx); err != nil {
		return errors.Wrap(err, "error in doSomething",
			loghelper.Attr(
				// usually one doesn't have to attach the context since caller
				// already has it, but it can be done.
				ctx,
				slog.String("b", "b"),
				slog.String("c", "c"),
			),
		)
	}
	return nil
}

func main() {

	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		),
	)

	app := Application{
		Name:  "vovan",
		Build: "20b8c3f",
		Version: AppVersion{
			Major: 1,
			Minor: 7,
			Patch: 2,
		},
	}

	// Use loghelper.Context to attach log attributes to pass them down to the callee.
	ctx := loghelper.Context(context.Background(), "application", app)

	// loghelper.Attr can be used instead of slog attribute constructors
	// if we want to extract log attributes from context or errors.
	slog.Info("application started",
		loghelper.Attr(
			ctx,
			slog.String("x", "x"),
		),
	)

	if err := doSomething(ctx); err != nil {
		slog.Error("error occurred",
			loghelper.Attr(ctx, err),
		)
	}
}
