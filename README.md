# jenkins-docker-api

## Endpoints

A list of endpoints this API supports

```sh
GET   /api/v1/health # Is the API up

GET   /api/v1/jenkins # List of running Jenkins masters

PUT   /api/v1/jenkins/restart/:name # Restarts a single Jenkins master

PUT   /api/v1/jenkins/update/:name # Update a single Jenkins master to the latest release

PUT   /api/v1/jenkins/shutdown/:name # Shutdown a single Jenkins master

PUT   /api/v1/jenkins/admin/update_all # Update all Jenkins masters to the latest release (admins only)
```

### Relevant Header

`-H "Authorization: token <your_token>"`

## Development

Dependencies:

1. Ensure that Docker is installed and running.
1. Ensure that Docker swarm is initialized: ```docker swarm init```.
1. Ensure that [golang](https://golang.org/dl/) is installed.
1. Ensure that [govendor](https://github.com/kardianos/govendor) is installed.

Setting up project:

1. Clone down the project.

    ```console
    cd $GOPATH/src/github.com/target/
    git clone git@github.com:Jenkins/jenkins-docker-api.git
    cd jenkins-docker-api
    ```
1. Make your code changes.
1. Update the version in `version/version.go` if applicable
1. Commit and push your changes.
