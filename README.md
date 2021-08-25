# PXE-Init

## Author
Thomas Mortensson - thomasmortensson@googlemail.com

## Description
This project forms the basis of a service to enable iPXE booting of cloud-init machines backed by iSCSI volumes.

The aim is to be able to boot both VMs and physical machines in the mortcloud.com domain diskless.

The current implementation is designed to show the basis of a client-server architecture application making use of both gRPC and HTTP endpoints.

The service is currently deployed in my home environment running in K8S exposed via MetalLB onto my main VLAN. The service is supported by a DNS server, a TFTP server, an Apache HTTP server serving images and a postgres DB backend. I am able to make use of it to perform basic PXE bootup of coreos machines (either VMs or physical HW)

## Features
- [X] Provide CLI client-server architecture
- [X] Provide liveness check /health
- [ ] Provide readiness check /ready
- [X] Initial chainload endpoint /boot.ipxe
- [X] ipxe chain script /ipxe?mac=uuid=
- [X] List images gRPC
- [ ] Register machine image (partial complete)
- [X] Unit test (Not all code, enough to show how I would mock datastores etc)
- [X] Docker container for server build / run
- [X] Helm scripts generated
- [X] Uploaded to local SSL docker registry
- [X] Running in K8S
- [X] Exposed in MetalLB
- [X] Successful diskless coreos machine boot
- [X] Usecases defined and mockery used to mock tests
- [ ] SSL between client / server
- [ ] SSL between server and DB
- [ ] Integration test (pytest framework)
- [ ] CI integration (drone.yaml)
- [ ] UML documents drawn up for workflows
- [ ] UML documents drawn up for entity diagrams

The features left incomplete in this list I have had to leave either unimplemented or incomplete due to time constraints, I wish I had some more time before interview to complete these but the reality is I also have to stay out of the doghouse :)

## TODO
- Generate PUML architectural documentation for actor based workflows of service
- Generate PUML architectural documentation for entities
- Implement fully RegisterMachineImage gRPC
- Implement a RegisterImage method
- Implement a null Image
- Split server entrypoint down, split to:
	- GetArgs
	- Setup DB
	- Setup HTTP server
	- Setup gRPC server
	- Run servers
	- Shutdown
- Secure client-server communication via SSL
- Secure DB communication via SSL
- Add svc pointer for DB in helm charts
- Implement liveness and readiness checks
- Setup integration testing framework for E2E testing in pytest
- Setup drone CI on project
- Vastly increase test coverage


## Building

It is expected that this application will be built on go 1.16. I've been using an alpine base.

Makefile usage:
```
➜  pxe-init git:(master) ✗ make
help           Output this message and exit.
build          build the project
tidy           tidy dependencies
vendor         download vendored dependencies
fmt            format the code
clean          clean up built code and vendor directory
trivy          scan for vulnerabilities
lint           run the project linters
test           Run unit tests
proto          Generate grpc code from proto
docker-build   Build the pxe-init-server image.
docker-push    Push the docker image.
```

To build natively:
```
make build
```

To build a docker container:
```
make docker-build
```

## Deployment

Helm charts exist for the deployment of the pxe-init server in ./helm. I have rendered a complete example in ./helm/example.yaml. To template the chart, run:
```
helm template --release-name pxe-init -n pxe-init --values ./pxe-init/values.yaml ./pxe-init
```

## Usage

To run the server, you will need a postgres database. To interact with this database you will need a user / password with permissions on this DB.

```
postgres=# create database pxe_init;
postgres=# create user pxe_init with encrypted password 'mypass';
postgres=# grant all privileges on database pxe_init to pxe_init;
```

You can then setup the example objects that I have provided for usage with the application
```
psql -h 127.0.0.1 -U pxe_init -f example/sql/example.sql
```

Application usage is available here:
```
➜  pxe-init git:(master) ✗ ./dist/pxe-init-server serve -h
Start the pxe-init-server

Usage:
  pxe-init-server serve [flags]

Flags:
  -h, --help   help for serve

Global Flags:
      --db-host string              Database host endpoint to use for communication (default "127.0.0.1")
      --db-name string              Database name (default "pxe_init")
      --db-password string          Database password
      --db-port int                 Database port to use for communication (default 5432)
      --db-ssl-mode string          Database SSL mode choice [disable allow prefer require verify-ca verify-full] (default "disable")
      --db-user string              Database user (default "pxe_init")
      --grpc-port int               gRPC port to serve pxe-init-server on (default 5000)
      --http-port int               HTTP port to serve pxe-init-server on (default 8080)
      --pxe-forward-server string   Forward server for PXE assets (default "http://zbox.mortcloud.com")
```

At a minimum, I'd expect that you'd want to set the DB endpoint (db-host), password (db-password) and potentially the HTTP port (http-port).

```
pxe-init-server serve --db-host 127.0.0.1 --db-password mypass --http-port 80
```

Once the server is running, you should be able to query the available images from the client:
```
pxe-init-client list-images --grpc-endpoint grpc://pxe-init.mortcloud.com:5000
+--------+------------------------------------------------------------------------------------------+-------------------------------------------------------------------------------------------------+----------------------------------------------------------------------------------------------+
|  NAME  |                                          KERNEL                                          |                                             INITRD                                              |                                            ROOTFS                                            |
+--------+------------------------------------------------------------------------------------------+-------------------------------------------------------------------------------------------------+----------------------------------------------------------------------------------------------+
| coreos | http://zbox.mortcloud.com/assets/coreos/fedora-coreos-34.20210725.3.0-live-kernel-x86_64 | http://zbox.mortcloud.com/assets/coreos/fedora-coreos-34.20210725.3.0-live-initramfs.x86_64.img | http://zbox.mortcloud.com/assets/coreos/fedora-coreos-34.20210725.3.0-live-rootfs.x86_64.img |
+--------+------------------------------------------------------------------------------------------+-------------------------------------------------------------------------------------------------+----------------------------------------------------------------------------------------------+
```

You can then simulate a bootup issing the following requests:

```
➜  pxe-init git:(master) ✗ curl 'http://pxe-init.mortcloud.com/boot.ipxe'
#!ipxe
chain ipxe?uuid=${uuid}&mac=${mac:hexhyp}&domain=${domain}&hostname=${hostname}&serial=${serial}

➜  pxe-init git:(master) ✗ curl 'http://pxe-init.mortcloud.com/ipxe?uuid=2f1214fe-59ba-9f42-a5c5-1af6f124aaf7&mac=08-00-27-c3-62-83&domain=&hostname=&serial=0'
#!ipxe
kernel http://zbox.mortcloud.com/assets/coreos/fedora-coreos-34.20210725.3.0-live-kernel-x86_64 coreos.live.rootfs_url=http://zbox.mortcloud.com/assets/coreos/fedora-coreos-34.20210725.3.0-live-rootfs.x86_64.img coreos.first_boot=1 coreos.autologin
initrd http://zbox.mortcloud.com/assets/coreos/fedora-coreos-34.20210725.3.0-live-initramfs.x86_64.img
boot%

```