package roundtest_test

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"testing"

	roundtest "github.com/fujiwaram/round-test"
)

func TestRoundByPosition(t *testing.T) {
	const (
		intDigitNum = 3
		decDigitNum = 4
	)

	tests := []struct {
		name     string
		from, to int // 四捨五入する位置
	}{
		{name: "1-3", from: 1, to: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roundTest(intDigitNum, decDigitNum, tt.from, tt.to,
				roundtest.RoundByPosition,
				func(got, want, base float64) {
					t.Errorf("round.RoundByPosition() = %v, want %v, base %v", got, want, base)
				},
			)
		})
	}
}

func TestRoundByPositionNG1(t *testing.T) {
	const (
		intDigitNum = 3
		decDigitNum = 4
	)

	tests := []struct {
		name     string
		from, to int // 四捨五入する位置
	}{
		{name: "1-3", from: 1, to: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roundTest(intDigitNum, decDigitNum, tt.from, tt.to,
				roundtest.RoundByPositionNG1,
				func(got, want, base float64) {
					t.Errorf("round.RoundByPositionNG1() = %v, want %v, base %v", got, want, base)
				},
			)
		})
	}
}

func TestRoundByPositionNG2(t *testing.T) {
	const (
		intDigitNum = 3
		decDigitNum = 4
	)

	tests := []struct {
		name     string
		from, to int // 四捨五入する位置
	}{
		{name: "1-3", from: 1, to: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roundTest(intDigitNum, decDigitNum, tt.from, tt.to,
				roundtest.RoundByPositionNG2,
				func(got, want, base float64) {
					t.Errorf("round.RoundByPositionNG2() = %v, want %v, base %v", got, want, base)
				},
			)
		})
	}
}

func roundTest(intDigitNum, decDigitNum int, posFrom, posTo int,
	roundFunc func(num float64, pos int32) float64,
	failCallback func(got, want, base float64),
) {
	intMax := int(math.Pow10(intDigitNum))
	decMax := int(math.Pow10(decDigitNum))
	format := fmt.Sprintf("%%d.%%0%dd", decDigitNum)

	limit := make(chan struct{}, 6)
	var wg sync.WaitGroup
	// Integer part
	for i := 0; i < intMax; i++ {
		wg.Add(1)
		go func(i int) {
			limit <- struct{}{}
			defer func() {
				<-limit
				wg.Done()
			}()
			// Decimal part
			for d := 0; d < decMax; d++ {
				numStr := fmt.Sprintf(format, i, d)

				for p := posFrom; p < posTo; p++ {
					wi, wd := calcWantRoundedNum(i, d, decDigitNum-p, decMax)
					wantNumStr := fmt.Sprintf(format, wi, wd)

					num, _ := strconv.ParseFloat(numStr, 64)
					wantNum, _ := strconv.ParseFloat(wantNumStr, 64)

					if got := roundFunc(num, int32(p)); got != wantNum {
						failCallback(got, wantNum, num)
					}
				}
			}
		}(i)
	}
	wg.Wait()
}

func calcWantRoundedNum(intNum, decNum, pos, decMax int) (wi, wd int) {
	wd = calcWantDecimal(decNum, pos)
	wi = intNum
	if wd >= decMax {
		wi++
		wd = 0
	}
	return wi, wd
}

func calcWantDecimal(num, pos int) int {
	roundCoef := int(math.Pow10(pos))
	q := num / roundCoef * roundCoef
	coef := roundCoef / 10
	if num-q < 5*coef {
		return q
	}
	return q + 10*coef
}
