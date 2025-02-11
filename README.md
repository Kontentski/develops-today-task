#### Architecture

**High-Level File Structure**

- `./cmd`: Contains the entry point for the program, allowing users to run and use its features.
- `./config`: Contains configuration files that define the behavior of the program and allow users to customize it to their needs.
- `./internal`: Contains the program's internal logic, the code that dictates how the software works and operates.
- `./pkg`: Incorporates the external logic of the program that can be utilized by other applications to add new features or refine existing ones.

**Internal Architecture**
The internal components of the program are defined by a clean architecture that establishes a high-level hierarchy.

Clean architecture is a software design that helps maintain an organized codebase and promotes development flexibility. For a better understanding of the architecture, it's recommended to read "Clean Architecture" by Robert Martin.

**File Structure of the `internal` Directory**:

- `./app`: Responsible for the program's initialization, and contains code related to the initialization of other levels.
- `./entity`: Represents the core business domain, reflecting the entities of the application and operations on them. It also includes the most critical business rules [high-level module].
- `./service`: Represents the business logic layer, responsible for executing business rules, such as creating spy cats. [mid-level module].
- `./api`: This layer is responsible for communication with external services [low-level module].
- `./storage`: The storage layer, responsible for performing CRUD operations on the data [low-level module].
- `./controller`: The transport layer, responsible for obtaining input data, passing it to the business logic level, and returning the output result. Here you can find directories like `http`, `pubsub`, and `queue`, as these are all data exchange mechanisms [low-level module].

![](https://blog.cleancoder.com/uncle-bob/images/2012-08-13-the-clean-architecture/CleanArchitecture.jpg)

## Run locally

#### Clone project from Github

###### https:

```sh
git clone github.com/Kontentski/develops-today-task
```

To run all necessary services, ensure that Docker is started, then execute the following command in the root folder:

```
docker-compose --env-file config/.env up --build postgresdb api
```

For stop containers and removing all images run command:

```
docker-compose down --rmi local
```

#### Testing

Postman collection link:

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://god.gw.postman.com/run-collection/34376513-261257ef-2adc-4811-bd15-53a14f18ffab?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D34376513-261257ef-2adc-4811-bd15-53a14f18ffab%26entityType%3Dcollection%26workspaceId%3Dbd0d64bc-f498-4e4f-84f6-537e25d829f1)
