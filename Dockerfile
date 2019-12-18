FROM scratch

MAINTAINER Tobias Pfandzelter

ADD frednode .

EXPOSE 9001

CMD ["./frednode", "--addr localhost:9001"]