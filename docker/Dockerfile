FROM erlang:latest

MAINTAINER Docker SongCF <fuxiao333@qq.com>

EXPOSE 9901
EXPOSE 9911

ENV SOURCE_DIR /jh/source/scene/
WORKDIR $SOURCE_DIR
COPY ./scene $SOURCE_DIR
RUN ["rm", "-rf", "log"]
RUN ["make", "clean"]

RUN ["chmod", "+x", "rebar"]
RUN ["make"]
RUN ["make", "g"]
RUN ["make"]

CMD ["/bin/bash", "./docker/start.sh"]
