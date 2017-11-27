# PiaaS

A heroku like PaaS for the Raspberry Pi or any linux system.

## Building

+ First of all fetch all the packages without installing.
    ```shell
    go get -d github.com/magrandera/PiaaS-go
    ```
+ `cd` into the project directory
    ```shell
    cd $GOPATH/src/github.com/magrandera/PiaaS-go
    ```
+ install dependencies with glide
    ```shell
    glide install --strip-vendor
    ```
+ build the binary using the makefile
    ```shell
    make
    ```
+ binary is in the build folder

## Deployment

Full explanation available [here](doc/GUIDE.md)

## Built With

* [Glide](https://github.com/Masterminds/glide) - Dependency Management
* [cli](https://github.com/urfave/cli) - CLI framework
* [moby](https://github.com/moby/moby) - Docker repository

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors

* **Malte Granderath** - *Initial work* - [magrandera](https://github.com/magrandera)

See also the list of [contributors](https://github.com/magrandera/PiaaS-go/graphs/contributors) who participated in this project.

## License

This project is licensed under the Apache2.0 License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments
