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

CMD ["./scene"]
