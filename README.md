# Raspberry Pi Poudriere Cluster with Ansible + Webhook

This repository automates setting up a FreeBSD-based Raspberry Pi 4 cluster for building FreeBSD packages with Poudriere. It includes:

- A jail-based poudriere build container on the master node
- Build nodes mounting NFS exports from the master
- Automated testing and report publishing
- GitHub webhook trigger via Go-based listener

## Requirements

- Raspberry Pi 4 devices with FreeBSD installed
- One `pi-master`, three build nodes (`pi-node1`, `pi-node2`, `pi-node3`)
- SSH access set up between nodes
- Ansible and Go installed on `pi-master`

## Inventory File

See `hosts.ini` for a sample inventory setup.

## Usage

1. **Provision Poudriere Jail on Master:**

```sh
ansible-playbook -i hosts.ini poudriere-jail.yml
```

2. **Configure Build Nodes:**

```sh
ansible-playbook -i hosts.ini build-nodes-setup.yml
```

3. **Run Test Build:**

```sh
ansible-playbook -i hosts.ini test-build-run.yml
```

4. **Set Up nginx to Publish Build Reports:**

```sh
ansible-playbook -i hosts.ini nginx-report-setup.yml
```

5. **Build and Run the Webhook Listener:**

```sh
export WEBHOOK_SECRET=yoursecret
go build -o webhook webhook.go
./webhook
```

You can daemonize the webhook using `daemon(8)` or the included rc script.

## Webhook

Set up a GitHub webhook pointing to:

```
http://<pi-master-ip>:9000
```

Use the same secret as `WEBHOOK_SECRET`.

