# kube-agent

## Building, Running and Deploying the Applications
The recommended way to build and run this application is using [Docker](https://www.docker.com/), this application
provides a ``Dockerfile`` with all the instructions needed to build a working image of the application and its
dependencies.

On a machine with Docker installed, cd into the application root and execute ``docker build --tag agent .``
this command will produce a working image of the application, depending on your network connection it will
take some time. After this command executes, you can create a container with the image in detached mode
issuing the following command ``docker run -d -p 8080:8080 agent``. Once this command executes, the user can
visit [http://localhost:8080/query](http://localhost:8080/query) on their machine and interact with the GraphiQL
client.

Once the image is built, please execute the following command (this assumes the user has installed Minikube) 
```eval $(minikube docker-env --unset)```, this will allow us to use the local image rather than try to pull from 
a remote registry. 

Assuming the user is at the project's root, please execute the following command to create the service-account for 
the system, this allows us to define roles and improve security by using the least privilege possible. 
```kubectl apply -f deploy/service-account.yaml```

Now, we need to create the deployment and service to access our application (this is merged into one file).
```kubectl apply -f deploy/deployment.yaml```

The resources should have been created successfully, this means that we now have routing from outside the cluster to the 
application. One way to obtain the specific node-ip and port where we can reach it:
```minikube service kube-agent-entrypoint --url -n operant```

Now, we can build the Operant Controller CLI. Assuming the user is at the controller directory, execute the following:
```go build -o ./out/kube-agent .```

Provide run permissions to the binary
```chmod +x ./out/kube-agent```

With the Agent running and correctly exposed to outside the cluster, you can issue commands with the
recently built binary as such:

```./out/controller --u=operant --p=secret --jobname=nazario-test --image=ubuntu:latest --command=ls --backoff=0 --url={$YOUR_IP:PORT}/query --namespace=operant```


## Architecture

![agent-arch](https://github.com/user-attachments/assets/588ee7ec-b7ca-443a-a5a4-a7b8b7eef228)


[Agent] - The agent is two components embedded into one. Firstly, we provide a graphql server as an entry-point
to Controllers seeking to schedule and deploy Kubernetes jobs. This means, we have significant
flexibility into what a Controller can be, it can even be graphical. In our case, we have decided to implement 
it as a CLI application. Secondly, it provides a Kubernetes Controller which
leverages the existing V1 APIs to create the resources needed and manage their lifecycle. 

[Controller] - The Controller (not to be confused with a Kubernetes Controller) is a CLI client application
which accepts relevant details about a job and creates the resource in the cluster.

## Security 
While a production grade implementation of the security components of this application is out of the scope of
the project, it mimics the security features one would expect from a zero thrust architected application in the real-world.

* Authentication: This application requires and enforces requests to be authenticated. While the implementation is Basic Auth,
we demonstrate how an Authentication middleware should operate and handle every incoming request, not thrusting authentication
performed outside the system, such as one done by a GraphQL Gateway like Apollo.
* Role Based Access Control (RBAC): In a production system, roles are essential to establishing clear boundaries between 
user access to resources. This application provides roles and actions along with a handler that acts as a middleware to
validate the current action is accessible to the given user.
* Limited role access with Kubernetes Service Accounts: This application implements the least amount of privilege possible,
we use a custom service account spec which gives enough permissions to create the kubernetes job, but not to operate over other resources.

# Reliability
The designed system provides resilience and reliability by exposing itself via a service to be managed by a cluster ingress.
The deployment definition would require 3 pods to be present, depending on load. Whilst out of the scope of the project,
a production level implementation could consider (if required by the use case) an HPA policy. This provides
scaling on demand, whilst also reducing the costs associated with the deployment, both economical and in terms of resources.

# Testing
This application provides a test suite that validates and tests the most critical parts of the system,
tests can be executed with the ```go test``` command.

# Copyright Notice
Copyright 2024 Victor Nazario.


