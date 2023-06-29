package main

import (
	"fmt"
	"github.com/magiconair/properties"
	"io"
	"log"
	"os"
)

//func main() {
//	// Path to the property file
//	filePath := "device.properties"
//
//	// Read properties from the file
//	propertiesMap, err := readProperties(filePath)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if len(propertiesMap) > 0 {
//		fmt.Println("Properties read from the file:")
//		for key, value := range propertiesMap {
//			fmt.Printf("%s = %s\n", key, value)
//		}
//	}
//
//	propertiesMap["xytest"] = "true"
//	// Write properties to the file
//	err = writeProperties(filePath, propertiesMap)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	value, exist := propertiesMap["database.host"]
//	if exist {
//		fmt.Println(value)
//	} else {
//		fmt.Println("not foudn")
//	}
//}

func writeProperties(filePath string, propertiesMap map[string]string) error {
	// Create a new properties object
	p := properties.NewProperties()

	// Set properties
	for key, value := range propertiesMap {
		p.Set(key, value)
	}

	// Save properties to the file
	// Open the output file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Create an io.Writer from the file
	writer := io.Writer(file)

	_, err =  p.Write(writer, properties.UTF8)
	return err
}

func readProperties(filePath string) (map[string]string, error) {
	// Create a new properties object
	// p := properties.NewProperties()

	// Load properties from the file
	p, err := properties.LoadFile(filePath, properties.UTF8)
	//err := p.LoadFromFile(filePath, properties.UTF8)
	if err != nil {
		fmt.Println(err.Error())
		// still return an empty map if file not found
		return make(map[string]string), nil
	}

	// Get all properties as a map
	propertiesMap := p.Map()

	return propertiesMap, nil
}
