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

func CountVote(election []ElectionData, localBody string) ElectionCalculation {
	var electionCal ElectionCalculation
	partyMap := map[int]string{
		1: "नेपाली काँग्रेस",
		2: "नेपाल कम्युनिष्ट पार्टी (एकीकृत मार्क्सवादी-लेनिनवादी)",
		3: "नेपाल कम्युनिष्ट पार्टी (माओवादी केन्द्र)",
		4: "राष्ट्रिय जनमोर्चा ",
		5: "संघीय समाजवादी फोरम, नेपाल ",
		6: "नयाँ शक्ति पार्टी, नेपाल",
		7: "राष्ट्रिय जनता पार्टी नेपाल",
	}
	var gov, opp, state int
	var district, winner, postname string
	for _, data := range election {
		state = data.Stateid
		district = data.Districtname
		postname = data.Postname
		if data.Rank == 1 {
			winner = data.Candidatename
		}
		if partyMap[2] == data.Politicalpartyname {
			opp = data.Totalvotesrecieved
		}

		if partyMap[1] == data.Politicalpartyname || partyMap[3] == data.Politicalpartyname || partyMap[4] == data.Politicalpartyname || partyMap[5] == data.Politicalpartyname || partyMap[6] == data.Politicalpartyname || partyMap[7] == data.Politicalpartyname {
			gov += data.Totalvotesrecieved
		}
	}
	if opp > gov {
		electionCal = ElectionCalculation{
			Localbodyname:          localBody,
			GovernmentAllianceVote: gov,
			UMLVote:                opp,
			State:                  state,
			District:               district,
			WonBy:                  winner,
			Postname:               postname,
		}
	}

	return electionCal

}

func convertJSONToCSV(electionData []ElectionCalculation, destination string) error {
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"S.No", "LocalBodyName", "Post", "District", "State", "Winner", "GovernmentAllianceVote", "UMLVote"}
	if err := writer.Write(header); err != nil {
		return err
	}

	var count = 0
	for _, e := range electionData {
		count++
		var csvRow []string
		csvRow = append(csvRow, strconv.Itoa(count), e.Localbodyname, e.Postname, e.District, strconv.Itoa(e.State), e.WonBy, strconv.Itoa(e.GovernmentAllianceVote), strconv.Itoa(e.UMLVote))
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	electionData, err := ReadAndParseData()

	var electionCal []ElectionCalculation

	if err != nil {
		fmt.Println(err)
	}

	var mainMap = make(map[string][]ElectionData)

	for _, post := range []string{"प्रमुख", "उपप्रमुख", "उपाध्यक्ष", "अध्यक्ष"} {
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
		csvFileName := fmt.Sprintf("%s.csv", post)
		os.Create(csvFileName)
		if err := convertJSONToCSV(electionCal, csvFileName); err != nil {
			log.Fatal(err)
		}
	}
}
