FROM golang:1.8

MAINTAINER Docker SongCF <fuxiao333@qq.com>

EXPOSE 9901
EXPOSE 9911
EXPOSE 9912
EXPOSE 9913

ENV WORKDIR ~


COPY $./scene $WORKDIR/scene
COPY ./conf.ini $WORKDIR/conf.ini

CMD ["./scene"]
