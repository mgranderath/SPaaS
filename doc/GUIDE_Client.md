# Client Guide

## General Info

##### We recommend to have a ssh key setup with the server.

Projects have to contain a Procfile that contain the following:

```
web: %command that starts the server%
```

Examples of projects:
- [Python](https://github.com/heroku/python-sample)
- [NodeJs](https://github.com/heroku/node-js-sample)

## Installing PiaaS

You can either download the precompiled binary called `PiaaS_cli` or use go install to custom compile it. You have to have the `$PATH` variable set to the directory to make the binary be executable system wide.

Once the server has been started you have setup the client. Do this by running the following:

```shell
PiaaS setup
```

## Creating a new application

To create a new application, execute the following:

```shell
PiaaS add %name
```

That will return output similar to the below.

```shell
-----> Task: Creating test2
-----> Info: Repository path: /home/malte/PiaaS-Data/Applications/test2/repo
-----> Success: Creating test2
```

To deploy the application take the above Repository path and do the following in the directory that you want to deploy:
```shell
git init
git remote add origin username@server:%repository_path%
git add .
git push -u origin master
```

## Deployment

There is 2 ways that a project can be deployed:
1. Git pushing the repository will automatically start the deployment process
2. We can call a command that will trigger a deploy

```shell
PiaaS deploy %name%
```

## Starting and Stopping Applications

You have to be logged in on the Raspberry Pi or server.
To stop an application:
```shell
PiaaS stop %name%
```
To start an application:
```shell
PiaaS start %name%
```

## Removing an application
You have to be logged in on the Raspberry Pi or server.
To remove an application:
```shell
PiaaS remove %name%
```

## Listing all applications

You can list all the existing applications by using the following result.
```shell
PiaaS list
```