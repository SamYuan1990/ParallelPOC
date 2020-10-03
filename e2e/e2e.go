package main

import "github.com/SamYuan1990/pipelinepoc"

func main() {
	Create_A := &pipelinepoc.Node{
		RKey:        "",
		RKeyVersion: -1,
		WKey:        "a",
		WKeyVersion: 0,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}

	Create_B := &pipelinepoc.Node{
		RKey:        "",
		RKeyVersion: -1,
		WKey:        "b",
		WKeyVersion: 0,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}

	Update_A := &pipelinepoc.Node{
		RKey:        "",
		RKeyVersion: -1,
		WKey:        "a",
		WKeyVersion: 1,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}

	Create_C := &pipelinepoc.Node{
		RKey:        "",
		RKeyVersion: -1,
		WKey:        "c",
		WKeyVersion: 0,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}

	Create_D := &pipelinepoc.Node{
		RKey:        "",
		RKeyVersion: -1,
		WKey:        "d",
		WKeyVersion: 0,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}
	RA_WD := &pipelinepoc.Node{
		RKey:        "a",
		RKeyVersion: 1,
		WKey:        "d",
		WKeyVersion: 1,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}
	Tx := make([]*pipelinepoc.Node, 0)
	Tx = append(Tx, Create_A)
	Tx = append(Tx, Create_B)
	Tx = append(Tx, Update_A)

	Tx = append(Tx, Create_C)
	Tx = append(Tx, Create_D)
	Tx = append(Tx, RA_WD)
	consumer := &pipelinepoc.Consumer{
		Nodes: make([]*pipelinepoc.Node, 0),
	}
	consumer.TreeMaking(Tx)
	ch := make(chan *pipelinepoc.Node, 10)
	consumer.IntoChan(ch)

	Processor1 := &pipelinepoc.Processor{Name: "p1"}
	Processor2 := &pipelinepoc.Processor{Name: "p2"}
	go Processor1.Process(ch)
	Processor2.Process(ch)
	defer close(ch)
}
