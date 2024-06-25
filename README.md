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
3. For calculating total sales `math/big` package is used to better handle operations with float numbers and precision.
4. `sync.Map` is used for safe concurrent access to the in-memory data store. This makes the solution scalable and efficient 
for handling simultaneous requests in a multithreaded environment.
5. Error handling ensures that the api will respond with proper status codes and messages.

### Use Cases

#### Get All Sales
Fetch all sales records from the database.

#### GET /data

**Example Request:**
```sh
curl -X GET http://localhost:8080/data
```
**Example Response:**

```bash
[
    {
        "id": "1",
        "product_id": "12345",
        "store_id": "6789",
        "quantity_sold": 10,
        "sale_price": 19.99,
        "sale_date": "2024-06-15T14:30:00Z"
    },
    {
        "id": "2",
        "product_id": "54321",
        "store_id": "9876",
        "quantity_sold": 5,
        "sale_price": 9.99,
        "sale_date": "2024-06-16T10:00:00Z"
    }
]
```

#### Add a Sale
Add a new sale record to the database.

#### POST /data

**Example Request:**
```sh
curl -X POST http://localhost:8080/data \
     -H "Content-Type: application/json" \
     -d '{
           "product_id": "12345",
           "store_id": "6789",
           "quantity_sold": 10,
           "sale_price": 19.99,
           "sale_date": "2024-06-15T14:30:00Z"
         }'
```

**Example Response:**

```bash
{
"status": "success"
}
```

#### Calculate Sales
Calculate total sales for a specific store within a given date range. If range is empty, return all sales for a provided storeId.

#### POST /calculate


**Example Request:**
```sh
curl -X POST http://localhost:8080/calculate \
     -H "Content-Type: application/json" \
     -d '{
           "operation": "total_sales",
           "store_id": "6789",
           "start_date": "2024-06-01T00:00:00Z",
           "end_date": "2024-06-16T00:00:00Z"
         }'
```
**Example Response:**

```bash
{
    "store_id": "6789",
    "start_date": "2024-06-01T00:00:00Z",
    "end_date": "2024-06-16T00:00:00Z",
    "total_sales": 199.9
}
```

**Example Request (without dates, returns all sales for store_id):**

```sh
curl -X POST http://localhost:8080/calculate \
     -H "Content-Type: application/json" \
     -d '{
           "operation": "total_sales",
           "store_id": "6789",
           "start_date": "",
           "end_date": "",
           
         }'
```
**Example Response:**
```bash
{
    "store_id": "6789",
    "start_date": "",
    "end_date": "",
    "total_sales": 199.9
}
```

