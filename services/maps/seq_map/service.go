package seq_map

import (
	"slices"
	"strconv"

	"github.com/kasaikou/is-your-human-resource-sufficient/models"
)

type SequentialMapService struct {
	acceptable   int
	targetPoints []int
}

func NewSequentialMapService(acceptable int, lengthes ...int) *SequentialMapService {

	if len(lengthes) == 0 {
		panic("points argument is required")
	} else if lengthes[0] < 1 {
		panic("point must be larger than 0")
	} else if acceptable < 0 {
		panic("acceptable must not be smaller than 0")
	}

	smallerThanAcceptable := slices.ContainsFunc(lengthes, func(i int) bool {
		return i <= acceptable
	})
	if smallerThanAcceptable {
		panic("length must be larger than acceptable")
	}

	targetPoints := make([]int, 0, len(lengthes))
	targetPoints = append(targetPoints, lengthes[0])

	for _, length := range lengthes[1:] {
		targetPoints = append(targetPoints, targetPoints[len(targetPoints)-1]+length)
	}

	return &SequentialMapService{
		acceptable:   acceptable,
		targetPoints: targetPoints,
	}
}

var Default = NewSequentialMapService(2, 50, 40, 30, 20, 10, 20, 30, 30, 30)

func (service *SequentialMapService) internalNearestTarget(pos int) (length int, targetPoint int) {
	for _, targetPoint := range service.targetPoints {
		if pos+service.acceptable < targetPoint {
			return targetPoint - pos, targetPoint
		}
	}

	targetPoint = service.targetPoints[len(service.targetPoints)-1]
	if pos < targetPoint {
		return targetPoint - pos, targetPoint
	}

	return 0, targetPoint
}

func (service *SequentialMapService) InitialPositionID() (positionID string) {
	return strconv.Itoa(0)
}

func (service *SequentialMapService) NearestTarget(positionID string) (length int, targetPositionID string) {
	position, err := strconv.Atoi(positionID)
	if err != nil {
		panic(err)
	}

	length, targetPoint := service.internalNearestTarget(position)
	return length, strconv.Itoa(targetPoint)
}

func (service *SequentialMapService) internalProgressPercent(pos int) int {
	goal := service.targetPoints[len(service.targetPoints)-1]
	if pos >= goal {
		return 100
	} else {
		return int(100 * float32(pos) / float32(goal))
	}
}

func (service *SequentialMapService) ProgressPercent(positionID string) int {
	position, err := strconv.Atoi(positionID)
	if err != nil {
		panic(err)
	}

	return service.internalProgressPercent(position)
}

func (service *SequentialMapService) ResultByAction(action models.Action) (resultType models.ActionResultType, nextPositionID string) {

	position, err := strconv.Atoi(action.FromPositionID)
	if err != nil {
		panic(err)
	}

	moved := 0
	for _, n := range action.MoveNum {
		moved += n
	}

	length, targetPos := service.internalNearestTarget(position)
	if moved < length {
		return models.ResultForward, strconv.Itoa(position + moved)
	} else if moved == length {
		if service.targetPoints[len(service.targetPoints)-1] == targetPos {
			return models.ResultGoaled, strconv.Itoa(targetPos)
		} else {
			return models.ResultReached, strconv.Itoa(targetPos)
		}
	} else {
		pos := max(0, targetPos-(moved-length))
		return models.ResultBounded, strconv.Itoa(pos)
	}
}
