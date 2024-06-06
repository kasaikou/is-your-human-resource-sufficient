package turn

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/kasaikou/is-your-human-resource-sufficient/models"
)

type MapService interface {
	ResultByAction(action models.Action) (resultType models.ActionResultType, nextPositionID string)
}

type TurnUseCases struct {
	mapService MapService
	rand       *rand.Rand
}

func NewTurnService(mapService MapService) *TurnUseCases {
	return &TurnUseCases{
		rand:       rand.New(rand.NewSource(time.Now().Unix())),
		mapService: mapService,
	}
}

type MakeActionRequest struct {
	Actor models.Player
}

type MakeActionResponse struct {
	Actor  models.Player
	Action models.Action
}

func (useCase TurnUseCases) MakeActionContext(ctx context.Context, req MakeActionRequest) MakeActionResponse {

	if req.Actor.IsGoaled {
		panic(fmt.Sprintf("actor '%s' (%s) has been already goaled", req.Actor.DisplayName, req.Actor.ID))
	}

	actor := req.Actor
	action := models.NewAction(&actor)
	action.MoveNum = make([]int, 0, actor.HumanResource)
	for i := 0; i < actor.HumanResource; i++ {
		actor.TotalCost++
		action.MoveNum = append(action.MoveNum, useCase.rand.Intn(actor.NumPower)+1)
	}

	resultType, nextPositionID := useCase.mapService.ResultByAction(action)
	actor.PositionID = nextPositionID

	action.Result = &models.ActionResult{
		Action:     &action,
		Type:       resultType,
		PositionID: nextPositionID,
	}

	if resultType == models.ResultGoaled {
		actor.IsGoaled = true
	}

	return MakeActionResponse{
		Actor:  actor,
		Action: action,
	}
}

type HumanResourceRangeRequest struct {
	Actor models.Player
}

type HumanResourceRangeResponse struct {
	Min, Max int
}

func (useCase TurnUseCases) HumanResourceRange(req HumanResourceRangeRequest) HumanResourceRangeResponse {

	actor := req.Actor
	return HumanResourceRangeResponse{
		Min: (actor.HumanResource + 1) / 2,
		Max: actor.HumanResource * 3,
	}
}
