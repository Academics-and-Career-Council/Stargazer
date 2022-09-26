# Note: You can use any Debian/Ubuntu based image you want. 
FROM mcr.microsoft.com/vscode/devcontainers/base:0-bullseye

# [Option] Install zsh
ARG INSTALL_ZSH="true"
# [Option] Upgrade OS packages to their latest versions
ARG UPGRADE_PACKAGES="false"
# [Option] Enable non-root Docker access in container
ARG ENABLE_NONROOT_DOCKER="true"
# [Option] Use the OSS Moby CLI instead of the licensed Docker CLI
ARG USE_MOBY="true"

# Enable new "BUILDKIT" mode for Docker CLI
ENV DOCKER_BUILDKIT=1

# Install needed packages and setup non-root user. Use a separate RUN statement to add your
# own dependencies. A user of "automatic" attempts to reuse an user ID if one already exists.
ARG USERNAME=automatic
ARG USER_UID=1000
ARG USER_GID=$USER_UID
COPY library-scripts/*.sh /tmp/library-scripts/
RUN apt-get update \
    && /bin/bash /tmp/library-scripts/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" "true" "true" \
    # Use Docker script from script library to set things up
    && /bin/bash /tmp/library-scripts/docker-debian.sh "${ENABLE_NONROOT_DOCKER}" "/var/run/docker-host.sock" "/var/run/docker.sock" "${USERNAME}" \
    # Clean up
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts/ \
    && wget https://download.libsodium.org/libsodium/releases/libsodium-1.0.18.tar.gz \
    && wget https://download.libsodium.org/libsodium/releases/libsodium-1.0.18.tar.gz.sig \
    && wget https://download.libsodium.org/jedi.gpg.asc \
    && gpg --import jedi.gpg.asc \
    && gpg --verify libsodium-1.0.18.tar.gz.sig libsodium-1.0.18.tar.gz \
    && tar zxvf libsodium-1.0.18.tar.gz \
    && cd libsodium-1.0.8 \
    && ./configure; make check \
    && sudo make install \
    && sudo ldconfig \
    && git clone https://github.com/zeromq/libzmq.git \
    && cd libzmq \
    && ./autogen.sh \
    && ./configure --with-libsodium \
    && make check \
    && sudo make install \
    && sudo ldconfig \
    && cd .. \
    && git clone https://github.com/zeromq/czmq.git \
    && cd libzmq \
    && ./autogen.sh \
    && ./configure \
    && make check \
    && sudo make install \
    && sudo ldconfig

# Setting the ENTRYPOINT to docker-init.sh will configure non-root access 
# to the Docker socket. The script will also execute CMD as needed.
ENTRYPOINT [ "/usr/local/share/docker-init.sh" ]
CMD [ "sleep", "infinity" ]

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>