FROM alpine

MAINTAINER Tobias Pfandzelter

ADD frednode frednode

EXPOSE 9001

RUN ls

ENTRYPOINT ["./frednode", "--addr localhost:9001"]