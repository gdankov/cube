package k8s

import (
	"context"

	"github.com/julz/cube/opi"
	"k8s.io/api/apps/v1beta1"
	batch "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	av1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Desirer struct {
	Client *kubernetes.Clientset
}

func (d *Desirer) Desire(ctx context.Context, lrps []opi.LRP) error {
	deployments, err := d.Client.AppsV1beta1().Deployments("default").List(av1.ListOptions{})
	if err != nil {
		return err
	}

	dByName := make(map[string]struct{})
	for _, d := range deployments.Items {
		dByName[d.Name] = struct{}{}
	}

	for _, lrp := range lrps {
		if _, ok := dByName[lrp.Name]; ok {
			continue
		}

		if _, err := d.Client.AppsV1beta1().Deployments("default").Create(toDeployment(lrp)); err != nil {
			// fixme: this should be a multi-error and deferred
			return err
		}
	}

	return nil
}

func toDeployment(lrp opi.LRP) *v1beta1.Deployment {
	deployment := &v1beta1.Deployment{
		Spec: v1beta1.DeploymentSpec{
			Replicas: int32ptr(lrp.TargetInstances),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name:  "web",
						Image: lrp.Image,
					}},
				},
			},
		},
	}

	deployment.Name = lrp.Name
	deployment.Spec.Template.Labels = map[string]string{
		"name": lrp.Name,
	}

	deployment.Labels = map[string]string{
		"cube": "cube",
		"name": lrp.Name,
	}

	return deployment
}

type TaskDesirer struct {
	Client *kubernetes.Clientset
}

func (d TaskDesirer) Desire(ctx context.Context, tasks []opi.Task) error {
	jobs, err := d.Client.BatchV1().Jobs("default").List(av1.ListOptions{})
	if err != nil {
		return err
	}

	dByName := make(map[string]struct{})
	for _, d := range jobs.Items {
		dByName[d.Name] = struct{}{}
	}

	for _, task := range tasks {
		//if _, ok := dByName[task.Name]; ok {
		//continue
		//}

		if _, err := d.Client.BatchV1().Jobs("default").Create(toJob(task)); err != nil {
			// fixme: this should be a multi-error and deferred
			return err
		}
	}

	return nil
}

func toJob(task opi.Task) *batch.Job {
	job := &batch.Job{
		Spec: batch.JobSpec{
			ActiveDeadlineSeconds: int64ptr(10),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name:  "bumblebee",
						Image: task.Image,
						Env: []v1.EnvVar{
							v1.EnvVar{
								Name:  "DOWNLOAD_URL",
								Value: task.Env["DOWNLOAD_URL"],
							},
							v1.EnvVar{
								Name:  "UPLOAD_URL",
								Value: task.Env["UPLOAD_URL"],
							},
							v1.EnvVar{
								Name:  "APP_ID",
								Value: task.Env["APP_ID"],
							},
							v1.EnvVar{
								Name:  "STAGING_GUID",
								Value: task.Env["STAGING_GUID"],
							},
						},
					}},
					RestartPolicy: "Never",
				},
			},
		},
	}

	job.Name = task.Env["STAGING_GUID"]

	job.Spec.Template.Labels = map[string]string{
		"name": task.Env["APP_ID"],
	}

	job.Labels = map[string]string{
		"cube": "cube",
		"name": task.Env["APP_ID"],
	}
	return job
}

func int32ptr(i int) *int32 {
	u := int32(i)
	return &u
}

func int64ptr(i int) *int64 {
	u := int64(i)
	return &u
}
