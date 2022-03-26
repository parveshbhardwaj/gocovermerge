package main

import (
	"fmt"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"gocovermerge/xmlmerge"
	"gocovermerge/jsonmerge"
)

func main(){
	xmlFlagPtr := flag.Bool("xml",false,"to merge coverage.xml for packages")
	jsonFlagPtr := flag.Bool("json",false,"to merge coverage.json for packages")
    flag.Parse()

	if *xmlFlagPtr {
		err := checkValidArguments(flag.Args(),".xml")
		if err != nil {
			panic(err)
		}
		xmlProcess := xmlmerge.XMLProcessor{}
		err = xmlProcess.ProcessXMLs(flag.Args())
		if err != nil {
			panic(err)
		}
	}else if *jsonFlagPtr {
		err := checkValidArguments(flag.Args(),".json")
		if err != nil {
			panic(err)
		}
		jsonProcess := jsonmerge.JsonProcessor{}
		err = jsonProcess.ProcessJSONs(flag.Args())
		if err != nil {
			panic(err)
		}
	}else{
		usage()
	}
}

func usage(){
	fmt.Println(`Usage : 
	For Json coverage merge : 
		gocovermerge -json coverage1.json coverage2.json > output.json
	For xml coverage merge : 
		gocovermerge -xml coverage1.xml coverage2.xml > output.xml`)
}

func checkValidArguments(arguments []string,ext string) error{
	for index :=  range arguments {
		file := arguments[index]
		fileInfo , err := os.Stat(file)
		if os.IsNotExist(err) {
			return fmt.Errorf("%v file does not exists",file)
		}
		if fileInfo.IsDir(){
			return fmt.Errorf("%v file is a directory",file)
		}
		if strings.ToLower(filepath.Ext(file)) != ext {
			return fmt.Errorf("%v is not %s file",file,ext)
		}
	}
 	return nil
}