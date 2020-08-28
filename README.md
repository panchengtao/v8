# V8 Bindings for Go [![Build Status](https://travis-ci.org/augustoroman/v8.svg?branch=master)](https://travis-ci.org/augustoroman/v8) [![Go Report Card](https://goreportcard.com/badge/github.com/augustoroman/v8)](https://goreportcard.com/report/github.com/augustoroman/v8) [![GoDoc](https://godoc.org/github.com/augustoroman/v8?status.svg)](https://godoc.org/github.com/augustoroman/v8)

The v8 bindings allow a user to execute javascript from within a go executable.

The bindings are tested to work with several recent v8 builds matching the
Chrome builds 54 - 60 (see the .travis.yml file for specific versions). For
example, Chrome 59 (dev branch) uses v8 5.9.211.4 when this was written.

Note that v8 releases match the Chrome release timeline:
Chrome 48 corresponds to v8 4.8.\*, Chrome 49 matches v8 4.9.\*. You can see
the table of current chrome and the associated v8 releases at:

http://omahaproxy.appspot.com/

# Using a pre-compiled v8

v8 is very slow to compile, it's a large project. If you want to go that route, there are building instructions below.

Fortunately, there's a project that pre-builds v8 for various platforms. It's packaged as a ruby gem called [libv8](https://rubygems.org/gems/libv8).

```bash
# Find the appropriate gem version for your OS,
# visit: https://rubygems.org/gems/libv8/versions

# Download the gem
# MacOS Sierra is darwin-16, for v8 6.3.292.48.1 it looks like:
curl https://rubygems.org/downloads/libv8-6.3.292.48.1-x86_64-darwin-16.gem > libv8.gem

# Extract the gem (it's a tarball)
tar -xf libv8.gem

# Extract the `data.tar.gz` within
cd libv8-6.3.292.48.1-x86_64-darwin-16
tar -xzf data.tar.gz

# Symlink the compiled libraries and includes
ln -s $(pwd)/data/vendor/v8/include ${SRCDIR}/darwin/include
ln -s $(pwd)/data/vendor/v8/out/x64.release ${SRCDIR}/darwin/libv8

# Run the tests to make sure everything works
go test ./...
```

# Using docker (linux only)

For linux builds, you can use pre-built libraries or build your own.

## Pre-built versions

To use a pre-built library, select the desired v8 version from https://hub.docker.com/r/augustoroman/v8-lib/tags/ and then run:

```bash
# Select the v8 version to use:
export V8_VERSION=6.3.292.48.1
docker pull augustoroman/v8-lib:$V8_VERSION           # Download the image, updating if necessary.
docker rm v8 ||:                                      # Cleanup from before if necessary.
docker run --name v8 augustoroman/v8-lib:$V8_VERSION  # Run the image to provide access to the files.
docker cp v8:/v8/include include/                     # Copy the include files.
docker cp v8:/v8/lib libv8/                           # Copy the library fiels.
```

## Build your own via docker

This takes a lot longer, but is still easy:

```bash
export V8_VERSION=6.3.292.48.1
docker build --build-arg V8_VERSION=$V8_VERSION --tag augustoroman/v8-lib:$V8_VERSION docker-v8-lib/
```

and then extract the files as above:

```bash
docker rm v8 ||:                                      # Cleanup from before if necessary.
docker run --name v8 augustoroman/v8-lib:$V8_VERSION  # Run the image to provide access to the files.
docker cp v8:/v8/include include/                     # Copy the include files.
docker cp v8:/v8/lib libv8/                           # Copy the library fiels.
```

# Reference

Also relevant is the v8 API release changes doc:

https://docs.google.com/document/d/1g8JFi8T_oAE_7uAri7Njtig7fKaPDfotU6huOa1alds/edit

# Credits

This work is based off of several existing libraries:

* https://github.com/fluxio/go-v8
* https://github.com/kingland/go-v8
* https://github.com/mattn/go-v8
