package loader

import (
	"context"
	"sync"

	"github.com/nicksrandall/dataloader"
	"github.com/tonyghita/graphql-go-example/swapi"
)

type VehicleGetter interface {
	Vehicle(ctx context.Context, url string) (swapi.Vehicle, error)
}

type VehicleLoader struct {
	get VehicleGetter
}

func NewVehicleLoader(client VehicleGetter) dataloader.BatchFunc {
	return VehicleLoader{get: client}.loadBatch
}

func (loader VehicleLoader) loadBatch(ctx context.Context, urls []string) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(ctx context.Context, i int, url string) {
			data, err := loader.get.Vehicle(ctx, url)
			results[i] = &dataloader.Result{Data: data, Error: err}
			wg.Done()
		}(ctx, i, url)
	}

	wg.Wait()

	return results
}
