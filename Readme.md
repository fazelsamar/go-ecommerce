# Go Eccomerce

## Description
Describe your project here, including its purpose, functionality, and any key features.

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

Clone the repository:
```shell
git clone https://github.com/fazelsamar/go-ecommerce.git
```

To run this application you have two aptions:

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
## Usage
Explain how to use your application here. Provide any necessary instructions or examples.

## Project Structure
Provide a high-level overview of your project's file structure, including the purpose of key files and directories.

## Contributing
Provide guidelines for contributing to your project, if applicable.

## License
Specify the license under which your project is distributed.

## Contact
Provide contact information for users to reach out for support or inquiries.