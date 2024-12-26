package application

import (
	"context"
	"fmt"
)

type Logger interface {
	Info(ctx context.Context, err string)
	Warning(ctx context.Context, err string)
	Error(ctx context.Context, err string)
	Critical(ctx context.Context, err string)
}

type PrintLnLogger struct {
}

func (l PrintLnLogger) Info(ctx context.Context, err string) {
	fmt.Println("info: " + err)
}

func (l PrintLnLogger) Warning(ctx context.Context, err string) {
	fmt.Println("warning: " + err)
}

func (l PrintLnLogger) Error(ctx context.Context, err string) {
	fmt.Println("error: " + err)
}

func (l PrintLnLogger) Critical(ctx context.Context, err string) {
	fmt.Println("critical: " + err)
}
