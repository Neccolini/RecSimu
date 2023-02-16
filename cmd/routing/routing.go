package routing

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math"
	"math/big"
)

const (
	Router        = "Router"
	Coordinator   = "Coordinator"
	BroadCastId   = -1
	CoordinatorId = 0
)

type RF struct {
	id       int
	nodeType string
	table    map[int][]int
}
type Packet struct {
	id     int
	distId int
	data   string
}

func (r *RF) Init(nodeType string) ([][]byte, error) {
	r.nodeType = nodeType
	if nodeType == Router {
		// parent request 送信
		randomId, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
		if err != nil {
			log.Fatal(err)
		}
		r.id = int(randomId.Int64())

		p := Packet{r.id, BroadCastId, "preq"}
		return [][]byte{p.Serialize()}, nil
	} else {
		r.id = CoordinatorId
	}
	return nil, nil
}

func (r *RF) GenMessageFromM(received []byte) ([][]byte, error) {
	return nil, nil
}

func (r *RF) GenMessageFromI(inst []byte) ([][]byte, error) {
	return nil, nil
}

func (p *Packet) Serialize() []byte {
	jsonData, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("error during packet serialization: %v", err)
	}
	return jsonData
}
