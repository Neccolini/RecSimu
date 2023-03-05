package injection

import "fmt"

type Injection struct {
	InjectionId string `json:"injectionId"`
	Cycle       int    `json:"cycle"`
	FromId      string `json:"fromId"`
	DistId      string `json:"distId"`
	Data        string `json:"data"`
}

type InjectionTable struct {
	InjectionMap map[string]Injection
	CycleQuery   map[int][]string
}

func (i *InjectionTable) QueryByCycle(cycle int) ([]Injection, error) {
	ids, ok := i.CycleQuery[cycle]
	if !ok {
		return nil, fmt.Errorf("injection(cycle = %d) is not found", cycle)
	}
	res := []Injection{}
	for _, id := range ids {
		res = append(res, i.InjectionMap[id])
	}
	return res, nil
}

func (i *InjectionTable) QueryById(id string) (Injection, error) {
	if injection, ok := i.InjectionMap[id]; ok {
		return injection, nil
	}
	return Injection{}, fmt.Errorf("injection(id = %s) is not found", id)
}
