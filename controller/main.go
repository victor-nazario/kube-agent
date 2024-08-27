package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const defaultApiUrl = "http://127.0.0.1:8080/query"

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data struct {
		DeployJob string `json:"DeployJob"`
	} `json:"data"`
}

func (g GraphQLResponse) String() string {
	if len(g.Errors) != 0 {
		sb := strings.Builder{}
		defer sb.Reset()

		for _, e := range g.Errors {
			sb.WriteString(fmt.Sprintf("Error in mutation %s: %s\n", e.Path, e.Message))

			// here we are not fortunate enough to have typed errors, best we can do is find the suffix -- this is not ideal
			if strings.HasSuffix(e.Message, "already exists") {
				sb.WriteString("HINT: Provide a different name in --jobName\n")
			}
		}
		return sb.String()
	} else {
		return fmt.Sprintf("Success!\nResource ID: %s", g.Data.DeployJob)
	}
}

type DeployJobInput struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Image        string `json:"image"`
	Command      string `json:"command"`
	BackOffLimit int    `json:"backOffLimit"`
}

func main() {

	// this is just some ASCII art to customize the program execution
	fmt.Println("\n┏┓             ┏┓         ┓┓    \n┃┃┏┓┏┓┏┓┏┓┏┓╋  ┃ ┏┓┏┓╋┏┓┏┓┃┃┏┓┏┓\n┗┛┣┛┗ ┛ ┗┻┛┗┗  ┗┛┗┛┛┗┗┛ ┗┛┗┗┗ ┛ \n  ┛                             ")

	user := flag.String("u", "operant", "username")
	password := flag.String("p", "secret", "password")
	url := flag.String("url", defaultApiUrl, "api url")
	jobName := flag.String("jobname", "test-job-1", "Name of the Kubernetes job to create")
	namespace := flag.String("namespace", "operant", "Namespace where the job should be created")
	containerImage := flag.String("image", "ubuntu:latest", "Name of the container image")
	entryCommand := flag.String("command", "ls", "Command to run inside the container")
	backOffLimit := flag.Int("backoff", 0, "Backoff limit for the job")

	flag.Parse()

	fmt.Printf("\nScheduling job. Job details: \n Name: %s\n Namespace: %s\n Image: %s\n Command:%s\n BackOffLimit:%d\n",
		*jobName, *namespace, *containerImage, *entryCommand, *backOffLimit)

	rea, err := buildMutationReader(*jobName, *namespace, *containerImage, *entryCommand, *backOffLimit)
	if err != nil {
		log.Fatalf("Error marshaling request: %v", err)
	}

	res, err := performGraphRequest(rea, *user, *password, *url)

	defer func(Body io.ReadCloser) {
		closeErr := Body.Close()
		if closeErr != nil {
			log.Fatalf("Error closing body: %v", closeErr)
		}
	}(res.Body)

	var graphQLResponse GraphQLResponse
	err = json.NewDecoder(res.Body).Decode(&graphQLResponse)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	fmt.Printf("\n%s", graphQLResponse)
}

// buildMutationReader receives the relevant parameters and builds an io.Reader of the mutation and its arguments
func buildMutationReader(jobName, namespace, image, cmd string, limit int) (io.Reader, error) {
	query := `mutation DeployJob($input: DeployJobInput!) {
  DeployJob(input: $input)
}`

	request := GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"input": DeployJobInput{
				Name:         jobName,
				Namespace:    namespace,
				Image:        image,
				Command:      cmd,
				BackOffLimit: limit,
			}},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonData), nil
}

// performGraphRequest receives the request body and login credentials to produce an HTTP request with a 10
func performGraphRequest(reader io.Reader, u, p, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		log.Fatalf("Error building job creation request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Operant Controller")
	req.SetBasicAuth(u, p)

	// set a 10s timeout, we don't want to hang if the service is unavailable
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error requesting job creation: %v", err)
	}

	return res, nil
}
