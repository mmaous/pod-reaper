package reaper

import (
	"fmt"
	"time"

	"pod-reaper/internal/numutil"
	"pod-reaper/pkg/config"

	corev1 "k8s.io/api/core/v1"
)

// StuckReason returns non-empty reason string if pod is stuck,
// and whether deletion must be forced (grace period 0).
func StuckReason(p *corev1.Pod, cfg *config.Config, now time.Time) (reason string, force bool) {
	// Catch pods where the VM clock jumped backwards, leading to "<invalid> ago"
	if now.Sub(p.CreationTimestamp.Time) < 0 {
		return "time traveled (CreationTimestamp in future)", true
	}

	if p.DeletionTimestamp != nil {
		if now.Sub(p.DeletionTimestamp.Time) > cfg.TerminatingThreshold || now.Sub(p.DeletionTimestamp.Time) < 0 {
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
				if now.Sub(c.LastTransitionTime.Time) > cfg.NotReadyThreshold || now.Sub(c.LastTransitionTime.Time) < 0 {
					return fmt.Sprintf("Running but NotReady >%s (0/N)", cfg.NotReadyThreshold), false
				}
			}
		}
	}

	for _, cs := range p.Status.ContainerStatuses {
		if cs.State.Waiting != nil {
			if cs.State.Waiting.Reason == "CrashLoopBackOff" && cs.RestartCount >= numutil.ClampToInt32(cfg.CrashLoopRestarts) {
				return fmt.Sprintf("CrashLoopBackOff restarts=%d", cs.RestartCount), false
			}
			if cs.State.Waiting.Reason == "CreateContainerError" || cs.State.Waiting.Reason == "CreateContainerConfigError" || cs.State.Waiting.Reason == "ContainerCreating" {
				// For ContainerCreating we ensure it's actually been waiting a while
				if now.Sub(p.CreationTimestamp.Time) > cfg.PendingThreshold {
					return fmt.Sprintf("stuck in %s", cs.State.Waiting.Reason), false
				}
			}
		}
	}

	return "", false
}

// ShouldSkip checks if pod belongs to excluded namespace or is a static mirror pod.
func ShouldSkip(p *corev1.Pod, cfg *config.Config) bool {
	if p.Annotations != nil && p.Annotations["kubernetes.io/config.mirror"] != "" {
		return true // Never delete static pods
	}
	return cfg.ExcludeNamespaces[p.Namespace]
}
