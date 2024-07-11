package adapters

import (
	"github.com/kmyokoyama/go-template/internal/models"
	"github.com/kmyokoyama/go-template/internal/wire"
)

func ToVersionResponse(m models.Version) wire.VersionResponse {
	return wire.VersionResponse{Version: m.Version}
}