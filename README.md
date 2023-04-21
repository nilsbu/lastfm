# Run

You need a file to configure the environment variables.

```
TLS_CERT_BASE=... # path to the directory containing the TLS certificates, e.g. /etc/letsencrypt
TLS_CERT_REL=... # path to the TLS certificate relative to TLS_CERT_BASE, e.g. live/<URL>
DATA_PATH=... # path to the directory containing the data, e.g. ~/.lastfm
```

The split between TLS_CERT_BASE and TLS_CERT_REL is necessary because the TLS certificates are mounted into the container as a volume. This can be problematic if the actual directory contains symbolic links. Ensure that the links are relative and that TLS_CERT_BASE is the directory containing both the link and actual certificate.

Build with docker using the following command:

```
docker-compose --env-file <config.yml> up --build -d
```

This will create two containers. One that contains the server and a worker that cleans up the data.
