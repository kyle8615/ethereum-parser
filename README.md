# Ethereum Parser

## Overview
Ethereum Parser is a service that allows users to subscribe to Ethereum addresses and check incoming and outgoing transactions by addresses. It uses Ethereum's JSON RPC API to interact with the blockchain and provides a simple REST API for users to interact with.

## Getting Started
### Prerequisites
Go (version 1.15 or later)
Access to the internet to interact with the Ethereum JSON RPC API
### Installation
Clone the repository:

```
git clone https://github.com/kyle8615/ethereum-parser.git
cd ethereum-parser
```
### Running the Service
To start the Ethereum Parser service, navigate to the root directory of the project and run:

```
go run ./cmd
```
The service will start and listen on http://localhost:8080.

## API Usage
### Subscribe to an Address
**Endpoint:** POST /subscribe

**Description:** Subscribe to an Ethereum address to monitor its transactions.

**Payload:**

```
{
    "address": "0xYourEthereumAddress"
}
```
**cURL Example:**

```
curl -X POST 'http://localhost:8080/subscribe' \
-H 'Content-Type: application/json' \
-d '{"address": "0xYourEthereumAddress"}'
```

### Get Transactions for an Address
**Endpoint:** GET /transactions?address=0xYourEthereumAddress

**Description:** Retrieve the list of transactions for a subscribed Ethereum address.

**cURL Example:**

```
curl 'http://localhost:8080/transactions?address=0xYourEthereumAddress'
```

### Get Current Block
**Endpoint:** GET /currentblock

**Description:** Get the number of the most recent block on the Ethereum blockchain.

**cURL Example:**

```
curl 'http://localhost:8080/currentblock'
```

Contributing
Contributions to the Ethereum Parser project are welcome! Please submit a pull request or issue on the project's GitHub page.

License
MIT License
