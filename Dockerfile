FROM scratch

COPY hosts-generator /

ENTRYPOINT [ "/hosts-generator" ]
