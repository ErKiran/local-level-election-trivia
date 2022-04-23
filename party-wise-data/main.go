package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type ElectionData struct {
	Partyid            int    `json:"PartyID"`
	Stateid            int    `json:"StateID"`
	Candidatename      string `json:"CandidateName"`
	Gender             string `json:"Gender"`
	Age                int    `json:"Age"`
	Politicalpartyname string `json:"PoliticalPartyName"`
	Districtname       string `json:"DistrictName"`
	Localbodyname      string `json:"LocalBodyName"`
	Wardno             string `json:"WardNo"`
	Postname           string `json:"PostName"`
	Serialno           int    `json:"SerialNo"`
	Totalvotesrecieved int    `json:"TotalVotesRecieved"`
	Estatus            string `json:"EStatus"`
	Rank               int    `json:"Rank"`
}

type PartyWiseData struct {
	Post           string  `json:"Post"`
	Party          string  `json:"Party"`
	TotalVote      int     `json:"TotalVote"`
	WinCount       int     `json:"WinCount"`
	TotalCandidate int     `json:"TotalCandidate"`
	WinPercentage  float64 `json:"WinPercentage"`
}

func ReadAndParseData() ([]ElectionData, error) {
	jsonFile, err := os.Open("./local-level-election/raw/alldata.json")

	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var electionData []ElectionData

	err = json.Unmarshal(byteValue, &electionData)

	if err != nil {
		return nil, err
	}

	return electionData, nil
}

func convertJSONToCSV(electionData []PartyWiseData, destination, post string) error {
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"S.No", "Post", "Party", "Total Vote", "Win Count", "Total Candidate", "Win Percentage"}
	if err := writer.Write(header); err != nil {
		return err
	}

	var count = 0
	for _, e := range electionData {
		count++
		var csvRow []string
		csvRow = append(csvRow, strconv.Itoa(count), post, e.Party, strconv.Itoa(e.TotalVote), strconv.Itoa(e.WinCount), strconv.Itoa(e.TotalCandidate), strconv.FormatFloat(e.WinPercentage, 'E', -1, 32))
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	electionData, err := ReadAndParseData()
	var casteData []PartyWiseData
	if err != nil {
		fmt.Println(err)
	}

	var mainMap = make(map[string]int)
	var winMap = make(map[string]int)
	var candidateMap = make(map[string]int)
	post := "अध्यक्ष"

	for _, data := range electionData {
		if data.Postname == post {
			mainMap[fmt.Sprintf("%s__%s", data.Politicalpartyname, data.Postname)] += data.Totalvotesrecieved
			candidateMap[fmt.Sprintf("%s__%s", data.Politicalpartyname, data.Postname)] += 1
		}

		if data.Estatus == "E" && data.Postname == post {
			winMap[fmt.Sprintf("%s__%s", data.Politicalpartyname, data.Postname)] += 1
		}
	}

	for key, value := range mainMap {
		name := strings.Split(key, "__")
		casteData = append(casteData, PartyWiseData{
			Post:           name[1],
			Party:          name[0],
			TotalVote:      value,
			WinCount:       winMap[key],
			TotalCandidate: candidateMap[key],
			WinPercentage:  float64(winMap[key]) / float64(candidateMap[key]) * 100,
		})
	}

	csvFileName := fmt.Sprintf("%s.csv", post)
	os.Create(csvFileName)
	if err := convertJSONToCSV(casteData, csvFileName, post); err != nil {
		log.Fatal(err)
	}
}
