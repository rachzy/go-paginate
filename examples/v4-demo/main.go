package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/rachzy/go-paginate/v4/paginate"
)

// Product is a minimal model for exercising the v4 builder.
type Product struct {
	ID        int     `json:"id"         paginate:"p.id"`
	Name      string  `json:"name"       paginate:"p.name"`
	Price     float64 `json:"price"      paginate:"p.price"`
	Active    bool    `json:"active"     paginate:"p.active"`
	CreatedAt string  `json:"created_at" paginate:"p.created_at"`
}

var defaultLoggerOpts = &slog.HandlerOptions{
	Level: slog.LevelWarn,
}

func defineLoggerWithinInit(ctx context.Context) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, defaultLoggerOpts))
	paginate.InitWithLogger(ctx, &paginate.GlobalConfig{
		DefaultLimit: 10,
		MaxLimit:     100,
		DebugMode:    false,
	}, logger)
}

func defineLoggerWithSetLogger() {
	paginate.SetLogger(slog.New(slog.NewJSONHandler(os.Stdout, defaultLoggerOpts)))
}

func defineLoggerWithEnvironmentVariables(ctx context.Context) {
	// os.Setenv("GO_PAGINATE_LOG_LEVEL", "WARN")
	paginate.Init(ctx, &paginate.GlobalConfig{
		DefaultLimit: 10,
		MaxLimit:     100,
		DebugMode:    false,
	})
}

func main() {
	ctx := context.Background()
	// defer cancel()
	defineLoggerWithinInit(ctx)
	// defineLoggerWithSetLogger()
	// defineLoggerWithEnvironmentVariables(ctx)
	printModuleSource()

	params := paginate.NewPaginationParams()
	params.Page = 2
	params.Limit = 10
	params.Eq = map[string][]any{"active": {true}}
	params.Gte = map[string]any{"price": 9.99}
	params.Sort = []string{"-created_at", "name"}

	result, err := paginate.NewBuilder().
		Table("products p").
		Model(&Product{}).
		FromStruct(params).
		Build()
	if err != nil {
		log.Fatalf("build: %v", err)
	}

	fmt.Println("--- SELECT ---")
	fmt.Println(result.Query)
	fmt.Printf("args: %v\n\n", result.Args)

	fmt.Println("--- COUNT ---")
	fmt.Println(result.CountQuery)
	fmt.Printf("args: %v\n", result.CountArgs)
}

func printModuleSource() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	const mod = "github.com/rachzy/go-paginate/v4"
	for _, dep := range info.Deps {
		if dep.Path != mod {
			continue
		}
		source := dep.Version
		if dep.Replace != nil {
			source = dep.Replace.Path
			if dep.Replace.Version != "" && dep.Replace.Version != "(devel)" {
				source += "@" + dep.Replace.Version
			}
		}
		fmt.Printf("using %s (%s)\n\n", mod, source)
		return
	}
}
