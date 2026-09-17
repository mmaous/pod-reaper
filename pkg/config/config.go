package config

import (
	"log"
	"os"
	"pod-reaper/internal/strutil"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DryRun               bool
	PendingThreshold     time.Duration
	TerminatingThreshold time.Duration
	NotReadyThreshold    time.Duration
	CrashLoopRestarts    int
	ExcludeNamespaces    map[string]bool
	ExcludeLabelSelector string
	KubeconfigPath       string
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("bad duration %s=%q, using default %s", key, strutil.SanitizeForLog(v), def) //nolint:gosec // G706: v is sanitized via strutil.SanitizeForLog above
		return def
	}
	return d
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("bad int %s=%q, using default %d", key, strutil.SanitizeForLog(v), def) //nolint:gosec // G706: v is sanitized via strutil.SanitizeForLog above
		return def
	}
	return i
}

func getEnvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Printf("bad bool %s=%q, using default %v", key, strutil.SanitizeForLog(v), def) //nolint:gosec // G706: v is sanitized via strutil.SanitizeForLog above
		return def
	}
	return b
}

func Load() *Config {
	excl := getEnv("EXCLUDE_NAMESPACES", "")
	m := map[string]bool{}
	for _, ns := range strings.Split(excl, ",") {
		ns = strings.TrimSpace(ns)
		if ns != "" {
			m[ns] = true
		}
	}
	return &Config{
		DryRun:               getEnvBool("DRY_RUN", true),
		PendingThreshold:     getEnvDuration("PENDING_THRESHOLD", 2*time.Minute),
		TerminatingThreshold: getEnvDuration("TERMINATING_THRESHOLD", 2*time.Minute),
		NotReadyThreshold:    getEnvDuration("NOT_READY_THRESHOLD", 2*time.Minute),
		CrashLoopRestarts:    getEnvInt("CRASHLOOP_RESTART_THRESHOLD", 3),
		ExcludeNamespaces:    m,
		ExcludeLabelSelector: getEnv("EXCLUDE_LABEL_SELECTOR", ""),
		KubeconfigPath:       getEnv("KUBECONFIG", ""),
	}
}
