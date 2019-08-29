package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
)

//generate csv
func main() {

	var idType = flag.String("id-type", "int", "can be int/string/uuid")
	var size = flag.Int("size", 1000000, "size of number of items")
	flag.Parse()
	file, err := os.Create("hierarchy.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"ID", "Nodeid", "Path"})
	if err != nil {
		panic(err)
	}

	level := 0
	parentPaths := []string{""}
	parentPathsIndex := 0
	childrenPaths := []string{}
	levelChildrenPaths := []string{}
	for i := 0; i < *size; i++ {
		//the number of nodes in hierarchy will be 2*level
		var record []string
		var nodePath string
		switch *idType {
		case "int":
			record, nodePath = getIDRecord(i, parentPaths[parentPathsIndex])
		case "string":
			panic("string unimplemented")
		case "uuid":
			panic("uuid unimplemented")
		default:
			fmt.Println("unknown id-type")
		}

		err := writer.Write(record)
		if err != nil {
			panic(err)
		}

		childrenPaths = append(childrenPaths, nodePath)

		if len(childrenPaths) >= level*2 {
			//next node

			levelChildrenPaths = append(levelChildrenPaths, childrenPaths...)

			childrenPaths = []string{}
			parentPathsIndex++
			if parentPathsIndex == 1 || parentPathsIndex > len(parentPaths) {
				level++
				parentPathsIndex = 0

				parentPaths = []string{}
				parentPaths = append(parentPaths, levelChildrenPaths...)
				levelChildrenPaths = []string{}
			}
		}
	}
}

func getIDRecord(i int, parentPath string) ([]string, string) {
	if parentPath == "" {
		return []string{fmt.Sprintf("%v", i), uuid.New().String(), fmt.Sprintf("%v", i)}, fmt.Sprintf("%v", i)
	}
	return []string{fmt.Sprintf("%v", i), uuid.New().String(), fmt.Sprintf("%v.%v", parentPath, i)}, fmt.Sprintf("%v.%v", parentPath, i)
}
