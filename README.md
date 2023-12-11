### Run Keycloak Locally ###
```
docker run -d --name keycloak \
    -e KEYCLOAK_DATABASE_VENDOR=dev-file \
    -e KEYCLOAK_ADMIN_USER=admin \
    -e KEYCLOAK_ADMIN_PASSWORD=admin \
    -p 8080:8080 \
    bitnami/keycloak:22.0.5
```
- After running the Keycloak,
  - Add Google as Identity Provider. Check the link below as tutorial. \
    https://medium.com/codeshakeio/configure-keycloak-to-use-google-as-an-idp-4e3c59583345
  - Create a Client in Keycloak and save the client id and client secret for later use in running this backend service. 
  
### Run MongoDB Locally ###
```
docker volume create mongo-data
docker run -d --name mongodb \
    -v mongo-data:/etc/mongo \
    -p 27017:27017 \
    mongo:5.0.22
```

### Environment Config ###
- To run this backend locally, set the environment values in the `.env` file
- To set the production config, set the environment variables in `k8/secret.yaml` file


### APIs ###
- Postman APIs example in available in the `Alap.postman_collection.json` file. You can import the file in your postman and check.

### Websocket Endpoint ###
- /ws?accessToken=[ACCESS_TOKEN]


### WS Message Send Format ###
```
{
    "to": "RECIPIENT_USER_ID",
    "chatWindow": "CHAT_WINDOW_ID_OF_THE_PARTICIPANTS",
    "msg": "TEXT_MESSAGE"
}
```


### WS Message Receive Format ###
```
{
    "sender": "SENDER_USER_ID",
    "chatWindow": "CHAT_WINDOW_ID_OF_THE_PARTICIPANTS",
    "msg": "TEXT_MESSAGE"
}
```

### Kubernetes Deployment ###
All Kubernetes deployment descriptors are available in the `k8` folder