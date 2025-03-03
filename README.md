<h1 align="center">Distributed lock</h1>

## 📜 Summary
- [About](#About)
- [Docs](#Docs)
- [Libs/Dependencies](#Libs/Dependencies)
- [Run](#Run)
- [Endpoints](#Endpoints)
- [Tracing](#Tracing)


<a id="About"></a> 
## 📃 About
This is a project to try and test distributed locks using Redis. There's 2 services(orders and invoices) that have CRUD operations and 
one update or delete request should reflect in both systems avoiding inconsistency. See the <a href="#Endpoints">endpoints</a> section
for more details.

<a id="Endpoints"></a> 
## 💻 Endpoints

In this section you will see informations about the endpoints. There's mermaids sequenceDiagram to ilustrate all the steps.
You can check the mermaid's doc <a href="https://mermaid.js.org/syntax/sequenceDiagram.html">here</a> and the online editor <a href="https://mermaid.js.org/syntax/sequenceDiagram.html"> here</a> 

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