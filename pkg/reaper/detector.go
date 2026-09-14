package reaper

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"pod-reaper/pkg/config"
)

// StuckReason returns non-empty reason string if pod is stuck,
// and whether deletion must be forced (grace period 0).
func StuckReason(p *corev1.Pod, cfg *config.Config, now time.Time) (reason string, force bool) {
	if p.DeletionTimestamp != nil {
		if now.Sub(p.DeletionTimestamp.Time) > cfg.TerminatingThreshold {
			return fmt.Sprintf("stuck Terminating >%s", cfg.TerminatingThreshold), true
		}
		return "", false
	}

	switch p.Status.Phase {
	case corev1.PodPending:
		if now.Sub(p.CreationTimestamp.Time) > cfg.PendingThreshold {
			return fmt.Sprintf("stuck Pending >%s (likely CNI/scheduling)", cfg.PendingThreshold), false
		}
	case corev1.PodRunning:
		for _, c := range p.Status.Conditions {
			if c.Type == corev1.PodReady && c.Status != corev1.ConditionTrue {
				if now.Sub(c.LastTransitionTime.Time) > cfg.NotReadyThreshold {
					return fmt.Sprintf("Running but NotReady >%s (0/N)", cfg.NotReadyThreshold), false
				}
			}
		}
	}

	for _, cs := range p.Status.ContainerStatuses {
		if cs.State.Waiting != nil && cs.State.Waiting.Reason == "CrashLoopBackOff" {
			if cs.RestartCount >= int32(cfg.CrashLoopRestarts) {
				return fmt.Sprintf("CrashLoopBackOff restarts=%d", cs.RestartCount), false
			}
		}
	}

	return "", false
}

// ShouldSkip checks if pod belongs to excluded namespace.
func ShouldSkip(p *corev1.Pod, cfg *config.Config) bool {
	return cfg.ExcludeNamespaces[p.Namespace]
}
