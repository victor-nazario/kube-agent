package agent

import (
	"context"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"strings"
)

type Deployer interface {
	Deploy(ctx context.Context, jobName, image, cmd, nameSpace string, backOffLimit int32) (string, error)
}

type Agent struct {
	client client.Client
}

// NewAgent returns a configured agent type, along with its initialized client
func NewAgent() (Agent, error) {
	mg, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		return Agent{}, err
	}

	var agent Agent
	agent.client = mg.GetClient()

	return agent, nil
}

// Deploy receives job information and behaves like a K8s Controller deploying a job using the v1.batch api
func (a *Agent) Deploy(ctx context.Context, jobName, image, cmd, nameSpace string, backOffLimit int32) (string, error) {
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: nameSpace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    jobName,
							Image:   image,
							Command: strings.Split(cmd, " "),
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	err := a.client.Create(ctx, jobSpec)
	if err != nil {
		return "", err
	}

	return string(jobSpec.GetUID()), nil
}
