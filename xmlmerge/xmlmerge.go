package xmlmerge

import (
	"os"
	"encoding/xml"
	"io/ioutil"
	"fmt"
	"strconv"
	"time"
)

type Coverage struct {
	XMLName         xml.Name `xml:"coverage"`
	LineRate        string   `xml:"line-rate,attr"`
	BranchRate      string   `xml:"branch-rate,attr"`
	LinesCovered    string   `xml:"lines-covered,attr"`
	LinesValid      string   `xml:"lines-valid,attr"`
	BranchesCovered string   `xml:"branches-covered,attr"`
	BranchesValid   string   `xml:"branches-valid,attr"`
	Complexity      string   `xml:"complexity,attr"`
	Version         string   `xml:"version,attr"`
	Timestamp       string   `xml:"timestamp,attr"`
	Packages        Packages `xml:"packages"`
	Sources 		Sources  `xml:"sources"`
} 

type Sources struct{
	Source []string `xml:"source"`
}

type Packages struct {
	Package []Package `xml:"package"`
}

type Package struct{
	Name       string  `xml:"name,attr"`
	LineRate   string  `xml:"line-rate,attr"`
	BranchRate string  `xml:"branch-rate,attr"`
	Complexity string  `xml:"complexity,attr"`
	LineCount  string  `xml:"line-count,attr"`
	LineHits   string  `xml:"line-hits,attr"`
	Classes    Classes `xml:"classes"`
} 

type Classes struct{
	Class []Class `xml:"class"`
} 

type Class struct{
	Name       string `xml:"name,attr"`
	Filename   string `xml:"filename,attr"`
	LineRate   string `xml:"line-rate,attr"`
	BranchRate string `xml:"branch-rate,attr"`
	Complexity string `xml:"complexity,attr"`
	LineCount  string `xml:"line-count,attr"`
	LineHits   string `xml:"line-hits,attr"`
	Methods    Methods `xml:"methods"`
	Lines 	   Lines `xml:"lines"`
}

type Methods struct {
	Method []Method `xml:"method"`
} 

type Method struct{
	Name       string `xml:"name,attr"`
	Signature  string `xml:"signature,attr"`
	LineRate   string `xml:"line-rate,attr"`
	BranchRate string `xml:"branch-rate,attr"`
	Complexity string `xml:"complexity,attr"`
	LineCount  string `xml:"line-count,attr"`
	LineHits   string `xml:"line-hits,attr"`
	Lines      Lines  `xml:"lines"`
}

type Lines struct{
	Line []Line `xml:"line"` 
}

type Line struct{
	Number string `xml:"number,attr"`
	Hits   string `xml:"hits,attr"`
}

const xmlHeaderTag = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE coverage SYSTEM "http://cobertura.sourceforge.net/xml/coverage-04.dtd">
`

type XMLProcess interface{
	ProcessXMLs([]string) error
}

type XMLProcessor struct{

}

func (xp *XMLProcessor) ProcessXMLs(files []string) error {
	var coverData []Coverage
	for index := range files{
		var coverage Coverage
		err := convertToStruct(files[index],&coverage)
		if err != nil{
			return err
		}
		coverData = append(coverData,coverage)
	}
	covMergedData := mergeXMLStruct(coverData)
	mergedXML := convertToXML(covMergedData)
	fmt.Println(xmlHeaderTag+mergedXML)
	return nil
}

func convertToXML(mergedCoverage *Coverage) string {
	xmlFile , err := xml.MarshalIndent(mergedCoverage,"","\t")
	if err != nil{
		fmt.Println(" error found ",err)
	}
	return string(xmlFile)
}

func mergeXMLStruct(coverData []Coverage) *Coverage{
	mrgCover := Coverage{}
	branchesCovered := 0
	branchesValid := 0
	linesCovered := 0
	linesValid := 0
	branchRate := 0
	lineRate := 0
	complexity :=0
	for index := range coverData{
		cov := coverData[index]
		bc, _ := strconv.Atoi(cov.BranchesCovered)
		bv, _ := strconv.Atoi(cov.BranchesValid)
		lc, _ := strconv.Atoi(cov.LinesCovered)
		lv, _ := strconv.Atoi(cov.LinesValid)
		cmplx, _ := strconv.Atoi(cov.Complexity)
		mrgCover.Packages.Package = append(mrgCover.Packages.Package , cov.Packages.Package...)
		mrgCover.Sources.Source = append(mrgCover.Sources.Source , cov.Sources.Source...)
		branchesCovered += bc
		branchesValid += bv
		linesCovered += lc
		linesValid += lv
		if cmplx > complexity  {
			complexity = cmplx
		}
	}
	if branchesValid != 0 {
		branchRate = branchesCovered / branchesValid
	}
    if linesValid != 0 {
		lineRate = linesCovered / linesValid
	}
	mrgCover.BranchesCovered = strconv.Itoa(branchesCovered)
	mrgCover.BranchesValid = strconv.Itoa(branchesValid)
	mrgCover.LinesCovered = strconv.Itoa(linesCovered)
	mrgCover.LinesValid = strconv.Itoa(linesValid)
	mrgCover.BranchRate = strconv.Itoa(branchRate)
	mrgCover.LineRate = strconv.Itoa(lineRate)
	mrgCover.Complexity = strconv.Itoa(complexity)
	mrgCover.Timestamp = fmt.Sprintf("%d",time.Now().Unix())
	return &mrgCover
}

func convertToStruct(file string, cov *Coverage)error{
	xmlFile ,err := os.Open(file)

	if err != nil{
		return fmt.Errorf("error opening xml file %v",err)
	}
	defer xmlFile.Close()

	xmlByte, err := ioutil.ReadAll(xmlFile)
	if err != nil{
		return fmt.Errorf("error converting xml file to byte array %v",err)
	}

	err = xml.Unmarshal(xmlByte,cov)
	if err != nil{
		return fmt.Errorf("error unmarshing byte array to xml struct %v",err)
	}
	return nil
}

