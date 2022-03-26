package jsonmerge

import(
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Coverage struct {
	Packages []Packages `json:"Packages"`
}

type Packages struct {
	Name      string `json:"Name"`
	Functions []Functions `json:"Functions"`
}

type Functions struct{
	Name       string `json:"Name"`
	File       string `json:"File"`
	Start      int    `json:"Start"`
	End        int    `json:"End"`
	Statements []Statements `json:"Statements"`
}

type Statements struct {
	Start   int `json:"Start"`
	End     int `json:"End"`
	Reached int `json:"Reached"`
}


type JsonProcess interface{
	ProcessJSONs([]string) error
}

type JsonProcessor struct{

}

func (jp *JsonProcessor) ProcessJSONs(files []string) error {
	var coverData []Coverage
	for index := range files{
		var coverage Coverage
		err := convertToStruct(files[index],&coverage)
		if err != nil{
			return err
		}
		coverData = append(coverData,coverage)
	}
	covMergedData := mergeJSONStruct(coverData)
	mergedJSON := convertToJSON(covMergedData)
	fmt.Println(mergedJSON)
	return nil
}

func convertToJSON(mergedCoverage *Coverage) string {
	jsonFile , err := json.MarshalIndent(mergedCoverage,"","\t")
	if err != nil{
		fmt.Println(" error found ",err)
	}
	return string(jsonFile)
}

func mergeJSONStruct(coverData []Coverage) *Coverage{
	mrgCover := Coverage{}
	for index := range coverData{
		cov := coverData[index]
		mrgCover.Packages = append(mrgCover.Packages,cov.Packages...)
	}
	return &mrgCover
}

func convertToStruct(file string, cov *Coverage)error{
	jsonFile ,err := os.Open(file)

	if err != nil{
		return fmt.Errorf("error opening json file %v",err)
	}
	defer jsonFile.Close()

	jsonByte, err := ioutil.ReadAll(jsonFile)

	if err != nil{
		return fmt.Errorf("error converting json file to byte array %v",err)
	}
	
	err = json.Unmarshal(jsonByte,cov)
	if err != nil{
		return fmt.Errorf("error unmarshing byte array to json struct %v",err)
	}

	return nil
}

