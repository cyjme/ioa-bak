FROM alpine:latest
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ADD release.tar.gz /var
ADD etcd-v3.3.12-linux-amd64.tar.gz /var
COPY run.sh /

EXPOSE 9991
EXPOSE 9992

ENTRYPOINT ["./run.sh"]
