package network

import (
	"encoding/json"
	"log"
)

type Packet struct {
	FromId string
	DistId string
	PrevId string
	NextId string
	Data   string
}

func (p *Packet) Serialize() []byte {
	jsonData, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("error during packet serialization: %v", err)
	}
	return jsonData
}

func DeserializeFrom(data []byte) Packet {
	var packet Packet
	if err := json.Unmarshal(data, &packet); err != nil {
		log.Fatalf("error during packet deserialization: %v", err)
	}
	return packet
}
