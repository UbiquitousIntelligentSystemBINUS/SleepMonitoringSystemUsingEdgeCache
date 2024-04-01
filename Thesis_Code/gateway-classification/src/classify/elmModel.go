package classify

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

type RealMatrix struct {
	Rows int
	Cols int
	Data [][]float64
}

type Weight struct {
	W  RealMatrix
	BW RealMatrix
}

type ELMModel struct {
	InputWeight     RealMatrix
	BiasInputWeight RealMatrix
	OutputWeight    RealMatrix
	Separator       string
}

func NewELMModel(inputWeightFilePath, outputWeightFilePath string) (*ELMModel, error) {
	var elmModel ELMModel

	inputWeightFileBytes, err := ioutil.ReadFile(inputWeightFilePath)
	if err != nil {
		return nil, err
	}

	outputWeightFileBytes, err := ioutil.ReadFile(outputWeightFilePath)
	if err != nil {
		return nil, err
	}

	inputWeightFileLines := strings.Split(string(inputWeightFileBytes), "\n")
	outputWeightFileLines := strings.Split(string(outputWeightFileBytes), "\n")

	weight, err := convertListCsvTo2dArr(inputWeightFileLines, true)
	if err != nil {
		return nil, err
	}
	elmModel.InputWeight = weight.W
	elmModel.BiasInputWeight = weight.BW

	weight, err = convertListCsvTo2dArr(outputWeightFileLines, false)
	if err != nil {
		return nil, err
	}
	elmModel.OutputWeight = weight.W

	return &elmModel, nil
}

func convertListCsvTo2dArr(input []string, useBias bool) (Weight, error) {
	var weight Weight

	listSize := len(input)
	if listSize == 0 {
		return weight, errors.New("empty input list")
	}

	csvElementSize := len(strings.Split(input[0], ","))
	for _, line := range input {
		if len(strings.Split(line, ",")) != csvElementSize {
			return weight, errors.New("invalid CSV length")
		}
	}

	weightData := make([][]float64, listSize)
	biasWeightData := make([][]float64, listSize)

	for i, line := range input {
		splittedLine := strings.Split(line, ",")
		temp := make([]float64, len(splittedLine))
		for j, val := range splittedLine {
			temp[j], _ = strconv.ParseFloat(strings.TrimSpace(val), 64)
		}
		if useBias {
			weightData[i] = temp[:len(temp)-1]
			biasWeightData[i] = []float64{temp[len(temp)-1]}
		} else {
			weightData[i] = temp
		}
	}

	weight.W = RealMatrix{Rows: listSize, Cols: len(weightData[0]), Data: weightData}
	weight.BW = RealMatrix{Rows: listSize, Cols: 1, Data: biasWeightData}

	return weight, nil
}
