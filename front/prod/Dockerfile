FROM node:10.9-slim AS build

RUN printf "deb http://archive.debian.org/debian/ jessie main\ndeb-src http://archive.debian.org/debian/ jessie main\ndeb http://security.debian.org jessie/updates main\ndeb-src http://security.debian.org jessie/updates main" > /etc/apt/sources.list
COPY ./ /app
WORKDIR /app

RUN yarn install && yarn build

FROM nginx
RUN mkdir /app
COPY --from=build /app/dist /app
COPY ngnix.conf /etc/nginx/nginx.conf
EXPOSE 80
