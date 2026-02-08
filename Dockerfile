## copied/modified from: https://github.com/drnay/docker-kmttg/blob/master/Dockerfile
# Alpine Linux with Oracle JRE
FROM sgrio/java:jre_8

RUN apk update
RUN apk upgrade

ENV APP_DIR /home/kmttg/app
ENV MOUNT_DIR /mnt/kmttg
ENV TOOLS_DIR /usr/local/bin

RUN mkdir -p ${APP_DIR} ${MOUNT_DIR} ${TOOLS_DIR} \
 && chown kmttg:kmttg ${APP_DIR} ${MOUNT_DIR} ${TOOLS_DIR}

RUN apk add gettext mediainfo mkvtoolnix mplayer \
 && ln -s /usr/bin/mediainfo ${TOOLS_DIR}/mediainfo \
 && ln -s /usr/bin/mencoder ${TOOLS_DIR}/mencoder \
 && ln -s /usr/bin/ffmpeg ${TOOLS_DIR}/ffmpeg

RUN curl -L -o ./AtomicParsleyAlpine.zip https://github.com/wez/atomicparsley/releases/download/20210715.151551.e7ad03a/AtomicParsleyAlpine.zip \
 && unzip -d ${TOOLS_DIR} ./AtomicParsleyAlpine.zip \
 && rm ./AtomicParsleyAlpine.zip

RUN apk add --no-cache g++ make \
 && curl -L -o ./tivodecode-ng.tar.gz https://github.com/wmcbrine/tivodecode-ng/archive/refs/tags/0.5.tar.gz \
 && tar xvfz tivodecode-ng.tar.gz -C /opt/ && rm tivodecode-ng.tar.gz \
 && cd /opt/tivodecode-ng-0.5 \
 && ./configure \
 && make \
 && make install 

RUN apk --no-cache add python ffmpeg tzdata bash \
 && apk --no-cache add --virtual=builddeps autoconf automake libtool git ffmpeg-dev wget tar build-base

RUN cd /root \
 && wget http://prdownloads.sourceforge.net/argtable/argtable2-13.tar.gz \
 && tar xzf argtable2-13.tar.gz \
 && cd argtable2-13/ && ./configure && make && make install

RUN cd /tmp \
 && git clone https://github.com/erikkaashoek/Comskip.git \
 && cd Comskip && ./autogen.sh && ./configure && make && make install \
 && wget -O /opt/PlexComskip.py https://raw.githubusercontent.com/ekim1337/PlexComskip/master/PlexComskip.py \
 && wget -O /opt/PlexComskip.conf https://raw.githubusercontent.com/ekim1337/PlexComskip/master/PlexComskip.conf.example \
 && sed -i "s#${TOOLS_DIR}/ffmpeg#/usr/bin/ffmpeg#g" /opt/PlexComskip.conf \
 && sed -i "/forensics/s/True/False/g" /opt/PlexComskip.conf \
 && apk del builddeps \
 && rm -rf /var/cache/apk/* /tmp/* /tmp/.[!.]*

COPY --from=plexinc/pms-docker /usr/lib/plexmediaserver/Resources/comskip.ini /opt/comskip.ini

USER kmttg

COPY --chown=kmttg:kmttg java/release/kmttg ${APP_DIR}/
COPY --chown=kmttg:kmttg java/release/kmttg.jar ${APP_DIR}/
COPY --chown=kmttg:kmttg input/* ${APP_DIR}/input/
COPY --chown=kmttg:kmttg kmttg.sh ${APP_DIR}/

CMD ["/bin/bash", "-c", "${APP_DIR}/kmttg.sh"]
