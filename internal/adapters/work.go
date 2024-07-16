package adapters

import (
	"github.com/kmyokoyama/go-template/internal/models"
	"github.com/kmyokoyama/go-template/internal/wire"
)

func ToWorkResponse(work models.Work) wire.WorkResponse {
	return wire.WorkResponse{Id: work.Id, Description: work.Description, Status: work.Status.String()}
}
