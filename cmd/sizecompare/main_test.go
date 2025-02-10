package main

import "testing"

func BenchmarkSerializeToJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		serializeToJSON(metadata)
	}
}

func BenchmarkSerializeTOXML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		seralizeToXml(metadata)
	}
}

func BenchmarkSerializeToProto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		seralizeToProto(genMetadata)
	}
}
