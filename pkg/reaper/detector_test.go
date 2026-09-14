package reaper

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"pod-reaper/pkg/config"
)

func TestStuckReason(t *testing.T) {
	cfg := &config.Config{
		PendingThreshold:     5 * time.Minute,
		TerminatingThreshold: 2 * time.Minute,
		NotReadyThreshold:    5 * time.Minute,
		CrashLoopRestarts:    6,
	}

	now := time.Now()

	// 1. Terminating stuck past threshold -> force = true
	t.Run("TerminatingStuck", func(t *testing.T) {
		deletionTime := metav1.NewTime(now.Add(-3 * time.Minute))
		p := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				DeletionTimestamp: &deletionTime,
			},
		}
		reason, force := StuckReason(p, cfg, now)
		if reason == "" || !force {
			t.Errorf("expected stuck terminating with force=true, got reason=%q, force=%v", reason, force)
		}
	})

	// 2. Terminating within threshold -> not stuck
	t.Run("TerminatingWithinThreshold", func(t *testing.T) {
		deletionTime := metav1.NewTime(now.Add(-1 * time.Minute))
		p := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				DeletionTimestamp: &deletionTime,
			},
		}
		reason, _ := StuckReason(p, cfg, now)
		if reason != "" {
			t.Errorf("expected not stuck, got reason=%q", reason)
		}
	})

	// 3. Pending stuck past threshold -> normal delete (force = false)
	t.Run("PendingStuck", func(t *testing.T) {
		p := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				CreationTimestamp: metav1.NewTime(now.Add(-6 * time.Minute)),
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodPending,
			},
		}
		reason, force := StuckReason(p, cfg, now)
		if reason == "" || force {
			t.Errorf("expected stuck pending with force=false, got reason=%q, force=%v", reason, force)
		}
	})

	// 4. Running but NotReady past threshold -> normal delete
	t.Run("RunningNotReady", func(t *testing.T) {
		p := &corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				Conditions: []corev1.PodCondition{
					{
						Type:               corev1.PodReady,
						Status:             corev1.ConditionFalse,
						LastTransitionTime: metav1.NewTime(now.Add(-6 * time.Minute)),
					},
				},
			},
		}
		reason, force := StuckReason(p, cfg, now)
		if reason == "" || force {
			t.Errorf("expected stuck not ready with force=false, got reason=%q, force=%v", reason, force)
		}
	})

	// 5. CrashLoopBackOff with high restarts
	t.Run("CrashLoopBackOff", func(t *testing.T) {
		p := &corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				ContainerStatuses: []corev1.ContainerStatus{
					{
						RestartCount: 6,
						State: corev1.ContainerState{
							Waiting: &corev1.ContainerStateWaiting{
								Reason: "CrashLoopBackOff",
							},
						},
					},
				},
			},
		}
		reason, force := StuckReason(p, cfg, now)
		if reason == "" || force {
			t.Errorf("expected stuck crashloop with force=false, got reason=%q, force=%v", reason, force)
		}
	})

	// 6. Healthy pod -> empty reason
	t.Run("HealthyPod", func(t *testing.T) {
		p := &corev1.Pod{
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				Conditions: []corev1.PodCondition{
					{
						Type:   corev1.PodReady,
						Status: corev1.ConditionTrue,
					},
				},
			},
		}
		reason, _ := StuckReason(p, cfg, now)
		if reason != "" {
			t.Errorf("expected healthy pod not stuck, got reason=%q", reason)
		}
	})
}

func TestShouldSkip(t *testing.T) {
	cfg := &config.Config{
		ExcludeNamespaces: map[string]bool{
			"kube-system": true,
			"ingress":     true,
		},
	}

	p1 := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "kube-system"}}
	p2 := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "default"}}

	if !ShouldSkip(p1, cfg) {
		t.Errorf("expected p1 in kube-system to be skipped")
	}
	if ShouldSkip(p2, cfg) {
		t.Errorf("expected p2 in default not to be skipped")
	}
}
