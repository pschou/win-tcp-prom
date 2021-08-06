# win-tcp-prom

Simple TCP stats exporter to send to prom collector.

Usage:

```bash
C:\Users\nobody\Downloads>win-tcp-prom.exe http://my.prom.collector:9550/collector/tcp/test
Prometheus TCP Metrics Scraper, Written by Paul Schou
Hostname = rdp
URI = http://my.prom.collector:9550/collector/tcp/test
Collecting TCP stats...
  Active_Opens 123305
  Passive_Opens 617
  Failed_Connection_Attempts 884
  Reset_Connections 8766
  Current_Connections 42
  Segments_Received 25402704
  Segments_Sent 10265836
  Segments_Retransmitted 0
Sending metrics...
Sent!  Response: 200
```
