package backend

import (
	"fmt"

	"github.com/julianstephens/warden/internal/backend/common"
	"github.com/julianstephens/warden/internal/backend/local"
)

func NewBackend(t common.BackendType, params common.Params) (common.Backend, error) {
	switch t {
	case common.LocalStorage:
		return local.NewLocalStorage(params.(common.LocalStorageParams))
	default:
		return nil, fmt.Errorf("invalid backend type: %s", t.String())
	}
}
