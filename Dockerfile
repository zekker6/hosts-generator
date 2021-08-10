FROM scratch

COPY traefik-hosts-generator /

ENTRYPOINT [ "/traefik-hosts-generator" ]