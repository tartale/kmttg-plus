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

RUN apk add gettext
ENV APP_DIR /home/kmttg/app
ENV OUTPUT_DIR /mnt/kmttg/output
ENV TOOLS_DIR /usr/local/bin 
COPY --chown=kmttg:kmttg auto.ini .
COPY --chown=kmttg:kmttg config.ini.personal .
COPY --chown=kmttg:kmttg config.ini.template .
COPY --chown=kmttg:kmttg comskip.ini.us-ota ./comskip.ini

RUN envsubst < ./config.ini.personal > ./config.ini \
 && envsubst < ./config.ini.template >> ./config.ini \
 && ln -f -s /mnt/kmttg/output/auto.history ./ \
 && ln -f -s /mnt/kmttg/output/auto.log.0 ./

USER root

CMD ["/home/kmttg/app/kmttg", "-a"]

### unsuccessful attempt to install handbrake vvv
# RUN apk add flatpak
# RUN flatpak remote-add --user --if-not-exists flathub https://flathub.org/repo/flathub.flatpakrepo

# RUN apk add --no-cache git \
#  && git clone https://github.com/HandBrake/HandBrake.git
# RUN apk add --no-cache meson nasm autoconf libtool pkgconfig \ 
#  && apk add --no-cache cmake=3.23.1-r0 --repository="http://dl-cdn.alpinelinux.org/alpine/edge/main" \
#  && cmake --version
# RUN cd HandBrake \
#  && ./configure --launch-jobs=$(nproc) --launch --disable-gtk \
#  && make --directory=/opt/handbrake install

    # libc-dev \
    # gcc \
    # gawk \
    # make \
    # python \
    # autoconf \
    # automake \
    # libtool \
    # git \
    # tar \
    # build-base \
    # libdvdread \
    # x264 \
    # x265 \
    # meson \
    # nasm \
    # cmake
