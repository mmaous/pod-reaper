package reaper

import (
	"context"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"pod-reaper/pkg/config"
)

type Reaper struct {
	client kubernetes.Interface
	cfg    *config.Config
}

func New(client kubernetes.Interface, cfg *config.Config) *Reaper {
	return &Reaper{
		client: client,
		cfg:    cfg,
	}
}

func (r *Reaper) Run(ctx context.Context) error {
	pods, err := r.client.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		LabelSelector: negateSelector(r.cfg.ExcludeLabelSelector),
	})
	if err != nil {
		log.Printf("ERROR listing pods: %v", err)
		return err
	}

	now := time.Now()
	checked := 0
	for i := range pods.Items {
		p := &pods.Items[i]
		if ShouldSkip(p, r.cfg) {
			continue
		}
		checked++
		reason, force := StuckReason(p, r.cfg, now)
		if reason == "" {
			continue
		}
		log.Printf("STUCK pod %s/%s: %s", p.Namespace, p.Name, reason)
		r.DeletePod(ctx, p, force)
	}
	log.Printf("check complete: %d pods evaluated", checked)
	return nil
}

func (r *Reaper) DeletePod(ctx context.Context, p *corev1.Pod, force bool) {
	if r.cfg.DryRun {
		log.Printf("[dry-run] would delete pod %s/%s force=%v", p.Namespace, p.Name, force)
		return
	}
	opts := metav1.DeleteOptions{}
	if force {
		zero := int64(0)
		opts.GracePeriodSeconds = &zero
	}
	err := r.client.CoreV1().Pods(p.Namespace).Delete(ctx, p.Name, opts)
	if err != nil && !apierrors.IsNotFound(err) {
		log.Printf("ERROR deleting pod %s/%s: %v", p.Namespace, p.Name, err)
		return
	}
	log.Printf("deleted pod %s/%s force=%v", p.Namespace, p.Name, force)
}

func negateSelector(_ string) string {
	return ""
}
