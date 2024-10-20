FROM drnay/kmttg:latest

USER root

RUN apk update
RUN apk upgrade
RUN apk add mediainfo mkvtoolnix mplayer \
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
 && apk --no-cache add --virtual=builddeps autoconf automake libtool git ffmpeg-dev wget tar build-base \
 && wget http://prdownloads.sourceforge.net/argtable/argtable2-13.tar.gz \
 && tar xzf argtable2-13.tar.gz \
 && cd argtable2-13/ && ./configure && make && make install \
 && cd /tmp && git clone https://github.com/erikkaashoek/Comskip.git \
 && cd Comskip && ./autogen.sh && ./configure && make && make install \
 && wget -O /opt/PlexComskip.py https://raw.githubusercontent.com/ekim1337/PlexComskip/master/PlexComskip.py \
 && wget -O /opt/PlexComskip.conf https://raw.githubusercontent.com/ekim1337/PlexComskip/master/PlexComskip.conf.example \
 && sed -i "s#/usr/local/bin/ffmpeg#/usr/bin/ffmpeg#g" /opt/PlexComskip.conf \
 && sed -i "/forensics/s/True/False/g" /opt/PlexComskip.conf \
 && apk del builddeps \
 && rm -rf /var/cache/apk/* /tmp/* /tmp/.[!.]*
COPY --from=plexinc/pms-docker /usr/lib/plexmediaserver/Resources/comskip.ini /opt/comskip.ini

# Get the latest kmttg version
RUN curl -L https://sourceforge.net/projects/kmttg/files/latest/download | busybox unzip -o - \
    && chmod +x /home/kmttg/app/kmttg

RUN apk add gettext
ENV APP_DIR /home/kmttg/app
ENV INPUT_DIR /mnt/kmttg/input
ENV OUTPUT_DIR /mnt/kmttg/output
ENV TOOLS_DIR /usr/local/bin 
COPY --chown=kmttg:kmttg kmttg.sh .

USER kmttg

CMD ["/home/kmttg/app/kmttg.sh"]
