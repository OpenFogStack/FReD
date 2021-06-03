FROM node:10-alpine

#RUN apk add python make gcc g++ 
RUN npm config set user root
RUN npm i grpcc -g

ENTRYPOINT [ "grpcc" ]
