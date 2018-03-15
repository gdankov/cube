package stager

import (
	"context"

	"github.com/julz/cube/opi"
)

type Stager struct {
	Desirer opi.TaskDesirer
}

func (s Stager) Run(task opi.Task) error {
	err := s.Desirer.Desire(context.Background(), []opi.Task{task})
	if err != nil {
		return err
	}
	return nil
}
