package roundtest

import (
	"math"

	"github.com/shopspring/decimal"
)

// RoundByPosition 倍精度浮動小数点数を指定桁で四捨五入する
func RoundByPosition(num float64, pos int32) float64 {
	result, _ := decimal.NewFromFloat(num).Round(pos).Float64()
	return result
}

// RoundByPositionNG1 四捨五入する桁を整数にして、整数で計算するパターン
func RoundByPositionNG1(num float64, pos int32) float64 {
	shift := math.Pow10(int(pos + 1))
	result := roundInt(num*shift) / shift
	return result
}

// 整数の1桁目を四捨五入する
func roundInt(num float64) float64 {
	t := num - math.Mod(num, 10)
	if math.Abs(num-t) >= 5 {
		return t + math.Copysign(10, num)
	}
	return t
}

//  RoundByPositionNG2 math.Round関数使うパターン
func RoundByPositionNG2(num float64, pos int32) float64 {
	shift := math.Pow10(int(pos))
	result := math.Round(num*shift) / shift
	return result
}
