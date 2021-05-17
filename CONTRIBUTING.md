# Contributing to The-Amazing-Ledger
We love your input! We want to make contributing to this project as easy and transparent as possible, whether it's:

* Reporting a bug
* Discussing the current state of the code
* Submitting a fix
* Proposing new features
* Becoming a maintainer

## An Guideline to Write a Feature

- Download dependencies
    - run `make setup` to get all required softwares.
- Add entities or vos
    - add or edit `app/domain/entities/` and `app/domain/vos/`.
- Repository interface
    - edit `app/domain/repository.go` to add your method to interface
    - add file in `app/gateways/db/` implementing this method.
- Ledger Usecase interface
    - edit `app/domain/usecase.go` to add yout method to interface
    - add file in `app/domain/usecases/` implementing this method
- Work with protofile
    - change the file `proto/ledger/ledger.proto` adding your services.
    - generate proto server with `$ make generate` shell command.
- Write SDK
    - add file in `clients/grpc/ledger/` with your sdk feature.
    - add file in `clients/grpc/examples/` to test your feature.
    - call your test feature in `clients/grpc/examples/main.go`.

_write a grpc examples are optional, we use like end-to-end tests_

### Running end-to-end tests

The amazing ledger needs a Postgres database running and ready to connections,
see `docker-compose` files to run a database, or use another way like install
postgres in your machine.

After up the database, you need compile the code with `$ make compile`, so a
binary file will be created at `build/` folder.

Now, you can run `the amazing ledger` with `./build/server` command.

Finally you need run `clients/grpc/examples` with `go run ./...` command.

Let us know if you have any questions, use [gitter][gitter] to send us a message

## We use [Github Flow](https://guides.github.com/introduction/flow/index.html), so all code changes happen through Pull Requests
Pull requests are the best way to propose changes to the codebase (we use [Github Flow](https://guides.github.com/introduction/flow/index.html)). We actively welcome your pull requests:

1.  Fork the repo and create your branch from `main`.
2.  If you've added code that should be tested, add tests.
3.  If you've changed APIs, update the documentation.
4.  Ensure the test suite passes.
5.  Make sure your code lints.
6.  Issue that pull request!

We encourage the use of [semantic commit
messages](https://seesparkbox.com/foundry/semantic_commit_messages) for better understanding of what
is being done in each commit.

We use a changelog based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).
We add new entries to it in the commit that makes the changes.

Our versioning follows [Semantic Versioning](http://semver.org).

For more best practices details, read
[stone-best-practices](https://github.com/stone-payments/stoneco-best-practices/).

## Report bugs using Github's [issues](https://github.com/stone-co/the-amazing-ledger/issues/new/choose)
All of our issues have a template with tips for describing the problem. Try to follow the proposed structure, this helps the reviewer a lot.

## License
By contributing, you agree that your contributions will be licensed under its [MIT license](LICENSE).

[gitter]:https://gitter.im/the-amazing-ledger/community#
