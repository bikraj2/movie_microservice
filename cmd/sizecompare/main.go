package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"bikraj.movie_microservice.net/gen"
	model "bikraj.movie_microservice.net/metadata/pkg"
	"github.com/gogo/protobuf/proto"
)

var metadata = &model.Metadata{
	ID:          "123",
	Title:       "The movie 2",
	Description: "A great movei",
	Director:    "Nrisopher cholan",
}

var genMetadata = &gen.Metadata{
	Id:          "123",
	Title:       "The movie 2",
	Description: "A great movei",
	Director:    "Nrisopher cholan",
}

func main() {
	jsonBytes, err := json.Marshal(metadata)
	if err != nil {
		panic(err)
	}
	xmlBytes, err := xml.Marshal(metadata)
	if err != nil {
		panic(err)
	}
	protoBytes, err := proto.Marshal(genMetadata)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Size of json is : %v\n", len(jsonBytes))
	fmt.Printf("Size of xml is : %v\n", len(xmlBytes))
	fmt.Printf("Size of proto is : %v\n", len(protoBytes))
}

func serializeToJSON(m *model.Metadata) ([]byte, error) {
	return json.Marshal(m)
}

func seralizeToXml(m *model.Metadata) ([]byte, error) {
	return xml.Marshal(m)
}

func seralizeToProto(m *gen.Metadata) ([]byte, error) {
	return proto.Marshal(m)
}
