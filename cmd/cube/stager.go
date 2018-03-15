package main

import (
	"log"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/julz/cube/k8s"
	"github.com/julz/cube/stager"
	"github.com/urfave/cli"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func stagingCmd(c *cli.Context) {
	config, err := clientcmd.BuildConfigFromFlags("", c.String("kubeconfig"))

	exitWithError(err)
	clientset, err := kubernetes.NewForConfig(config)
	exitWithError(err)

	taskDesirer := k8s.TaskDesirer{Client: clientset}

	st8r := stager.Stager{
		taskDesirer,
	}

	logger := lager.NewLogger("st8r")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	backend := stager.NewBackend(logger)

	handler := stager.New(st8r, backend, logger)

	log.Fatal(http.ListenAndServe("0.0.0.0:8085", handler))
}
