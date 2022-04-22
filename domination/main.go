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

type ElectionCalculation struct {
	Localbodyname          string `json:"LocalBodyName"`
	GovernmentAllianceVote int    `json:"GovernmentAllianceVote"`
	District               string `json:"District"`
	State                  int    `json:"State"`
	WonBy                  string `json:"WonBy"`
	Postname               string `json:"PostName"`
	UMLVote                int    `json:"UMLVote"`
}

type ElectionCompetation struct {
	Localbodyname     string `json:"LocalBodyName"`
	WinnerVote        int    `json:"WinnerVote"`
	LoserVote         int    `json:"LoserVote"`
	VoteDifference    int    `json:"VoteDifference"`
	District          string `json:"District"`
	State             int    `json:"State"`
	WonBy             string `json:"WonBy"`
	WinningParty      string `json:"WinningParty"`
	SecondParty       string `json:"SecondParty"`
	NearestCompetitor string `json:"NearestCompetitor"`
	Postname          string `json:"PostName"`
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

func CountVote(election []ElectionData, localBody string) ElectionCompetation {
	var eleCom ElectionCompetation
	var winVote, losVote, state, diff int
	var district, winner, nearset, winningParty, secondParty, postname string
	for _, data := range election {
		state = data.Stateid
		district = data.Districtname
		postname = data.Postname

		if data.Rank == 1 {
			winner = data.Candidatename
			winVote = data.Totalvotesrecieved
			winningParty = data.Politicalpartyname
		}

		if data.Rank != 1 {
			nearset = data.Candidatename
			losVote += data.Totalvotesrecieved
			secondParty = data.Politicalpartyname
		}
	}
	fmt.Println("losVote", losVote, "winVote", winVote)

	if winVote > losVote {
		diff = winVote - losVote
		eleCom = ElectionCompetation{
			Localbodyname:     localBody,
			State:             state,
			WinnerVote:        winVote,
			LoserVote:         losVote,
			District:          district,
			WonBy:             winner,
			NearestCompetitor: nearset,
			Postname:          postname,
			WinningParty:      winningParty,
			SecondParty:       secondParty,
			VoteDifference:    diff,
		}
	}
	return eleCom
}

func convertJSONToCSV(electionData []ElectionCompetation, destination string) error {
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"S.No", "LocalBodyName", "Post", "District", "State", "Winner", "Winning Party", "Winner Vote", "Opposition Vote", "VoteDifference"}
	if err := writer.Write(header); err != nil {
		return err
	}

	var count = 0
	for _, e := range electionData {
		count++
		var csvRow []string
		csvRow = append(csvRow, strconv.Itoa(count), e.Localbodyname, e.Postname, e.District, strconv.Itoa(e.State), e.WonBy, e.WinningParty, strconv.Itoa(e.WinnerVote), strconv.Itoa(e.LoserVote), strconv.Itoa(e.VoteDifference))
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	electionData, err := ReadAndParseData()

	var electionCal []ElectionCompetation

	if err != nil {
		fmt.Println(err)
	}

	var mainMap = make(map[string][]ElectionData)
	// for _, post := range []string{"प्रमुख", "उपप्रमुख", "अध्यक्ष", "उपाध्यक्ष"} {
	post := "उपाध्यक्ष"

	for _, data := range electionData {
		if data.Postname == post {
			mainMap[data.Localbodyname] = append(mainMap[data.Localbodyname], data)
		}
	}

	for key, value := range mainMap {
		t := CountVote(value, key)
		if t.Localbodyname != "" {
			electionCal = append(electionCal, t)
		}
	}
	// sort.SliceStable(electionCal, func(i, j int) bool {
	// 	return electionCal[i].VoteDifference > electionCal[j].VoteDifference
	// })
	csvFileName := fmt.Sprintf("%s.csv", post)
	os.Create(csvFileName)
	if err := convertJSONToCSV(electionCal, csvFileName); err != nil {
		log.Fatal(err)
	}
}
