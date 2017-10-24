package main

import (
	"encoding/json"
	//"fmt"
)

type JsonNode struct {
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
}

type JsonEdge struct {
	Data struct {
		Id       string `json:"id"`
		Source   string `json:"source"`
		Target   string `json:"target"`
		Distance uint8  `json:"distance"`
	} `json:"data"`
}

func (g *Graph) toJSON() (string, error) {
	tree := make([]interface{}, 0)

	for _, node := range g.nodes {
		jnode := JsonNode{}
		jnode.Data.Id = node.Name
		//fmt.Printf("%v\n", jnode)
		tree = append(tree, jnode)
	}

	for _, edge := range g.edges {
		jedge := JsonEdge{}
		jedge.Data.Id = edge.Src.Name + "-" + edge.Dst.Name
		jedge.Data.Source = edge.Src.Name
		jedge.Data.Target = edge.Dst.Name
		jedge.Data.Distance = edge.Dist * 10
		tree = append(tree, jedge)
	}

	data, err := json.Marshal(tree)
	if err != nil {
		return "", err
	}

	//fmt.Printf("%v %v\n", tree, string(data))
	return string(data), nil
}
