# earthquake
Automated live migration of the game Quake3

## Environment
Several packages are needed to setup an environment for running earthquake.
The go dependency management included should be enough for `go build` to work.

## Building a quake container image
Obviously, quake3 is not free and I am not endorsing breaking copyright.
Specifically, a certain file must be obtained from the disk.
On ubuntu the rest of the code can be pulled and installed from the `ioquake3*` packages.

See https://github.com/jberrenberg/docker-quake3 for Docker tooling for creating the quake3 server image.
