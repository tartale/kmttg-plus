FROM drnay/kmttg:latest

USER root

RUN apk update
RUN apk upgrade

RUN apk add gettext mediainfo mkvtoolnix mplayer \
 && ln -s /usr/bin/mediainfo /usr/local/bin/mediainfo \
 && ln -s /usr/bin/mencoder /usr/local/bin/mencoder \
 && ln -s /usr/bin/ffmpeg /usr/local/bin/ffmpeg

RUN curl -L -o ./AtomicParsleyAlpine.zip https://github.com/wez/atomicparsley/releases/download/20210715.151551.e7ad03a/AtomicParsleyAlpine.zip \
 && unzip -d /usr/local/bin ./AtomicParsleyAlpine.zip \
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
 && sed -i "s#/usr/local/bin/ffmpeg#/usr/bin/ffmpeg#g" /opt/PlexComskip.conf \
 && sed -i "/forensics/s/True/False/g" /opt/PlexComskip.conf \
 && apk del builddeps \
 && rm -rf /var/cache/apk/* /tmp/* /tmp/.[!.]*

COPY --from=plexinc/pms-docker /usr/lib/plexmediaserver/Resources/comskip.ini /opt/comskip.ini

COPY --chown=kmttg:kmttg .bashrc.docker /home/kmttg/app/.bashrc

RUN apk add doas \
 && echo 'permit nopass :wheel' >> /etc/doas.conf \
 && addgroup kmttg wheel

ARG KMTTG_VERSION
COPY --chown=kmttg:kmttg input/* /home/kmttg/app/input/
COPY --chown=kmttg:kmttg kmttg.sh /home/kmttg/app/
COPY --chown=kmttg:kmttg java/dist/kmttg_${KMTTG_VERSION}.zip /home/kmttg/app/

RUN cd /home/kmttg/app \
 && unzip -o -q kmttg_${KMTTG_VERSION}.zip \
 && chown -R kmttg:kmttg /home/kmttg/app \
 && rm -rf kmttg_${KMTTG_VERSION}.zip

VOLUME [ /mnt/kmttg ]

USER kmttg
CMD ["/bin/bash", "-c", "/home/kmttg/app/kmttg.sh", "-a"]
