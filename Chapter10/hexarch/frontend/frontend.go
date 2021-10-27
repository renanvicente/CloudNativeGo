package frontend

import "github.com/renanvicente/CloudNativeGo/Chapter10/hexarch/core"

type FrontEnd interface {
	Start(kv *core.KeyValueStore) error
}
