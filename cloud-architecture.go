package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// With this once we have this working and that we can simulate a network, then I want to be able to create a graph of this with a visualisation. Do this in typescript.

type Node struct {
	Id              string
	Name            string
	Description     string
	OutNode         []*Node
	InNode          []*Node
	ContainedByNode *Node // a node can only be contained by one other node
	ContainNode     []*Node
}

type Graph struct {
	nodes map[int]*Node
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[int]*Node),
	}
}

func (g *Graph) AddNode(nodeID int, name string) {
	if _, exists := g.nodes[nodeID]; !exists {
		newNode := &Node{
			Name: name,
		}
		g.nodes[nodeID] = newNode
		fmt.Println("New node added to graph")
	} else {
		fmt.Println("Node already exists!")
	}
}

func (g *Graph) AddEdge(from, to int) {
	from_node := g.nodes[from]
	to_node := g.nodes[to]

	from_node.OutNode = append(from_node.OutNode, to_node)
	to_node.InNode = append(to_node.InNode, from_node)
}

func (g *Graph) AddContainedEdge(contained int, container_graph *Graph, container int) {
	contained_node := g.nodes[contained]
	container_node := container_graph.nodes[container]

	contained_node.ContainedByNode = container_node
	container_node.ContainNode = append(container_node.ContainNode, contained_node)
}

// get the key in the map of a node, useful for getting the key that a node in the OutNode or InNode is pointing to
func getNodeId(g *Graph, node *Node) int {

	for key, val := range g.nodes {
		if node == val {
			return key
		}
	}
	return -99
}

// implement the deletion of edges and nodes later as we don't need that at the momement

func (g *Graph) PrintGraph() {

	// we should really sort the map before printing it, this is where implementing a bubblesort would be useful

	for key, n := range g.nodes {
		fmt.Print("(", key, ") ", n.Name, " -> ")
		for _, in := range n.OutNode {
			// get the key in the graph map that this outnode is pointing to
			//	for key, val := range g.nodes {
			//	r = if
			fmt.Print("(", getNodeId(g, in), ") ", in.Name, "; ")
		}
		fmt.Println()
		fmt.Print("(", key, ") ", n.Name, " <- ")
		for _, out := range n.InNode {
			fmt.Print("(", getNodeId(g, out), ") ", out.Name, "; ")
		}
		fmt.Println()
	}
}

func renderGraph1(g *Graph, chart_title string) *charts.Graph {

	var nodes []opts.GraphNode
	var edges []opts.GraphLink

	for key, n := range g.nodes {
		nodes = append(nodes, opts.GraphNode{
			Name:       strconv.Itoa(key) + " " + n.Name,
			Symbol:     "circle",
			SymbolSize: []interface{}{40, 40}})

	}

	for key, n := range g.nodes {
		for _, out := range n.OutNode {
			edges = append(edges, opts.GraphLink{
				Source: strconv.Itoa(key) + " " + n.Name,
				Target: strconv.Itoa(getNodeId(g, out)) + " " + out.Name})
		}
	}

	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: chart_title,
			}),
		charts.WithLegendOpts(
			opts.Legend{
				Show: true,
			}),
		charts.WithInitializationOpts(
			opts.Initialization{
				Width:  "1000px",
				Height: "800px",
			}),
	)
	graph.AddSeries("data", nodes, edges,
		charts.WithGraphChartOpts(
			opts.GraphChart{
				Layout:         "force",
				Roam:           true,
				Draggable:      true,
				Force:          &opts.GraphForce{Repulsion: 4000, Gravity: 0.5},
				EdgeSymbol:     []interface{}{"none", "arrow"},
				EdgeSymbolSize: 10,
			}),
		charts.WithLabelOpts(
			opts.Label{
				Show:          true,
				Color:         "black",
				Position:      "below",
				FontWeight:    "normal",
				Align:         "right",
				VerticalAlign: "bottom",
			}),
		charts.WithLineStyleOpts(
			opts.LineStyle{
				Opacity:   1,
				Curveness: 0,
				Width:     1,
				Color:     "gray",
			}))
	return graph
}

func renderGraph2(g *Graph, chart_title string) *charts.Graph {

	var nodes []opts.GraphNode
	var edges []opts.GraphLink
	var series []string
	colours := [...]string{"#ff595e", "#ff924c", "#ffca3a", "#c5ca30", "#8ac926", "#52a675", "#1982c4", "#4267ac", "#6a4c93", "#b5a6c9"}
	for _, n := range g.nodes {
		contains := false
		for _, sn := range series {
			if n.ContainedByNode.Name == sn {
				contains = true
			}
		}
		if contains == false {
			series = append(series, n.ContainedByNode.Name)
		}
	}

	var col string

	for key, n := range g.nodes {
		for i, sn := range series {
			if n.ContainedByNode.Name == sn {
				col = colours[i]
				nodes = append(nodes, opts.GraphNode{
					Name:       strconv.Itoa(key) + " " + n.Name,
					ItemStyle:  &opts.ItemStyle{Color: col},
					Symbol:     "circle",
					SymbolSize: []interface{}{40, 40}})
			}
		}
	}

	for key, n := range g.nodes {
		for _, out := range n.OutNode {
			edges = append(edges, opts.GraphLink{
				Source: strconv.Itoa(key) + " " + n.Name,
				Target: strconv.Itoa(getNodeId(g, out)) + " " + out.Name})
		}
	}

	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: chart_title,
			}),
		charts.WithLegendOpts(
			opts.Legend{
				Show: true,
				Data: []interface{}{"Foo", "Bar"},
			}),
		charts.WithInitializationOpts(
			opts.Initialization{
				Width:  "1000px",
				Height: "800px",
			}),
	)
	graph.AddSeries("data", nodes, edges,
		charts.WithGraphChartOpts(
			opts.GraphChart{
				Layout:         "force",
				Roam:           true,
				Draggable:      true,
				Force:          &opts.GraphForce{Repulsion: 4000, Gravity: 0.5},
				EdgeSymbol:     []interface{}{"none", "arrow"},
				EdgeSymbolSize: 10,
			}),
		charts.WithLabelOpts(
			opts.Label{
				Show:          true,
				Color:         "black",
				Position:      "below",
				FontWeight:    "normal",
				Align:         "right",
				VerticalAlign: "bottom",
			}),
		charts.WithLineStyleOpts(
			opts.LineStyle{
				Opacity:   1,
				Curveness: 0,
				Width:     1,
				Color:     "gray",
			}))
	return graph
}

func main() {

	// define our workspaces
	wgraph := NewGraph()

	wgraph.AddNode(1, "CMS")
	wgraph.AddNode(2, "Finance")
	wgraph.AddNode(3, "Student")
	wgraph.AddNode(4, "Patient records")
	wgraph.AddNode(5, "Enterprise data platform")
	wgraph.AddNode(6, "Secured data platform")
	wgraph.AddNode(7, "Reporting")
	wgraph.AddNode(8, "TRE Workspace 1")
	wgraph.AddNode(9, "TRE Workspace 2")
	wgraph.AddNode(10, "TRE Workspace 3")

	wgraph.AddEdge(1, 5)
	wgraph.AddEdge(2, 5)
	wgraph.AddEdge(3, 5)
	wgraph.AddEdge(4, 6)
	wgraph.AddEdge(5, 7)
	wgraph.AddEdge(5, 6)
	wgraph.AddEdge(6, 8)
	wgraph.AddEdge(6, 9)
	wgraph.AddEdge(6, 10)

	fmt.Println("Workspace network topology")
	wgraph.PrintGraph()

	// define our services
	sgraph := NewGraph()

	sgraph.AddNode(1, "Dynamics")
	sgraph.AddNode(2, "Dynamics")
	sgraph.AddNode(3, "Dynamics")
	sgraph.AddNode(4, "Virtual machine")
	sgraph.AddNode(5, "ADF")
	sgraph.AddNode(6, "Databricks")
	sgraph.AddNode(7, "Datalake")
	sgraph.AddNode(8, "ADF")
	sgraph.AddNode(9, "Databricks")
	sgraph.AddNode(10, "Datalake")
	sgraph.AddNode(11, "Power BI")
	sgraph.AddNode(12, "Databricks")
	sgraph.AddNode(13, "Virtual machine")
	sgraph.AddNode(14, "Virtual machine")
	sgraph.AddNode(15, "Azure ML")
	sgraph.AddNode(16, "Virtual machine")
	sgraph.AddNode(17, "Function App")

	sgraph.AddEdge(1, 5)
	sgraph.AddEdge(2, 5)
	sgraph.AddEdge(3, 5)

	sgraph.AddEdge(4, 8)

	sgraph.AddEdge(5, 6)
	sgraph.AddEdge(6, 7)

	sgraph.AddEdge(7, 8)
	sgraph.AddEdge(8, 9)
	sgraph.AddEdge(9, 10)

	sgraph.AddEdge(7, 11)

	sgraph.AddEdge(10, 12)
	sgraph.AddEdge(10, 13)
	sgraph.AddEdge(10, 14)
	sgraph.AddEdge(10, 15)
	sgraph.AddEdge(10, 16)
	sgraph.AddEdge(16, 17)

	sgraph.AddEdge(12, 13)
	sgraph.AddEdge(13, 12)
	sgraph.AddEdge(14, 15)
	sgraph.AddEdge(15, 14)

	sgraph.AddContainedEdge(1, wgraph, 1)
	sgraph.AddContainedEdge(2, wgraph, 2)
	sgraph.AddContainedEdge(3, wgraph, 3)
	sgraph.AddContainedEdge(4, wgraph, 4)
	sgraph.AddContainedEdge(5, wgraph, 5)
	sgraph.AddContainedEdge(6, wgraph, 5)
	sgraph.AddContainedEdge(7, wgraph, 5)
	sgraph.AddContainedEdge(8, wgraph, 6)
	sgraph.AddContainedEdge(9, wgraph, 6)
	sgraph.AddContainedEdge(10, wgraph, 6)
	sgraph.AddContainedEdge(11, wgraph, 7)
	sgraph.AddContainedEdge(12, wgraph, 8)
	sgraph.AddContainedEdge(13, wgraph, 8)
	sgraph.AddContainedEdge(14, wgraph, 9)
	sgraph.AddContainedEdge(15, wgraph, 9)
	sgraph.AddContainedEdge(16, wgraph, 10)
	sgraph.AddContainedEdge(17, wgraph, 10)

	fmt.Println("Service network topology")
	sgraph.PrintGraph()

	page := components.NewPage()
	page.Initialization.PageTitle = "Network diagrams"

	page.AddCharts(
		renderGraph1(wgraph, "Workspace network topology"),
		renderGraph2(sgraph, "Service network topology"),
	)

	f, err := os.Create("graph.html")
	if err != nil {
		panic(err)
	}

	page.Render(io.MultiWriter(f))
}
