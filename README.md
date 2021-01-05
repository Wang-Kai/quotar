# quotar

quotar 是一个部署在 NFS 虚拟机节点上的一个常驻进程，该虚拟机必须绑定 xfs 文件系统的磁盘。quotar 利用 xfs 可以限制目录大小的能力，来实现具备 `limit` 能力的 K8s PVC 分配。

