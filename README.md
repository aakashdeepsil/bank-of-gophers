# Bank of Gophers

Bank of Gophers is a simple banking application built with Go.

## Features

- Account creation
- Deposit and withdrawal
- Balance inquiry
- Transaction history

## Installation

1. Clone the repository:
2. Run a docker container using this command - `docker run --name bank-of-gophers-postgres -e POSTGRES_PASSWORD=bank-of-gophers -p 5432:5432 -d postgres`
3. Run the Makefile - `make run`
4. Open Postman and perform the operations on the endpoints for development and testing purpose.