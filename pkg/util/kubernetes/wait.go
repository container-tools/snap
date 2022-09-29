package kubernetes

import (
	"context"
	"errors"
	"log"
	"time"

	snapclient "github.com/container-tools/snap/pkg/client"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func WaitForPodReady(ctx context.Context, client snapclient.Client, ns string, labelSelector string) (string, error) {
	for {
		pod, err := getReadyPod(ctx, client, ns, labelSelector)
		if err != nil {
			log.Print("error looking up target pod: ", err)
		}
		if err == nil && pod != "" {
			return pod, nil
		}

		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			return "", errors.New("unable to find target pod")
		}
	}
}

func getReadyPod(ctx context.Context, client snapclient.Client, ns string, labelSelector string) (string, error) {
	podList, err := client.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", err
	}
	for _, pod := range podList.Items {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.ContainersReady && condition.Status == corev1.ConditionTrue {
				return pod.Name, nil
			}
		}
	}
	return "", nil
}
