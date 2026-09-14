# pod-reaper 💀

Dead simple k8s janitor for homelabs and laptops. 

Ever sleep your laptop running minikube, or suspend your VMs, and wake up to half your pods stuck in `0/1 Running` or `CrashLoopBackOff` with `<invalid> ago` timestamps? k8s networking and storage providers often hate huge time jumps. 

`pod-reaper` runs once on boot, finds these zombie pods, and deletes them. Deployments instantly spin up fresh pods with working CNI/CSI attachments. Boom, cluster healed.

## What it kills

| State | Condition |
|---|---|
| Terminating | Stuck deleting past `2m` (force kills `grace=0`) |
| Pending | Stuck unscheduled or waiting on CNI past `5m` |
| Running | `NotReady` (0/1, 0/2) past `5m` |
| CrashLoopBackOff | Restarts exceed `6` |

## How to use

Run it as a oneshot systemd service on your nodes so it fires 3 minutes after the VM boots.

1. Build it:
```bash
go build -o pod-reaper .
```

2. Drop the binary on your node:
```bash
sudo mv pod-reaper /usr/local/bin/
sudo chmod +x /usr/local/bin/pod-reaper
```

3. Drop the systemd unit `pod-reaper.service` into `/etc/systemd/system/`.

4. Enable it:
```bash
sudo systemctl daemon-reload
sudo systemctl enable pod-reaper.service
```

Service starts in `DRY_RUN=true`. Check `journalctl -u pod-reaper.service` to make sure it's not killing anything important. When ready, edit `/etc/systemd/system/pod-reaper.service`, flip to `DRY_RUN=false`, and `systemctl daemon-reload`.

> **Note:** Don't run this in production. Masking crashloops in prod is a bad idea. This is a sledgehammer for dev labs.
