# This how we want to name the binary output
BINARY=afro

# These are the values we want to pass for VERSION and BUILD
VERSION=1.0.0
BUILD=`git rev-parse HEAD`

# These are the values we want to pass for BOOL_TRUE and BOOL_FALSE
BOOL_TRUE=1
BOOL_FALSE=2

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.True=${BOOL_TRUE} -X main.False=${BOOL_FALSE}"

# Builds the project
build:
	go build ${LDFLAGS} -o ${BINARY}

# Installs our project: copies binaries
install:
	go install ${LDFLAGS}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install
