# dcc-ex-exporter

A simpler [Prometheus](https://prometheus.io) exporter for [DCC-EX](https://www.dcc-ex.com) to expose track power usage.

## Compiling

```bash
go build
```

This will generate `dcc-ex-exporter`.

## Usage

The exporter accepts the following command line arguments:

| Option       | Required | Description                                                 |
|--------------|----------|-------------------------------------------------------------|
| --dccex-host | Yes      | The hostname or IP address of your DCC-EX command station   |
| --dccex-port | No       | The port used by your DCC-EX command station (default 2560) |
| --port       | No       | The port the exporter will listen on (default 9378)         |

Example:

```bash
./dcc-ex-exporter --dccex-host dccex
```

## Installation

Clone this repository:

```bash
git clone https://github.com/neilmunday/dcc-ex-exporter
```

Build the executable:

```bash
cd dcc-ex-exporter
go build
```

### Install executable

```bash
sudo mkdir -p /usr/local/bin
sudo cp dcc-ex-exporter /usr/local/bin
```

### Create systemd service

In this repository you will find `dcc-ex-exporter.service`. Edit this file and set the executable location and command line options to suit your set-up.

Now copy the file to `/etc/systemd/system`:

```bash
sudo cp dcc-ex-exporter.service /etc/systemd/system/
```

Reload `systemd`:

```bash
sudo systemctl daemon-reload
```

Check the exporter is running ok:

```bash
sudo systemctl status dcc-ex-exporter.service
```

You should see output similar to:

```
● dcc-ex-exporter.service - DCC-EX Prometheus Exporter
   Loaded: loaded (/etc/systemd/system/dcc-ex-exporter.service; disabled; vendor preset: disabled)
   Active: active (running) since Wed 2026-06-03 21:55:39 BST; 41min ago
 Main PID: 1084801 (dcc-ex-exporter)
    Tasks: 6 (limit: 22895)
   Memory: 7.6M
   CGroup: /system.slice/dcc-ex-exporter.service
           └─1084801 /opt/bin/dcc-ex-exporter --dccex-host dccex

Jun 03 21:55:39 myhost systemd[1]: Started DCC-EX Prometheus Exporter.
Jun 03 21:55:39 myhost dcc-ex-exporter[1084801]: 2026/06/03 21:55:39 Exporter running on :9378/metrics
Jun 03 21:55:40 myhost dcc-ex-exporter[1084801]: 2026/06/03 21:55:40 Connected to DCC-EX at dccex:2560
```

Now enable the service at start-up:

```bash
sudo systemctl enable dcc-ex-exporter
```

### Adding to Prometheus

In your `prometheus.yml` file add to your `scrape_configs` section:

```yml
scrape_configs:

  - job_name: 'dcc-ex-exporter'
    static_configs:
      - targets: ['localhost:9378']
```

Now restart Prometheus:

```bash
sudo systemctl restart prometheus
```

You should now be able see your DCC-EX metrics are being scraped by Prometheus at https://localhost:9090/targets

From here you can create your own dashboard using [Grafana](https://www.grafana.com) for example.