package main

import (
	"context"
	"log"
	"os/signal"
	"pod-reaper/pkg/config"
	"pod-reaper/pkg/k8s"
	"pod-reaper/pkg/reaper"
	"syscall"
)

func main() {
	cfg := config.Load()

	log.Printf("pod-reaper starting: dryRun=%v pendingThreshold=%s terminatingThreshold=%s notReadyThreshold=%s crashLoopRestarts=%d excludeNamespaces=%v",
		cfg.DryRun, cfg.PendingThreshold, cfg.TerminatingThreshold, cfg.NotReadyThreshold, cfg.CrashLoopRestarts, cfg.ExcludeNamespaces)

	client, err := k8s.NewClient(cfg.KubeconfigPath)
	if err != nil {
		log.Fatalf("failed to build k8s client: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := reaper.New(client, cfg)
	if err := r.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("reaper stopped with error: %v", err)
	}
}
