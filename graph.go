package main

import (
	"github.com/mfonda/simhash"
)

type Graph struct {
	nodes []*HashedFile
	edges []*Edge
}

type Edge struct {
	Src *HashedFile
	Dst *HashedFile

	Dist uint8
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make([]*HashedFile, 0),
		edges: make([]*Edge, 0),
	}
}

func (g *Graph) AddNode(node *HashedFile) {
	g.nodes = append(g.nodes, node)
}

func (g *Graph) BuildEdges() {
	for i, nodeS := range g.nodes[:len(g.nodes)-1] {
		minDist := g.dist(nodeS, g.nodes[i+1])
		minNode := g.nodes[i+1]

		for _, nodeD := range g.nodes[i+1:] {
			dist := g.dist(nodeS, nodeD)
			if dist < minDist {
				minDist = dist
				minNode = nodeD
			}
		}

		g.edges = append(g.edges, &Edge{
			Src:  nodeS,
			Dst:  minNode,
			Dist: minDist,
		})
	}
}

func (g *Graph) dist(node1 *HashedFile, node2 *HashedFile) uint8 {
	return simhash.Compare(node1.Hash, node2.Hash)
}
