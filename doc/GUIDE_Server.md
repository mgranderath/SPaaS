# Server Guide

## Prerequisites

- Docker has to be installed

## Installing

You can either download the precompiled binary called `PiaaS_cli` or use go install to custom compile it. You have to have the `$PATH` variable set to the directory to make the binary be executable system wide.

To run the server:
```shell
PiaaS # if installed system wide
./PiaaS # if installed locally
```

## Configuration

You can create a configuration file `$HOME\.config\piaas\config` that can contain the following options in JSON format.

#### **nginx**: can be set to true or false

If set to false no nginx proxy will be setup and youll be able to access the deployed applications on respective ports. Otherwise nginx will reverse proxy and you'll be able to access the project under a respective url like `projectname.yourdomain.com`. For this you cannot have any other application running on port 80.