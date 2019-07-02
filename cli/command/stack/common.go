package stack

import (
	"golang.org/x/net/context"

	"github.com/alcideio/moby/api/types"
	"github.com/alcideio/moby/api/types/filters"
	"github.com/alcideio/moby/api/types/swarm"
	"github.com/alcideio/moby/cli/compose/convert"
	"github.com/alcideio/moby/client"
	"github.com/alcideio/moby/opts"
)

func getStackFilter(namespace string) filters.Args {
	filter := filters.NewArgs()
	filter.Add("label", convert.LabelNamespace+"="+namespace)
	return filter
}

func getStackFilterFromOpt(namespace string, opt opts.FilterOpt) filters.Args {
	filter := opt.Value()
	filter.Add("label", convert.LabelNamespace+"="+namespace)
	return filter
}

func getAllStacksFilter() filters.Args {
	filter := filters.NewArgs()
	filter.Add("label", convert.LabelNamespace)
	return filter
}

func getServices(
	ctx context.Context,
	apiclient client.APIClient,
	namespace string,
) ([]swarm.Service, error) {
	return apiclient.ServiceList(
		ctx,
		types.ServiceListOptions{Filters: getStackFilter(namespace)})
}

func getStackNetworks(
	ctx context.Context,
	apiclient client.APIClient,
	namespace string,
) ([]types.NetworkResource, error) {
	return apiclient.NetworkList(
		ctx,
		types.NetworkListOptions{Filters: getStackFilter(namespace)})
}

func getStackSecrets(
	ctx context.Context,
	apiclient client.APIClient,
	namespace string,
) ([]swarm.Secret, error) {
	return apiclient.SecretList(
		ctx,
		types.SecretListOptions{Filters: getStackFilter(namespace)})
}
