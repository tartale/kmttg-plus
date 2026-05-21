FROM jlesage/handbrake

RUN apk update
RUN apk upgrade

RUN apk add --no-cache g++ make curl
RUN curl -L -o ./tivodecode-ng.tar.gz https://github.com/wmcbrine/tivodecode-ng/archive/refs/tags/0.5.tar.gz \
 && tar xvfz tivodecode-ng.tar.gz -C /opt/ && rm tivodecode-ng.tar.gz \
 && cd /opt/tivodecode-ng-0.5 \
 && ./configure \
 && make \
 && make install 

VOLUME /sys/fs/cgroup

RUN apk add --no-cache openrc avahi avahi-dev avahi-tools dbus \
 && echo " @testing https://dl-cdn.alpinelinux.org/alpine/edge/testing " >> /etc/apk/repositories \
 && apk add avahi2dns@testing \
 && mkdir -p /run/openrc \
 && touch /run/openrc/softlevel \
 && openrc

RUN apk add --no-cache doas \
 && adduser -D kmttg -s /bin/bash \
 && echo 'permit nopass :wheel' >> /etc/doas.conf \
 && addgroup kmttg wheel \
 && mkdir -p /mnt/kmttg \
 && chown -R kmttg:kmttg /mnt/kmttg

COPY --chown=kmttg:kmttg .bashrc.docker /home/kmttg/app/.bashrc
COPY --chown=kmttg:kmttg .bashrc.docker /home/kmttg/.bashrc
COPY --chown=kmttg:kmttg input/ /home/kmttg/app/input/
COPY --chown=kmttg:kmttg dist/kmttg /home/kmttg/app
COPY --chown=kmttg:kmttg dist/kmttg.sh /home/kmttg/app

ARG KMTTG_VERSION="v0.0.1"
ENV KMTTG_LOG_LEVEL="DEBUG"
ENV KMTTG_MEDIA_ACCESS_KEY=""
ENV KMTTG_PORT=7676

USER kmttg
WORKDIR /home/kmttg

CMD [ "/home/kmttg/app/kmttg.sh" ]
