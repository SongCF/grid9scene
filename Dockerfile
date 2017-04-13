FROM 139.198.2.55/soalib/golang:1.8

MAINTAINER Docker SongCF <fuxiao333@qq.com>

EXPOSE 9901
EXPOSE 9911
EXPOSE 9912
EXPOSE 9913

ENV BIN /go/bin
WORKDIR $BIN

COPY ./scene $BIN/scene
COPY ./conf.ini $BIN/conf.ini
COPY ./scripts/docker_boot.sh $BIN/docker_boot.sh

CMD ["/bin/bash", "./docker_boot.sh"]
