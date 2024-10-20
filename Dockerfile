## copied/modified from: https://github.com/drnay/docker-kmttg/blob/master/Dockerfile
# Alpine Linux with Oracle JRE
FROM sgrio/java:jre_8

RUN apk update
RUN apk upgrade
RUN apk add --no-cache \
      font-noto \
      gettext \
      gtk+2.0 \
      jq \
      libxtst \
      tzdata \
      wget \
      xterm

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

# server port
EXPOSE 8181

# Setup environment variables
ENV LD_LIBRARY_PATH=/lib:/usr/lib

# Run as a non-root (appending to /etc files directly to override 
# potential pre-existing users/groups with the given IDs)
ARG USER_ID
ARG GROUP_ID
RUN echo "kmttg:x:${GROUP_ID}:kmttg" >> /etc/group \
 && echo "kmttg:x:${USER_ID}:${GROUP_ID}:kmttg:/home/kmttg:/sbin/nologin" >> /etc/passwd \
 && mkdir -p /home/kmttg \
 && chown kmttg /home/kmttg

USER kmttg

# Set the working directory
ENV APP_DIR /home/kmttg/app
RUN mkdir -p ${APP_DIR}/web/cache
WORKDIR ${APP_DIR}

# Get the latest kmttg version
ARG APP_VERSION
RUN url=$(curl -L -H "Accept: application/vnd.github+json" -H "X-GitHub-Api-Version: 2022-11-28" "https://api.github.com/repos/lart2150/kmttg/releases/tags/${APP_VERSION}" \
  | jq -r '.assets[0].browser_download_url') \
 && wget -O kmttg.zip "${url}" \
 && busybox unzip -o kmttg.zip \
 && chmod +x kmttg

ENV INPUT_DIR /mnt/input
ENV OUTPUT_DIR /mnt/output
ENV TOOLS_DIR /usr/local/bin 

# mount points for input/output files
VOLUME ${OUTPUT_DIR}

# persist the install dir
VOLUME ${APP_DIR}

COPY --chmod=777 kmttg.sh .

CMD ["/home/kmttg/app/kmttg.sh"]
