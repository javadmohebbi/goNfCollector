<!-- # Grafana

- Dir to create
```mkdir /opt/nfcollector/vendors/grafana -v```
- Run container
```
docker run -d --network tick-net \
--user $(id -u) \
-p 3000:3000 \
-v /opt/nfcollector/vendors/grafana:/var/lib/grafana \
--name=grafana \
grafana/grafana
``` -->