package main

import (
	"fmt"
	"strconv"
	"io"
	"os"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func CreateCharts(a *Graph, b *Graph) {
	page := components.NewPage()
	page.Initialization.PageTitle = "Network diagrams"

	page.AddCharts(
		renderGraph1(a, "Workspace network topology"),
		renderGraph2(b, "Service network topology"),
	)

	f, err := os.Create("graph.html")
	if err != nil {
		panic(err)
	}

	page.Render(io.MultiWriter(f))
}


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





