package agent

import (
	"context"
	"log"
	"sync"
)

type Service interface {
	Name() string
	Run(context.Context) error
}

func Run(ctx context.Context, services []Service) error {
	var wg sync.WaitGroup
	var err error

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, service := range services {
		wg.Add(1)

		go func(s Service) {
			defer wg.Done()
			if err = s.Run(ctx); err != nil {
				log.Printf("Failing to start service %s: %v", s.Name(), err)
				cancel()
			}
		}(service)
	}

	wg.Wait()
	return err
}
