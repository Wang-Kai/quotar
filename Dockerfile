# Build the manager binary
# FROM golang:1.13 as builder

# COPY . /workspace/
# WORKDIR /workspace
# RUN go build -v -o /bin/nfs-provisioner nfs-provisioner/*.go

FROM centos:7.6.1810
# COPY --from=builder /bin//nfs-provisioner /bin/nfs-provisioner
COPY bin/nfs-provisioner /bin/nfs-provisioner
CMD ["nfs-provisioner"]