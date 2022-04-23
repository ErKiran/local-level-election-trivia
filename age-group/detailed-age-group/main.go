package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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

type CasteData struct {
	Age   string `json:"Age"`
	Count int    `json:"Count"`
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

func convertJSONToCSV(electionData []CasteData, destination, post string) error {
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"S.No", "Post", "Age", "Count"}
	if err := writer.Write(header); err != nil {
		return err
	}

	var count = 0
	for _, e := range electionData {
		count++
		var csvRow []string
		csvRow = append(csvRow, strconv.Itoa(count), post, e.Age, strconv.Itoa(e.Count))
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	electionData, err := ReadAndParseData()
	var casteData []CasteData
	if err != nil {
		fmt.Println(err)
	}

	var mainMap = make(map[int]int)
	post := "उपप्रमुख"

	for _, data := range electionData {
		if data.Postname == post && data.Estatus == "E" {
			mainMap[data.Age] += 1
		}
	}
	var dataRangeMap = make(map[string]int)
	for k, v := range mainMap {
		if k >= 21 && k <= 29 {
			dataRangeMap["21-29"] += v
		} else if k >= 30 && k <= 39 {
			dataRangeMap["30-39"] += v
		} else if k >= 40 && k <= 49 {
			dataRangeMap["40-49"] += v
		} else if k >= 50 && k <= 59 {
			dataRangeMap["50-59"] += v
		} else if k >= 60 && k <= 69 {
			dataRangeMap["60-69"] += v
		} else if k >= 70 && k <= 79 {
			dataRangeMap["70-79"] += v
		} else if k >= 80 && k <= 89 {
			dataRangeMap["80-89"] += v
		}
	}
	fmt.Println("dataRangeMap", dataRangeMap)
	// fmt.Println(mainMap)

	for key, value := range dataRangeMap {
		casteData = append(casteData, CasteData{
			Age:   key,
			Count: value,
		})
	}

	csvFileName := fmt.Sprintf("%s.csv", post)
	os.Create(csvFileName)
	if err := convertJSONToCSV(casteData, csvFileName, post); err != nil {
		log.Fatal(err)
	}
}
