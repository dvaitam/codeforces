package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type APIResponse struct {
	Status string `json:"status"`
	Result struct {
		Problems []Problem `json:"problems"`
	} `json:"result"`
}

type Problem struct {
	ContestId int    `json:"contestId"`
	Index     string `json:"index"`
	Name      string `json:"name"`
}

func main() {
	// 1. Get all problems from the Codeforces API
	apiProblems, err := getAllProblemsFromAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting problems from API: %v\n", err)
		os.Exit(1)
	}

	// 2. Get all existing problems from the file system
	localProblems := make(map[int]bool)
	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			base := filepath.Base(path)
			if contestId, err := strconv.Atoi(base); err == nil {
				localProblems[contestId] = true
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error walking directory: %v\n", err)
		os.Exit(1)
	}

	// 3. Compare the two lists
	var missingProblems []Problem
	for _, apiProblem := range apiProblems {
		if !localProblems[apiProblem.ContestId] {
			missingProblems = append(missingProblems, apiProblem)
		}
	}

	// 4. Write missing problems to a CSV file
	file, err := os.Create("truly_missing_problems.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create csv file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ContestID", "ProblemIndex", "ProblemName", "URL"}
	if err := writer.Write(header); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write header to csv: %v\n", err)
		os.Exit(1)
	}

	// Write data
	for _, p := range missingProblems {
		url := fmt.Sprintf("https://codeforces.com/problemset/problem/%d/%s", p.ContestId, p.Index)
		row := []string{strconv.Itoa(p.ContestId), p.Index, p.Name, url}
		if err := writer.Write(row); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write row to csv: %v\n", err)
		}
	}

	fmt.Println("Truly missing problems have been written to truly_missing_problems.csv")
}

func getAllProblemsFromAPI() ([]Problem, error) {
	resp, err := http.Get("https://codeforces.com/api/problemset.problems")
	if err != nil {
		return nil, fmt.Errorf("making API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling JSON: %w", err)
	}

	if apiResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned status: %s", apiResponse.Status)
	}

	return apiResponse.Result.Problems, nil
}