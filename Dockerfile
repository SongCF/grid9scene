FROM golang:1.8

MAINTAINER Docker SongCF <fuxiao333@qq.com>

EXPOSE 9901
EXPOSE 9911
EXPOSE 9912
EXPOSE 9913

ENV WORKDIR ~


COPY $GOPATH/bin/scene $WORKDIR/scene
COPY ./conf.ini $WORKDIR/conf.ini
COPY ./scripts/boot.sh $WORKDIR/boot.sh

CMD ["/bin/bash", "boot.sh"]
