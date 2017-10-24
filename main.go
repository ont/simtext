package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
)

var (
	debugFlag     = flag.Bool("debug", false, "Run debug mode")
	cleanFlag     = flag.Bool("clean", false, "Clear all files from target dir and save them with .clean suffix")
	cytoscapeFlag = flag.Bool("cytoscape", false, "Output for cytoscape.js")
	gephiFlag     = flag.String("gephi", "", "Path for CSV files for Gephi to output to.")
	threshFlag    = flag.Int("thresh", 0, "Threshold for groupping")
	usageFunc     = func() {
		fmt.Printf("Usage: %s [OPTIONS] directory\n", os.Args[0])
		flag.PrintDefaults()
	}
)

func main() {
	flag.Usage = usageFunc
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	files, err := readHashedFiles(flag.Arg(0))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	println(*cleanFlag)
	switch {
	case *debugFlag:
		debugFiles(files)

	case *cleanFlag:
		cleanFiles(files) // TODO: not optimal (clean + hash + final clean)

	case *cytoscapeFlag:
		cytoscapeOutput(files)

	case *gephiFlag != "":
		outputGephiCSV(*gephiFlag, files)

	default:
		flag.Usage()
		os.Exit(1)

	}
}

func readHashedFiles(dir string) ([]*HashedFile, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	hashedFiles := make([]*HashedFile, 0)

	for _, file := range files {
		hf, err := NewHashedFile(dir + "/" + file.Name())
		if err != nil {
			return nil, err
		}
		hashedFiles = append(hashedFiles, hf)

		//if i > 500 {
		//	break
		//}
	}

	return hashedFiles, nil
}

func outputGephiCSV(path string, files []*HashedFile) {
	graph := buildGraph(files)

	outputNodes(path, graph)
	outputEdges(path, graph)

}

func outputNodes(path string, graph *Graph) {
	file, err := os.Create(path + "/graph_nodes.csv")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	cw := csv.NewWriter(file)
	defer cw.Flush()

	cw.Write([]string{"Id", "Path", "Hash"})
	for _, node := range graph.nodes {
		cw.Write([]string{node.Name, node.path, fmt.Sprintf("%064b", node.Hash)})
	}
}

func outputEdges(path string, graph *Graph) {
	file, err := os.Create(path + "/graph_edges.csv")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	cw := csv.NewWriter(file)
	defer cw.Flush()

	cw.Write([]string{"Source", "Target", "Id", "Weight"})
	for _, edge := range graph.edges {
		cw.Write([]string{edge.Src.Name, edge.Dst.Name, edge.Src.Name + "-" + edge.Dst.Name, fmt.Sprintf("%f", math.Exp(-float64(edge.Dist)/10)*20)})
	}
}

func cytoscapeOutput(files []*HashedFile) {
	graph := buildGraph(files)

	json, err := graph.toJSON()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(json)
}

func buildGraph(files []*HashedFile) *Graph {
	graph := NewGraph()

	for _, file := range files {
		graph.AddNode(file)
	}
	graph.BuildEdges()

	return graph
}

func debugFiles(files []*HashedFile) {
	//for _, file := range files {
	//	data, err := ioutil.ReadFile(dir + "/" + file.Name())

	//	if err != nil {
	//		fmt.Println(err.Error())
	//		os.Exit(1)
	//	}

	//	fmt.Printf("------------\n")
	//	hash := calcHash(data, true)

	//	fmt.Printf("%s --> %s\n", file.Name(), fmt.Sprintf("%064b", hash))
	//}
}

func cleanFiles(files []*HashedFile) {
	println("here")
	for _, file := range files {
		text := file.GetCleanText()

		ext := filepath.Ext(file.path)
		path := strings.TrimSuffix(file.path, ext)

		println(path + ".clean" + ext)
		err := ioutil.WriteFile(path+".clean"+ext, []byte(text), 0644)
		if err != nil {
			panic(err) // TODO: normal err handling
		}
	}
}
