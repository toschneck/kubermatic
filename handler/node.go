package handler

import (
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/kubermatic/api/provider"
	"golang.org/x/net/context"
)

func Nodes(
	ctx context.Context,
	kp provider.KubernetesProvider,
	cps map[int]provider.CloudProvider,
) http.Handler {
	return httptransport.NewServer(
		ctx,
		nodesEndpoint(kp, cps),
		decodeNodesReq,
		encodeJSON,
	)
}

func nodesEndpoint(
	kp provider.KubernetesProvider,
	cps map[int]provider.CloudProvider,
) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*nodesReq)

		c, err := kp.Cluster(req.dc, req.cluster)
		if err != nil {
			return nil, err
		}

		cpIdx, err := provider.ClusterCloudProvider(c)
		if err != nil {
			return nil, err
		}

		cp, found := cps[cpIdx]
		if !found {
			return nil, fmt.Errorf("unsupported cloud provider")
		}

		return cp.Nodes(c)
	}
}

type nodesReq struct {
	dcReq
	cluster string
}

func decodeNodesReq(r *http.Request) (interface{}, error) {
	var req nodesReq

	dr, err := decodeDcReq(r)
	if err != nil {
		return nil, err
	}
	req.dcReq = *dr.(*dcReq)

	req.cluster = mux.Vars(r)["cluster"]

	return req, nil
}
