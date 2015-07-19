# Flamingo

![Flamingo Logo](/logo.png)

Flamingo is a lightweight contextualization tool that aims to handle
initialization of cloud instances.

It is meant to be a replacement for cloud-init in Atomic Host, a lightweight 
operating system designed to run applications in Docker containers with 
orchestration capabilities (Kubernetes). It has 3 different 
flavors; Fedora, RHEL, and CentOS.

# Goals

The aim of this project is to create a fast, documented, extensible
tool to use

## Why not use cloud-init?

- In order for an image to be built with cloud-init as its contextualization tool,
 it needs to contain the python interpreter, it's dependencies and the dependencies of
 cloud-init. It fattens up the image and increases complexity for the distributions.
- Due to the dynamic nature of scripting languages they are slow. In an I/O boud 
application this is not a big problem however initializing VM images, is different.
- The documentation is lacking at best. The tests are inadequate and it is a big chunk
of initialization scripts. Making it hard to work with cloud-init as a developer and
maintainer.

## Solution

- Golang excels at building big applications (e.g. distributed computing). 
  It provides the necessary tooling for the job. It doesn't have long 
  dependency chains like interpreted languages and it is fast.

- Flamingo will essentially be a single binary coupled with all of its
  dependencies. In addition, cloud images will be smaller as well.

- Testability, Documentation and Extensibility is heavily emphasized in Flamingo.

# Status

Flamingo is in **Alpha** stage at the moment. Build images containing
Flamingo as the default contextualization tool (swapped with cloud-init)
will be available soon. In the meanwhile API is volatile, so expect changes!

## What has been done

- Docker Test Environment
- User & Password Management
- Running Custom Scripts
- SSH Configuration
- Manage Services at Boot
- Manage and Modify IPtables/Filter Rules
- RESTful Client

## What is being worked on

- Meta-Data
- Cloud-Config
- Config-Drive
- Integration with Systemd (using config-drive.mount instead of manual checks)
- Test Images

# Discussions

In the meanwhile if you'd like to share your opinions, learn more,
or contribute please feel free to open an issue, mail to centos-devel, atomic,
user-list or come to #atomic #centos-devel IRC channels to have a chat.
