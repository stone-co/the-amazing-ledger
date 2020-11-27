<h1 align="center">
  <br>
  The Amazing Ledger
</h1>

The Amazing Ledger is a double-entry accounting system used to record transactions. The entries are processed, validated and persisted to a database and their original data is **never** modified.

# Concepts

The Amazing Ledger records a journal that keeps track of all credits and debits to and from accounts (an account being anywhere money can go, not necessarily a bank account).

Since it is a double-entry accounting system, there must always be a debit from at least one account for every credit made to another account, therefore all **entries** must balance to zero within a **transaction**.

For example,

| Account          | Amount |
|------------------|--------|
| Satoshi Nakamoto | $ -50  |
| Hal Finney       | $ 50   |

The system supports the following account types:
- **equity**: Represents your net worth. This will be used to start the ledger and its values are naturally negative.
- **asset**: Represents the money you have, where money sits. The values of theses accounts are naturally positive.
- **liability**: Represents the money you owe and is a synonym to debt. The values of theses accounts are naturally negative.
- **expense**: Represents the money you spent, where money goes. The values of theses accounts are naturally positive.
- **income**: Represents the money you have earned, where money comes from. The values of theses accounts are naturally negative.

# Development

To init development environment and run the amazing ledger, flow these steps.

Database setup

```bash
$ docker-compose -f docker-compose-dev.yml up
```

Wait for database setup and run in another shell this commands.

```bash
$ make compile
$ ./build/server
```

## Setup

Download and install dependency tools

```bash
$ make setup
```

## Running tests

```bash
$ make test
```

# Usage

A Web API is used to issue commands and query for data in the ledger.

```bash
curl -i -X POST localhost:3000/api/v1/transactions -d \
'{"id":"28c547fa-4dd4-2593-945c-495678d7a123", "entries":[{"id":
"16b23084-686b-434a-8323-db483ce1e584", "operation":"OPERATION_DEBIT",
"account_id":"liability:clients:available:16b23084-686b-434a-8323-db483ce1e589", "version": 0, "amount":
123},{"id": "16b23084-686b-434a-8323-db483ce1e586", "operation":"OPERATION_CREDIT",
"account_id":"liability:clients:available:16b23084-686b-434a-8323-db483ce1e581", "version": 0, "amount":
123}]}'
```
