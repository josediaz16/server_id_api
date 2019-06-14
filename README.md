# README

# Server ID api

Server ID api is a service that allows you to check different data about domains, including Ssl grade, owners and servers country.

## Site
Check out the app here [Server ID Api](http://34.220.54.27)

Interact with the api at http://34.220.54.27/api

Available Endpoints
  - GET /domains   **Get a list of domains**
  - GET /domains/search?domainName=somedomain.com   **Get the current status of a domain**
 
## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

  - Golang v1.12.5
  * Using Docker for Development
  - Docker version 17.03.0-ce or higher
  - Docker Compose version 1.21.2 or higher

### Installing

A step by step series of examples that tell you how to get a development env running

  1. Run `sudo docker-compose run db_init` to create the database.`
  3. Run `sudo docker-compose run -e DATABASE=servers_test db_init /setup_db.sh` to create the test database.

### Running the app

  - Run app in development mode `sudo docker-compose up front`
  
### Running the app in prod mode

  1. Create a .env file at project root with DATABASE_URL variable.
  2. Run app locally but with production configuration `sudo docker-compose -f docker-compose.prod.yml up front`

## Test Suite

  - Run `sudo docker-compose run api ./run_tests.sh` for Golang api tests.


## Build prod images
  - Login to Dockerhub
  ```
  sudo docker build -f front/prod/Dockerfile -t jldiazb16/server_id_web front/
  sudo docker push jldiazb16/server_id_web
  ```
  to build and push vue web app.
  ```
  sudo docker build -f api/docker/prod/Dockerfile -t jldiazb16/server_id_api api/
  sudo docker push jldiazb16/server_id_api
  ```
  to build and push golang api.
     
## Deployment instructions

This repository use [Capistrano](https://capistranorb.com/) gem for deployment. Please read the documentation first.
  - Run `deploy/full_up` to deploy the app to AWS.
  
Deployed using Docker containers on Prod.
  
## CD/CI

This app CD/CI has been configured with Semaphoreci v2.0. Config files are available at .semaphore/ folder
in project root. Check out the docs at [Semaphoreci](https://docs.semaphoreci.com/category/48-reference)

### Built With

* [Golang](https://golang.org/) - Golang Language
* [Chi-Router](https://github.com/go-chi/chi) - Chi Router
* [CockroachDB](https://www.cockroachlabs.com/docs/stable/) - Cockroach DB
* [VueJS](https://vuejs.org/) - Vue JS

## TODO
- Generate API documentation
- Add integration tests
- Deploy the app in a multi-host environment (Kubernetes Maybe)

## Contributing

This is a private repository.

## Authors

* **Jose Diaz** - [josediaz16](https://github.com/josediaz16)
