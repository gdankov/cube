package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/julz/cube/opi"
	"k8s.io/api/apps/v1beta1"
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
	//command, args := splitCommandAndArgs(lrp.Command[0])
	deployment := &v1beta1.Deployment{
		Spec: v1beta1.DeploymentSpec{
			Replicas: int32ptr(lrp.TargetInstances),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name:  "web",
						Image: lrp.Image,
						Command: []string{
							"/bin/bash",
						}, //command,
						Args: []string{
							"-c",
							fmt.Sprintf("%s %s", launch, lrp.Command[0]),
						}, //args,
						Env:             mapToEnvVar(lrp.Env),
						ImagePullPolicy: "Always",
					}},
				},
			},
		},
	}

	fmt.Println("PROVIDED ENV VARS:", lrp.Env)

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

func splitCommandAndArgs(cmd string) ([]string, []string) {
	all := strings.Split(cmd, " ")
	return []string{"app/bin/" + all[0]}, all[1:len(all)]
}

func int32ptr(i int) *int32 {
	u := int32(i)
	return &u
}

func int64ptr(i int) *int64 {
	u := int64(i)
	return &u
}

const launch = `
cd /home/vcap/app

HOME=/home/vcap
PATH=/usr/local/bin:/usr/bin:/bin
USER=vcap
BUNDLE_GEMFILE=/home/vcap/app/Gemfile
CF_INSTANCE_INTERNAL_IP=0.0.0.0
CF_INSTANCE_IP=0.0.0.0
CF_STACK=cflinuxfs2
HOME=/home/vcap/app

echo "launching app.."

DEPS_DIR=../deps

if [ -n "$(ls ../profile.d/* 2> /dev/null)" ]; then
  for env_file in ../profile.d/*; do
    source $env_file
	  done
fi

if [ -n "$(ls .profile.d/* 2> /dev/null)" ]; then
  for env_file in .profile.d/*; do
    source $env_file
	  done
fi

if [ -f .profile ]; then
  source .profile
fi

echo "deps dir: $DEPS_DIR"

echo "executing command $@"
exec bash -c "bundle exec rackup config.ru -p 8080"
`
