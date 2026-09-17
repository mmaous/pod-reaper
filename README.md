# pod-reaper

[![Build Status](https://github.com/mmaous/pod-reaper/actions/workflows/release.yml/badge.svg)](https://github.com/mmaous/pod-reaper/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmaous/pod-reaper)](https://goreportcard.com/report/github.com/mmaous/pod-reaper)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Dead simple k8s janitor for homelabs and laptops. 

Ever sleep your laptop running minikube, or suspend your VMs, and wake up to half your pods stuck in `0/1 Running` or `CrashLoopBackOff` with `<invalid> ago` timestamps? k8s networking and storage providers often hate huge time jumps. 

`pod-reaper` runs once on boot, finds these zombie pods, and deletes them. Deployments instantly spin up fresh pods with working CNI/CSI attachments. Boom, cluster healed.

## What it kills

| State | Condition |
|---|---|
| Time Traveled | Pod has a negative age (future timestamp) caused by host sleep/wake cycle |
| Terminating | Stuck deleting past `2m` (force kills `grace=0`) |
| Pending | Stuck unscheduled or waiting on CNI past `2m` |
| Running | `NotReady` (0/1, 0/2) past `2m` |
| CrashLoopBackOff | Restarts exceed `3` |
| ContainerCreating | Pod is stuck trying to create a container past `2m` |
| CreateContainerError | Pod is stuck trying to create a container |
| CreateContainerConfigError | Pod is stuck trying to create a container due to bad config |

## How to use

Run it as a oneshot systemd service on your nodes. 3 minutes after the VM boots, the systemd unit will:
1. Restart `containerd` and `kubelet` to unfreeze the network and static pods.
2. Wait 60 seconds for the cluster to re-sync.
3. Automatically ignore any static pods (identified by the `kubernetes.io/config.mirror` annotation).
4. Force delete the remaining broken application pods.

1. Download release tarball:
```bash
# Example for arm64 (M1/M2/arm VMs)
wget https://github.com/mmaous/pod-reaper/releases/download/v1.0.0/pod-reaper-linux-arm64.tar.gz
tar -xzvf pod-reaper-linux-arm64.tar.gz
```

2. Install binary and systemd units:
```bash
sudo mv pod-reaper /usr/local/bin/
sudo chmod +x /usr/local/bin/pod-reaper
sudo mv pod-reaper.service pod-reaper.timer /etc/systemd/system/
```

3. Enable timer for automatic boot cleanup:
```bash
sudo systemctl daemon-reload
sudo systemctl enable --now pod-reaper.timer
```

4. Run manually on-demand:
```bash
sudo systemctl start pod-reaper.service
```

Service starts in `DRY_RUN=true`. Check `journalctl -u pod-reaper.service` to make sure it's not killing anything important. When ready, edit `/etc/systemd/system/pod-reaper.service`, flip to `DRY_RUN=false`, and `systemctl daemon-reload`.

> **Note:** Don't run this in production. Masking crashloops in prod is a bad idea. This is a sledgehammer for dev labs.
