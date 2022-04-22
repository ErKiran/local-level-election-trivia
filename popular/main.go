package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
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

type ElectionWardPrime struct {
	Localbodyname            string `json:"LocalBodyName"`
	District                 string `json:"District"`
	State                    int    `json:"State"`
	Postname                 string `json:"PostName"`
	WonBy                    string `json:"WonBy"`
	WinnerVote               int    `json:"WinnerVote"`
	WinningParty             string `json:"WinningParty"`
	HighestWardHeadVote      int    `json:"HighestWardHeadVote"`
	HighestWardHeadVoteParty string `json:"HighestWardHeadVoteParty"`
	VoteDifference           int    `json:"VoteDifference"`
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

func CountVote(wardData []ElectionData, localHeadData map[string]ElectionData, localBody string) ElectionWardPrime {
	var wardPartyCountMap = make(map[string]int)
	var eleCom ElectionWardPrime

	for _, data := range wardData {
		wardPartyCountMap[fmt.Sprintf("%s_%s", data.Localbodyname, data.Politicalpartyname)] += data.Totalvotesrecieved
	}

	maxNumber := math.MinInt32
	var pname string
	for key, value := range wardPartyCountMap {
		name := strings.Split(key, "_")
		if localBody == name[0] {
			if value > maxNumber {
				maxNumber = value
				pname = name[1]
			}
		}
	}

	if localHeadData[localBody].Politicalpartyname != pname {
		eleCom = ElectionWardPrime{
			Localbodyname:            localBody,
			District:                 localHeadData[localBody].Districtname,
			State:                    localHeadData[localBody].Stateid,
			WonBy:                    localHeadData[localBody].Candidatename,
			WinnerVote:               localHeadData[localBody].Totalvotesrecieved,
			WinningParty:             localHeadData[localBody].Politicalpartyname,
			HighestWardHeadVote:      maxNumber,
			HighestWardHeadVoteParty: pname,
			Postname:                 localHeadData[localBody].Postname,
			VoteDifference:           localHeadData[localBody].Totalvotesrecieved - maxNumber,
		}
	}
	return eleCom
}

func convertJSONToCSV(electionData []ElectionWardPrime, destination string) error {
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"S.No", "LocalBodyName", "Post", "District", "State", "Winner", "Winning Party", "Winner Vote", "Highest Ward Vote", "Highest Ward Vote Party"}
	if err := writer.Write(header); err != nil {
		return err
	}

	var count = 0
	for _, e := range electionData {
		count++
		var csvRow []string
		csvRow = append(csvRow, strconv.Itoa(count), e.Localbodyname, e.Postname, e.District, strconv.Itoa(e.State), e.WonBy, e.WinningParty, strconv.Itoa(e.WinnerVote), strconv.Itoa(e.HighestWardHeadVote), e.HighestWardHeadVoteParty)
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	electionData, err := ReadAndParseData()

	var electionCal []ElectionWardPrime

	if err != nil {
		fmt.Println(err)
	}

	var mainMap = make(map[string]ElectionData)
	var wardVoteMap = make(map[string][]ElectionData)
	post := "उपाध्यक्ष"

	for _, data := range electionData {
		if data.Postname == post && data.Estatus == "E" {
			mainMap[data.Localbodyname] = data
		}

		if data.Postname == "वडा अध्यक्ष" {
			wardVoteMap[data.Localbodyname] = append(wardVoteMap[data.Localbodyname], data)
		}
	}

	for key, value := range wardVoteMap {
		t := CountVote(value, mainMap, key)
		if t.Postname != "" {
			electionCal = append(electionCal, t)
		}
	}

	csvFileName := fmt.Sprintf("%s.csv", post)
	os.Create(csvFileName)
	if err := convertJSONToCSV(electionCal, csvFileName); err != nil {
		log.Fatal(err)
	}
}
