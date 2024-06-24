### Setup
go version 1.22
#### Clone the Repository
```bash
git clone https://github.com/vlhvetrava/dataflow.git
cd dataflow
```
#### Install Dependencies
```bash
go mod tidy
```
#### Run Tests
```bash
go test ./...
```
#### Run the Application
```bash
go run main.go
```

### Architectural remarks
1. Layered project structure is used, with separate handlers, services and repository levels.
Service layer contains business logic, making it reusable and easier to test independently of the HTTP layer.
Separating the repo level and using interface there allows us to abstract the data storage mechanism, 
making the code more modular and easier to test. It also allows for flexibility in changing the storage implementation 
without affecting the business logic.
2. Gin framework is used for its performance and simplicity in handling HTTP requests.
3. For calculating total sales we're using `math/big` package to better handle operations with float numbers and precision.
4. `sync.Map` is used for safe concurrent access to the in-memory data store. This makes the solution scalable and efficient 
for handling simultaneous requests in a multithreaded environment.
5. Error handling ensures that the api will respons with proper status codes and messages.