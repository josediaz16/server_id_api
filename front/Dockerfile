FROM node:10.9-slim

RUN printf "deb http://archive.debian.org/debian/ jessie main\ndeb-src http://archive.debian.org/debian/ jessie main\ndeb http://security.debian.org jessie/updates main\ndeb-src http://security.debian.org jessie/updates main" > /etc/apt/sources.list

RUN apt-get update \
  && apt-get install -y git

RUN yarn global add @vue/cli -g

WORKDIR /app

RUN apt-get autoremove -y \
  && apt-get autoclean -y \
  && apt-get clean -y \
  && rm -rf /var/lib/apt/lists/*

EXPOSE 5000

CMD ["yarn", "serve", "--port", "5000"]
