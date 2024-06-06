package models

type Action struct {
	Actor          *Player
	HumanResource  int
	MoveNum        []int
	FromPositionID string
	Result         *ActionResult
}

func NewAction(actor *Player) Action {
	return Action{
		Actor:          actor,
		HumanResource:  actor.HumanResource,
		FromPositionID: actor.PositionID,
	}
}

type ActionResultType string

const (
	ResultForward  ActionResultType = "FORWARD"
	ResultBounded  ActionResultType = "BOUNDED"
	ResultReached  ActionResultType = "REACHED"
	ResultGoaled   ActionResultType = "GOALED"
	ResultRollback ActionResultType = "ROLLBACK"
)

type ActionResult struct {
	Action     *Action
	Type       ActionResultType
	PositionID string
}
