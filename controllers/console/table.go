package console

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mattn/go-runewidth"
)

func (gh *gameHandler) internalMakeTable() string {

	maxWidthName := 0
	nearestLengthExprs := make([]string, 0, len(gh.players))
	maxWidthNearestLengthExprs := 0
	for i, player := range gh.players {
		if player.DisplayName == "" {
			maxWidthName = max(runewidth.StringWidth("-"), maxWidthName)
		} else {
			maxWidthName = max(runewidth.StringWidth(player.DisplayName), maxWidthName)
		}

		length, nearest := gh.mapService.NearestTarget(player.PositionID)
		nearestLengthExprs = append(nearestLengthExprs, fmt.Sprintf("'%s' まで %d", nearest, length))
		maxWidthNearestLengthExprs = max(maxWidthNearestLengthExprs, runewidth.StringWidth(nearestLengthExprs[i]))
	}

	columnNames := []string{"プレイヤー名", "最大移動距離", "進捗率 [%]", "次の目標地点までの移動距離", "現在の人数", "ここまでの総コスト [人*時間]", "直前の状態"}
	columnWidthes := make([]int, 0, len(columnNames))
	for _, columnName := range columnNames {
		columnWidthes = append(columnWidthes, runewidth.StringWidth(columnName))
	}

	columnWidthes[0] = max(maxWidthName, columnWidthes[0])
	columnWidthes[2] = max(maxWidthNearestLengthExprs, columnWidthes[2])

	builder := strings.Builder{}
	builder.WriteString(runewidth.FillLeft(columnNames[0], columnWidthes[0]))

	for i := 1; i < len(columnNames); i++ {
		builder.WriteString(" | ")
		builder.WriteString(runewidth.FillLeft(columnNames[i], columnWidthes[i]))
	}

	length := runewidth.StringWidth(builder.String())
	builder.WriteByte('\n')
	builder.WriteString(strings.Repeat("=", length))
	builder.WriteByte('\n')

	for i, player := range gh.players {
		builder.WriteString(runewidth.FillLeft(player.DisplayName, columnWidthes[0]))
		builder.WriteString(" | ")
		builder.WriteString(runewidth.FillRight(strconv.Itoa(player.NumPower), columnWidthes[1]))
		builder.WriteString(" | ")
		builder.WriteString(runewidth.FillRight(strconv.Itoa(gh.mapService.ProgressPercent(player.PositionID)), columnWidthes[2]))
		builder.WriteString(" | ")
		builder.WriteString(runewidth.FillLeft(nearestLengthExprs[i], columnWidthes[3]))
		builder.WriteString(" | ")
		builder.WriteString(runewidth.FillRight(strconv.Itoa(player.HumanResource), columnWidthes[4]))
		builder.WriteString(" | ")
		builder.WriteString(runewidth.FillRight(strconv.Itoa(player.TotalCost), columnWidthes[5]))
		builder.WriteString(" | ")
		builder.WriteString(runewidth.FillLeft(func() string {
			if player.IsGoaled {
				return "GOALED"
			} else {
				return "INPROGRESS"
			}
		}(), columnWidthes[6]))
		builder.WriteByte('\n')
	}

	return builder.String()
}
