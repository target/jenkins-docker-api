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

`-H "Authorization: token <github token>"`

## Contributing / Development

Dependencies:

1. Ensure that Docker is installed _and_ running.
1. Ensure that Docker swarm is initialized: ```docker swarm init```.
1. Ensure that [golang](https://golang.org/dl/) is installed.
1. Ensure that [govendor](https://github.com/kardianos/govendor) is installed.

Setting up project:

1. Clone down the project:

    ```sh
     # Make sure your go paths are set if they aren't already
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    
     # Clone the project
    cd $GOPATH/src/github.com/target/
    git clone git@github.com:Jenkins/jenkins-docker-api.git
    cd jenkins-docker-api
    ```
    
1. Add/update golang vendor packages:

    ```sh
    govendor fetch +outside
    ```
    
1. Set environment variables:

    - `GITHUB_API_URL` - your GitHub API URL
    - `GITHUB_ADMIN_ORG` - name of your GitHub admin organization, e.g. `Jenkins`
    - `GITHUB_ADMIN_TEAM` - name of your GitHub admin team, e.g. `Admins`
    - `JENKINS_ENV` - your Jenkins environment, e.g. `test`, `prod`, etc.
    - `JENKINS_IMAGE` - your Jenkins Docker image
    - `JENKINS_USER_CONFIG_PATH` - default path is `/jenkins/user-configs/`
    - `JENKINS_ADMIN_CONFIG_PATH` - default path is `/jenkins/secret-configs/`
    
1. Make your code changes and ensure all tests pass

    ```sh
    # Checkout a branch for your work
    git checkout -b name_of_your_branch

    # Code away!
    ```
    
1. Update the version in `version/version.go` if applicable

1. Submit a PR for your changes.
