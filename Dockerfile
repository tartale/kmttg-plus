FROM drnay/kmttg:latest

USER root

RUN apk update
RUN apk upgrade
RUN apk add mkvtoolnix mplayer

RUN curl -L -o ./AtomicParsleyAlpine.zip https://github.com/wez/atomicparsley/releases/download/20210715.151551.e7ad03a/AtomicParsleyAlpine.zip \
 && unzip -d /usr/local/bin ./AtomicParsleyAlpine.zip \
 && rm ./AtomicParsleyAlpine.zip

RUN apk add --no-cache libc-dev gcc gawk make \
 && curl -L -o ./tivodecode.tar.gz http://sourceforge.net/projects/tivodecode/files/tivodecode/0.2pre4/tivodecode-0.2pre4.tar.gz \
 && tar xvfz tivodecode.tar.gz -C /opt/ && rm tivodecode.tar.gz \
 && cd /opt/tivodecode-0.2pre4 \
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

USER kmttg
