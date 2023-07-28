# Go Eccomerce
This Go application is a backend implementation for an eccomerce application. It provides a range of features such as user authentication, products management, creating carts and orders. 

## Technologies Used
- Go
- Fiber
- PostgreSQL
- Docker
- nginx
- pgAdmin

## Installation
### Prerequisites
- Docker
- Docker Compose
- Go

Clone the repository:
```shell
git clone https://github.com/fazelsamar/go-ecommerce.git
```

To run this application you have two options:

1. First is to run it using docker-compose for production :

    ```shell
    docker-compose -f docker-compose-prod.yml up -d
    ```
2. Second is to run just the dependecy in docker-compose for development :

	```shell
    docker-compose up -d
    ```
    And then run the go:

    ```shell
    go run ./cmd
    ```

## Features

### This application provides a range of features including:

- #### User authentication: users can sign up, log in and log out.
- #### Create and delete products by admin.
- #### Create and delete cart by anonymous user.
- #### Add products to the cart.
- #### Creating the order base on the carts.

## Code Structure

- `cmd`: contains the main go entry file.
- `internal/`
    - `entity`: contains database models.
    - `handlers`: contains http handlers for each of the supported endpoints.
    - `middleware`: contains custom middleware used in the application.
    - `repositories`: contains implementations of the persistence interface to interact with the database.
    - `routes`: contains the main application logic for the backend.
    - `services`: contains the business logic and use case implementation of the application.
- `pkg`: contains the packages such as database connection and managing env files.
- `proxy`: contains the nginx config file.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

<p style="text-align: center; width: 100%; ">Copyright&copy; 2023 <a href="https://github.com/fazelsamar">Fazelsamar</a>, Licensed under MIT</p>