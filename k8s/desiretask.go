package k8s

import (
	"context"

	"github.com/julz/cube/opi"
	"k8s.io/api/core/v1"
	av1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	batch "k8s.io/api/batch/v1"
	"k8s.io/client-go/kubernetes"
)

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

func (d *TaskDesirer) DeleteJob(job string) error {
	return d.Client.BatchV1().Jobs("default").Delete(job, nil)
}

func toJob(task opi.Task) *batch.Job {
	job := &batch.Job{
		Spec: batch.JobSpec{
			ActiveDeadlineSeconds: int64ptr(600),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name:  "st8ge",
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
							v1.EnvVar{
								Name:  "CF_USERNAME",
								Value: task.Env["CF_USERNAME"],
							},
							v1.EnvVar{
								Name:  "CF_PASSWORD",
								Value: task.Env["CF_PASSWORD"],
							},
							v1.EnvVar{
								Name:  "API_ADDRESS",
								Value: task.Env["API_ADDRESS"],
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
