package frontend

import "github.com/renanvicente/grpc_sample/hexarch/core"

type FrontEnd interface {
	Start(kv *core.KeyValueStore) error
}
