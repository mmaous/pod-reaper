package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	os.Clearenv()

	cfg := Load()
	if !cfg.DryRun {
		t.Errorf("expected DryRun=true, got %v", cfg.DryRun)
	}
	if cfg.PendingThreshold != 5*time.Minute {
		t.Errorf("expected PendingThreshold=5m, got %v", cfg.PendingThreshold)
	}
	if cfg.TerminatingThreshold != 2*time.Minute {
		t.Errorf("expected TerminatingThreshold=2m, got %v", cfg.TerminatingThreshold)
	}
	if cfg.NotReadyThreshold != 5*time.Minute {
		t.Errorf("expected NotReadyThreshold=5m, got %v", cfg.NotReadyThreshold)
	}
	if cfg.CrashLoopRestarts != 6 {
		t.Errorf("expected CrashLoopRestarts=6, got %d", cfg.CrashLoopRestarts)
	}
	if len(cfg.ExcludeNamespaces) != 0 {
		t.Errorf("expected no excluded namespaces by default")
	}
}

func TestLoadCustomEnv(t *testing.T) {
	os.Clearenv()
	t.Setenv("DRY_RUN", "false")
	t.Setenv("CRASHLOOP_RESTART_THRESHOLD", "3")
	t.Setenv("EXCLUDE_NAMESPACES", "kube-system,default,monitoring")

	cfg := Load()
	if cfg.DryRun {
		t.Errorf("expected DryRun=false, got true")
	}
	if cfg.CrashLoopRestarts != 3 {
		t.Errorf("expected CrashLoopRestarts=3, got %d", cfg.CrashLoopRestarts)
	}
	if !cfg.ExcludeNamespaces["default"] || !cfg.ExcludeNamespaces["monitoring"] {
		t.Errorf("expected default and monitoring excluded")
	}
}
