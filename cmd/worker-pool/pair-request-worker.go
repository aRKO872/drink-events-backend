package worker

import (
	"context"
	"time"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/dao"
	"github.com/drink-events-backend/pkg/business_logic/rao"
	"github.com/gocraft/work"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func (c *Context) AddPairRequest(job *work.Job) error {
	from := job.ArgString("from")
	to := job.ArgString("to")
	if err := job.ArgError(); err != nil {
		return err
	}

	id := uuid.NewString()
	createdAt := time.Now().Format(literals.DATE_FORMAT)

	pairReq := models.PairRequest{
		Id: id,
		SenderID: from,
		ReceiverID: to,
		IsActive: true,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	dao, getDBErr := dao.GetDBAccessOperator()
	if getDBErr != nil {
		return getDBErr
	}

	redisUserOp, redisInitErr := rao.GetRedisAccessOperator()
	if redisInitErr != nil {
		return redisInitErr
	}

	wg := new(errgroup.Group)

	wg.Go(func() error {
		return redisUserOp.AddPairRequest(context.Background(), &pairReq)
	})

	wg.Go(func() error {
		return dao.AddPairRequest(&pairReq)
	})

	if err := wg.Wait(); err != nil {
		return err
	}

	return nil
}