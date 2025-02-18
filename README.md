# Products API

## Instructions

### Running production
To run the application in production mode, use the following command:

```sh
  ./bin/products_app
```

Note that if you modify the code you will need to recompile the application with the following command:

```sh
  make build
```

### Running the application
To run execute
    
```sh
  make run
```

This will start the Golang application with the current code

### Running Tests
To run the tests for the Golang application, use the following command:

```sh
  make test
```

### Technical decisions

#### Domain as logic center
The domain inside products reside all the code needed to validate product input integrity and calculate product discount.
I did it this way so in the innermost layer of the application we have the business logic and the domain rules detached of any external dependencies such as databases or http requests.

This also includes the repository interface that is used to abstract the database (infra) layer from the domain layer.

#### SQLite
I chose SQLite as the database for this project because it is a simple database that it is easy to setup and use in this case it is entirely self contained
this means that no installation is needed by the user since sqlite is self contained in the library.

Besides offers easily an in memory database that can be used for testing purposes so integrations tests are easy to do without any extra setup.

#### Clean Architecture

A clean architecture was applied here to separate the concerns of the application in layers, the main layers are:
- Domain
- Use cases
- Infra

The main goal with this was abstract every layer from the other so that the application can be easily tested and maintained following the dependency direction rule
thus outermost layer can be tested with or without dependencies of the innermost layer since the details of these dependencies are hidden through interfaces.

#### Streaming of json data at initialization

JsonStream type allows function ReadJson to send a stream of json data to the channel, this is useful to read a large json file without loading it all in memory at once. so if files
grows to 20k rows or any arbitrary number of rows the application will not crash due to memory issues since it will be kept stable, besides thanks to concurrency use, json read and product insertion
can potentially happen simultaneously.

#### init function at Database connection generation
This function allows optionally to execute any function at database initialization, in this case it is used to create the table if it does not exist.

#### Detach domain model from response model

Domain model only contains the fields needed by its own nature of being a product in the another hand the response model contains the fields needed to be returned to the client the more
obvious is price field which contains the price with and discount information, business model offers a method to calculate the discount on the fly so response can contain the price with discount
but how it will be finally represented is http response responsibility and not a domain concerned.