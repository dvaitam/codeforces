package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
)

// DataSet represents the structure of the training data.
type DataSet struct {
	Features [][]float64 `json:"Features"`
	Labels   []int       `json:"Labels"`
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func main() {
	// try to locate the dataset file
	path := os.Getenv("DATASET_PATH")
	if path == "" {
		path = "dataset.json"
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		// dataset not found; nothing to output
		return
	}
	var ds DataSet
	if err := json.Unmarshal(data, &ds); err != nil {
		return
	}

	if len(ds.Features) == 0 || len(ds.Features) != len(ds.Labels) {
		return
	}

	w := [2]float64{0, 0}
	b := 0.0
	lr := 0.1
	for iter := 0; iter < 1000; iter++ {
		gradW := [2]float64{0, 0}
		gradB := 0.0
		for i, x := range ds.Features {
			y := ds.Labels[i]
			z := w[0]*x[0] + w[1]*x[1] + b
			p := sigmoid(z)
			diff := p - float64(y)
			gradW[0] += diff * x[0]
			gradW[1] += diff * x[1]
			gradB += diff
		}
		n := float64(len(ds.Features))
		w[0] -= lr * gradW[0] / n
		w[1] -= lr * gradW[1] / n
		b -= lr * gradB / n
	}

	out := struct {
		Weights [2]float64 `json:"weights"`
		Bias    float64    `json:"bias"`
	}{w, b}
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(out)
}
