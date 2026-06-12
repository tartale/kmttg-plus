FROM alpine:3.23.4

RUN apk update \
 && apk upgrade \
 && apk add \
  avahi \
  avahi-tools \
  bash \
  curl \
  dbus \
  doas \
  ffmpeg \
  mkvtoolnix

RUN adduser -D kmttg -s /bin/bash \
 && echo 'permit nopass :wheel' >> /etc/doas.conf \
 && addgroup kmttg wheel \
 && mkdir -p /mnt/kmttg \
 && chown -R kmttg:kmttg /mnt/kmttg

COPY --chown=kmttg:kmttg .bashrc.docker /home/kmttg/app/.bashrc
COPY --chown=kmttg:kmttg .bashrc.docker /home/kmttg/.bashrc
COPY --chown=kmttg:kmttg dist/kmttg /home/kmttg/app
COPY --chown=kmttg:kmttg dist/kmttg.sh /home/kmttg/app

ARG KMTTG_VERSION="v0.0.1"
ENV KMTTG_LOG_LEVEL="DEBUG"
ENV KMTTG_MEDIA_ACCESS_KEY=""
ENV KMTTG_WEBUI_DIR=""

EXPOSE 7676
USER kmttg
WORKDIR /home/kmttg

CMD [ "/home/kmttg/app/kmttg.sh" ]
