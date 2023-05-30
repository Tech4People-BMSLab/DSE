FROM node:18-alpine

WORKDIR /api
COPY ./api/package*.json ./

RUN npm i -g npm@latest
RUN npm i -g pnpm

COPY ./api/src            ./src
COPY ./api/nodemon.json   ./nodemon.json
COPY ./api/package.json   ./package.json
COPY ./api/tsconfig.json  ./tsconfig.json
COPY ./api/*.lock         ./

RUN pnpm i

# Start the application using ts-node
CMD ["npx", "nodemon"]
EXPOSE 5004