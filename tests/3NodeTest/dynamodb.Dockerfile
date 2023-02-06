FROM amazon/dynamodb-local:1.21.0

WORKDIR /home/dynamodblocal

# add prepared database file to local dynamodb image
RUN mkdir data
COPY --chown=dynamodblocal:dynamodblocal local-fred-dynamo.db data/shared-local-instance.db
RUN chmod -R 777 data

EXPOSE 8000

ENTRYPOINT [ "java", "-jar", "DynamoDBLocal.jar", "-sharedDb", "-dbPath", "./data" ]