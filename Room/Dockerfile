FROM node:10

ENV token_room = 1
WORKDIR /code

USER root

COPY package*.json ./

RUN npm install

COPY . .

EXPOSE 8080

CMD [ "node", "server.js" ]