# Pull base image.
FROM quay.io/eris/build
MAINTAINER Eris Industries <support@erisindustries.com>

#-----------------------------------------------------------------------------
# install epm

# set the repo and copy in files
ENV REPO $GOPATH/src/github.com/eris-ltd/eris-pm
COPY . $REPO

# use glide; remove its traces
WORKDIR $REPO
RUN glide install --strip-vcs --strip-vendor && \
	glide vc

# install eris-pm
WORKDIR $REPO/cmd/epm
RUN go build -o $INSTALL_BASE/epm
RUN chown --recursive $USER:$USER $REPO
# should be able to remove /usr/lib/go as well.
RUN glide cc && rm -rf $GOPATH && apk del --no-cache --purge go git gmp-dev gcc musl-dev 

#-----------------------------------------------------------------------------
# root dir

# persist data, set user
RUN chown --recursive $USER:$USER /home/$USER
VOLUME /home/$USER/.eris
WORKDIR /home/$USER/.eris
USER $USER
CMD ["epm", "--chain", "chain:46657", "--sign", "keys:4767" ]
