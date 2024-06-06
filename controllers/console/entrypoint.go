package console

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/kasaikou/is-your-human-resource-sufficient/models"
	"github.com/kasaikou/is-your-human-resource-sufficient/usecases/turn"
)

type MapService interface {
	turn.MapService
	InitialPositionID() (positionID string)
	NearestTarget(positionID string) (length int, targetPositionID string)
	ProgressPercent(positionID string) int
}

type gameHandler struct {
	numUpdateHRSteps int
	players          []models.Player
	mapService       MapService
	turnUseCases     *turn.TurnUseCases
	actionLogs       []models.Action
}

func NewGameContext(ctx context.Context, mapService MapService) {

	handler := &gameHandler{
		numUpdateHRSteps: 3,
		mapService:       mapService,
		turnUseCases:     turn.NewTurnService(mapService),
	}

	// プレイヤーの登録
	var numPlayers int
	fmt.Print("プレイヤーの人数を選択してください >")
	fmt.Scanf("%d", &numPlayers)
	handler.players = make([]models.Player, 0, numPlayers)

	for i := 0; i < numPlayers; i++ {
		var displayName string

		fmt.Print("\033[2J\033[1;1H")
		fmt.Println(handler.internalMakeTable())
		fmt.Println("")
		fmt.Printf("%d 人目のプレイヤ名を入力 >", i+1)
		fmt.Scanf("%s", &displayName)

		handler.players = append(handler.players, models.NewPlayer(displayName, handler.mapService.InitialPositionID(), 2))
	}

	for i := 0; i < len(handler.players); i++ {
		fmt.Print("\033[2J\033[1;1H")
		fmt.Println(handler.internalMakeTable())
		fmt.Println("")
		fmt.Printf("%s の一人当たりの最大移動距離を指定してください >", handler.players[i].DisplayName)
		fmt.Scanf("%d", &handler.players[i].NumPower)
	}

	goaled := 0
	for i := 0; goaled < len(handler.players); i++ {
		if i%handler.numUpdateHRSteps == 0 {
			for j := range handler.players {
				if handler.players[j].IsGoaled {
					continue
				}

				fmt.Print("\033[2J\033[1;1H")
				fmt.Println(handler.internalMakeTable())
				fmt.Println("")

				ranges := handler.turnUseCases.HumanResourceRange(turn.HumanResourceRangeRequest{Actor: handler.players[j]})
				fmt.Printf("[Turn %d] このターンでははじめにリソース調整を行います\n", i+1)
				for {
					fmt.Printf("%s の人数を指定してください [%d-%d] >", handler.players[j].DisplayName, ranges.Min, ranges.Max)
					var result int
					fmt.Scanf("%d", &result)
					if result >= ranges.Min && result <= ranges.Max {
						handler.players[j].HumanResource = result
						break
					}
				}

			}
		}

		for j := range handler.players {

			if handler.players[j].IsGoaled {
				continue
			}

			fmt.Print("\033[2J\033[1;1H")

			res := handler.turnUseCases.MakeActionContext(ctx, turn.MakeActionRequest{
				Actor: handler.players[j],
			})
			handler.players[j] = res.Actor

			fmt.Println(handler.internalMakeTable())
			fmt.Println("")

			switch res.Action.Result.Type {
			case models.ResultForward, models.ResultReached:
				concatStr := []string{}
				for _, move := range res.Action.MoveNum {
					concatStr = append(concatStr, strconv.Itoa(move))
				}
				fmt.Printf("[Turn %d] %s はこのターンで前進しました", i+1, handler.players[j].DisplayName)
				fmt.Printf("%s の %d 人はそれぞれ %s の仕事をして %s にきました", handler.players[j].DisplayName, handler.players[j].HumanResource, strings.Join(concatStr, ", "), handler.players[j].PositionID)

			case models.ResultBounded:
				concatStr := []string{}
				for _, move := range res.Action.MoveNum {
					concatStr = append(concatStr, strconv.Itoa(move))
				}
				fmt.Printf("[Turn %d] %s はこのターンで前進しましたが少し戻ってしまいました", i+1, handler.players[j].DisplayName)
				fmt.Printf("%s の %d 人はそれぞれ %s の仕事をして %s にきました", handler.players[j].DisplayName, handler.players[j].HumanResource, strings.Join(concatStr, ", "), handler.players[j].PositionID)

			case models.ResultGoaled:
				fmt.Printf("[Turn %d] %s はこのターンでゴールしました", i+1, handler.players[j].DisplayName)
				goaled++
			}

			fmt.Scanf("%s")
		}
	}
}
