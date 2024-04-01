package classify

import (
	"math"
	"sort"
)

type HRVFeature struct {
	F01_AVNN               float64
	F02_SDNN               float64
	F03_RMSSD              float64
	F04_SDSD               float64
	F05_NNx                float64
	F06_PNNx               float64
	F07_HRV_TRIANGULAR_IDX float64
	F08_SD1                float64
	F09_SD2                float64
	F10_SD1_SD2_RATIO      float64
	F11_S                  float64
	F12_TP                 float64
	F13_pLF                float64
	F14_pHF                float64
	F15_LFHFratio          float64
	F16_VLF                float64
	F17_LF                 float64
	F18_HF                 float64
}

func NewHRVFeature(rrIntervalSet RRIntervalSet) *HRVFeature {
	var hrv HRVFeature

	rrIntervalValue := rrIntervalSet.RRIntervalValue
	rrIntervalsValueDiff := rrIntervalSet.RRIntervalsValueDiff

	hrv.F01_AVNN = f01_AVNN(rrIntervalValue)
	hrv.F02_SDNN = f02_SDNN(rrIntervalValue)
	hrv.F03_RMSSD = f03_RMSSD(rrIntervalsValueDiff)
	hrv.F04_SDSD = f04_SDSD(rrIntervalsValueDiff)
	hrv.F05_NNx = f05_NNx(rrIntervalsValueDiff, 50)
	hrv.F06_PNNx = f06_PNNx(rrIntervalValue, hrv.F05_NNx)
	hrv.F07_HRV_TRIANGULAR_IDX = f07_HRV_TRIANGULAR_IDX(rrIntervalValue)
	hrv.F08_SD1 = f08_SD1(hrv.F04_SDSD)
	hrv.F09_SD2 = f09_SD2(hrv.F02_SDNN, hrv.F04_SDSD)
	hrv.F10_SD1_SD2_RATIO = f10_SD1_SD2_RATIO(hrv.F08_SD1, hrv.F09_SD2)
	hrv.F11_S = f11_S(hrv.F08_SD1, hrv.F09_SD2)

	feature12To18 := f12_18(rrIntervalValue, 2)
	hrv.F12_TP = feature12To18.TP
	hrv.F13_pLF = feature12To18.pLF
	hrv.F14_pHF = feature12To18.pHF
	hrv.F15_LFHFratio = feature12To18.LFHFratio
	hrv.F16_VLF = feature12To18.VLF
	hrv.F17_LF = feature12To18.LF
	hrv.F18_HF = feature12To18.HF

	return &hrv
}

func f01_AVNN(rrIntervalValue []float64) float64 {
	return mean(rrIntervalValue)
}

func f02_SDNN(rrIntervalValue []float64) float64 {
	return sampleStandardDeviation(rrIntervalValue)
}

func f03_RMSSD(rrIntervalsValueDiff []float64) float64 {
	return math.Sqrt(mean(powList(rrIntervalsValueDiff, 2)))
}

func f04_SDSD(rrIntervalsValueDiff []float64) float64 {
	return sampleStandardDeviation(rrIntervalsValueDiff)
}

func f05_NNx(rrIntervalsValueDiff []float64, x float64) float64 {
	count := filterCount(mulList(rrIntervalsValueDiff, 1000), func(y float64) bool {
		return y > x
	})
	return float64(count)
}

func f06_PNNx(rrIntervalValue []float64, NNx float64) float64 {
	return (NNx / (float64(len(rrIntervalValue)) - 1)) * 100
}

func f07_HRV_TRIANGULAR_IDX(rrIntervalValue []float64) float64 {
	binSize := 7.812
	var tempRr []float64
	for _, val := range rrIntervalValue {
		tempRr = append(tempRr, val*1000)
	}

	sort.Float64s(tempRr)
	maxVal := tempRr[len(tempRr)-1]
	minVal := tempRr[0]

	binCount := math.Ceil((maxVal - minVal) / binSize)
	edges := make([]float64, int(binCount)+1)
	var Nds []float64
	edges[0] = minVal

	for i := 1; i <= int(binCount); i++ {
		edges[i] = edges[i-1] + binSize
		var d float64
		for _, x := range tempRr {
			if x >= edges[i-1] && x < edges[i] {
				d++
			}
		}
		if d != 0 {
			Nds = append(Nds, d)
		}
	}

	return max(Nds) / sum(Nds)
}

func f08_SD1(sdsd float64) float64 {
	return math.Sqrt(math.Pow(sdsd, 2) / 2)
}

func f09_SD2(sdnn, sdsd float64) float64 {
	return math.Sqrt(2*math.Pow(sdnn, 2) - math.Pow(sdsd, 2)/2)
}

func f10_SD1_SD2_RATIO(sd1, sd2 float64) float64 {
	return sd1 / sd2
}

func f11_S(sd1, sd2 float64) float64 {
	return math.Pi * sd1 * sd2
}

func f12_18(rrIntervalValue []float64, Fs float64) Feature12To18 {
	var feature Feature12To18
	// Implement the logic for feature calculation
	return feature
}

func filterF(YY []float64, f []float64, predicate func(float64) bool) []float64 {
	var result []float64
	for i, val := range f {
		if predicate(val) {
			result = append(result, YY[i])
		}
	}
	return result
}

func nanzscore(input []float64) []float64 {
	m := nanmean(input)
	s := nanstd(input)

	var z []float64
	for _, val := range input {
		z = append(z, (val-m)/s)
	}
	return z
}

func nanmean(input []float64) float64 {
	return mean(input)
}

func nanstd(input []float64) float64 {
	return populationStandardDeviation(input)
}

type Feature12To18 struct {
	TP        float64
	pLF       float64
	pHF       float64
	LFHFratio float64
	VLF       float64
	LF        float64
	HF        float64
}

type RRIntervalSet struct {
	RRIntervalValue      []float64
	RRIntervalsValueDiff []float64
}

func mean(arr []float64) float64 {
	sum := 0.0
	for _, val := range arr {
		sum += val
	}
	return sum / float64(len(arr))
}

func sampleStandardDeviation(arr []float64) float64 {
	mean := mean(arr)
	variance := 0.0
	for _, val := range arr {
		variance += math.Pow(val-mean, 2)
	}
	return math.Sqrt(variance / float64(len(arr)-1))
}

func powList(arr []float64, exp float64) []float64 {
	var result []float64
	for _, val := range arr {
		result = append(result, math.Pow(val, exp))
	}
	return result
}

func filterCount(arr []float64, predicate func(float64) bool) int {
	count := 0
	for _, val := range arr {
		if predicate(val) {
			count++
		}
	}
	return count
}

func sum(arr []float64) float64 {
	total := 0.0
	for _, val := range arr {
		total += val
	}
	return total
}

func max(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}
	max := arr[0]
	for _, val := range arr {
		if val > max {
			max = val
		}
	}
	return max
}

func populationStandardDeviation(input []float64) float64 {
	meanVal := mean(input)
	variance := 0.0
	for _, val := range input {
		variance += math.Pow(val-meanVal, 2)
	}
	return math.Sqrt(variance / float64(len(input)))
}

func mulList(input []float64, multiplier float64) []float64 {
	result := make([]float64, len(input))
	for i, val := range input {
		result[i] = val * multiplier
	}
	return result
}
