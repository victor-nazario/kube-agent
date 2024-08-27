# kube-agent

## Building and Running the Application
The recommended way to build and run this application is using [Docker](https://www.docker.com/), this application
provides a ``Dockerfile`` with all the instructions needed to build a working image of the application and its
dependencies.

On a machine with Docker installed, cd into the application root and execute ``docker build --tag agent .``
this command will produce a working image of the application, depending on your network connection it will
take some time. After this command executes, you can create a container with the image in detached mode
issuing the following command ``docker run -d -p 8080:8080 agent``. Once this command executes, the user can
visit [http://localhost:8080/graphql](http://localhost:8080/graphql) on their machine and interact with the GraphiQL
client.

## Architecture

![agent-arch](https://github.com/user-attachments/assets/588ee7ec-b7ca-443a-a5a4-a7b8b7eef228)


[Agent] - The agent is two components embedded into one. Firstly, we provide a graphql server as an entry-point
to Controllers seeking to schedule and deploy Kubernetes jobs. Secondly, it provides a Kubernetes Controller which
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

# Reliability
The designed system provides resilience and reliability by exposing itself via a service to be managed by a cluster ingress.
The deployment definition would require 3 pods to be present, depending on load. Whilst out of the scope of the project,
a production level implementation could consider (if required by the use case) an HPA policy. This provides
scaling on demand, whilst also reducing the costs associated with the deployment, both economical and in terms of resources.

