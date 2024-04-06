# isol8r-backend

## Description

isol8r enables deployment of feature changes separate from dev, UAT, prod. This means that long running projects, expriments can be tested prior to merging the code.

## Installation

To setup this code locally you can follow the follwing steps:

1. Clone the repo

```go
git clone git@github.com:Swechhya/isol8r-backend.git
```

2. Install all the go modules

```go
go get .
```

3. Copy the `.env.example` file and rename it as `.env`. Add add all the env config.

4. To run the program execute the following code:

```go
 go run cmd/app/main.go
```

## Features

- Developers can test their changes in a controlled environment without affecting other ongoing development work or the stability of existing features
- By testing changes in an isolated environment, teams can identify and address potential issues or bugs early in the development process.
- Developers can easily rollback or roll-forward changes within the isolated environment without affecting the main development or production environments
- Developers can allocate resources specifically for testing and experimentation, ensuring efficient use of computing power
- Designed to be scalable and adaptable to the needs of different projects and teams

## Future Enhancements

- Support database other than postgres
- Support single repo that contain both client and server code
- Support other access control platforms like Gitlab, BitBucket etc.
- Add support of secret manager in feature environment
