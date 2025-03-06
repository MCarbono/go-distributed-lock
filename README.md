<h1 align="center">Distributed lock</h1>

## 📜 Summary
- [About](#About)
- [Docs](#Docs)
- [Run](#Run)
- [Tracing](#Tracing)


<a id="About"></a> 
## 📃 About
This is a project to try and test distributed locks using Redis. There's 2 services(orders and invoices) that have CRUD operations and 
one update or delete request should reflect in both systems avoiding inconsistency. The idea is that for any reason there's a concurrent requests to update/delete these resources, one should not affect the other. See the <a href="#Docs">Docs</a> section
for more details.

<a id="Docs"></a> 
## 🧾 Docs

There's a .json file that you can import on Postman to make requests to the endpoints.
Down below you will see informations about the endpoints. There's mermaids sequenceDiagram to ilustrate all the steps.
You can check the mermaid's doc <a href="https://mermaid.js.org/syntax/sequenceDiagram.html">here</a> and the online editor <a href="https://mermaid.js.org/syntax/sequenceDiagram.html"> here</a> 
The intentions of these sequenceDiagrams are to show the communications between the components and demonstrating the 'happy path' of each flow. For more details, you should look into the code.

<h4>Create order</h4>

```mermaid
sequenceDiagram
  client->>+order_service: POST /orders
  order_service->>+invoice_service: POST /invoices
  invoice_service->>+invoice_database: save invoice
  invoice_database-->>-invoice_service: invoice saved
  invoice_service-->>-order_service: invoice created
  order_service->>+order_database: save order 
  order_database-->>-order_service: order saved
  order_service-->>-client: order created
```
---

<h4>Delete order</h4>

```mermaid
    sequenceDiagram
    client->>+order_service: DELETE /orders/:id
    order_service->>+order_database: find order by ID
    order_database->>-order_service: return order
    order_service->>+lock_database: Try lock for order and invoice by it's respective ID
    lock_database-->>-order_service: locks resource
    order_service->>+invoice_service: DELETE /invoices/:id
    invoice_service->>+lock_database: Checks if invoice is locked
    lock_database-->>-invoice_service: if invoice is locked, continue
    invoice_service->>+invoice_database: deletes invoice
    invoice_database-->>-invoice_service: invoice deleted
    invoice_service-->>-order_service: invoice deleted
    order_service->>+order_database: delete order 
    order_database-->>-order_service: order deleted
    order_service->>+lock_database: remove locks
    lock_database-->>-order_service: locks removed
    order_service-->>-client: order deleted
```
---
<h4>Update order</h4>

```mermaid
sequenceDiagram
  client->>+order_service: PATCH /orders/:id
  order_service->>+order_database: find order by ID
  order_database-->>-order_service: return order
  order_service->>+redis_database: Try lock for order and invoice by it's respective ID
  redis_database-->>-order_service: locks resources
  order_service->>+invoice_service: PATCH /invoices/:id
  invoice_service->>+redis_database: Checks if invoice is locked
  redis_database-->>-invoice_service: if invoice is locked, continue
  invoice_service->>+invoice_database: update invoice
  invoice_database-->>-invoice_service: invoice updated
  invoice_service-->>-order_service: invoice updated
  order_service->>+order_database: update order 
  order_database-->>-order_service: order updated
  order_service->>+redis_database: remove locks
  redis_database-->>-order_service: locks removed
  order_service-->>-client: order updated
```

---

<a id="Run"></a> 
## ⚙️ Run

There's two ways of starting this project: using docker to start the webserver and the infrastructure or starting the infrastructure via docker and
starting the server on your machine.


### Using docker-compose for infra and webserver

Run one of the commands below to build golang image:

```bash
docker-compose -f docker-compose.production.yml build
```

```bash
make build
```

Then run one of the commands below to start the containers:

```bash
docker-compose -f docker-compose.production.yml up -d
```

```bash
make run_prod
```

If you want to destroy it all, run one of the commands below:

```bash
docker-compose -f docker-compose.production.yml down
```

```bash
make infra_down_prod
```

### Using docker-compose for infrastructure and webserver locally

To run the webserver locally, you need to have installed Go on version 1.22.6.<br>

run one of the commands below: 

```bash
docker compose up -d
```

```bash
make infra
```

then start the api:

```bash
go run main.go --env=local
```

```bash
make run
```

If you want to destroy the database's docker, run one of the commands below:

```bash
docker compose down 
```

```bash
make infra_down
```