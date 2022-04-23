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

type CasteData struct {
	Caste string `json:"Caste"`
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

	header := []string{"S.No", "Post", "Caste", "Count"}
	if err := writer.Write(header); err != nil {
		return err
	}

	var count = 0
	for _, e := range electionData {
		count++
		var csvRow []string
		csvRow = append(csvRow, strconv.Itoa(count), post, e.Caste, strconv.Itoa(e.Count))
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

	var mainMap = make(map[string]int)
	post := "वडा अध्यक्ष"

	for _, data := range electionData {
		if data.Postname == post && data.Estatus == "E" {
			name := strings.Split(data.Candidatename, " ")
			mainMap[name[len(name)-1]] += 1
		}
	}
	fmt.Println(mainMap)

	for key, value := range mainMap {
		casteData = append(casteData, CasteData{
			Caste: key,
			Count: value,
		})
	}

	csvFileName := fmt.Sprintf("%s.csv", post)
	os.Create(csvFileName)
	if err := convertJSONToCSV(casteData, csvFileName, post); err != nil {
		log.Fatal(err)
	}
}
