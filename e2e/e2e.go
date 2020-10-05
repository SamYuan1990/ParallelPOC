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
		WKeyVersion: 2,
		Input:       0,
		Left:        nil,
		Right:       nil,
	}

	Read_A := &pipelinepoc.Node{
		RKey:        "a",
		RKeyVersion: 1,
		WKey:        "",
		WKeyVersion: 3,
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
	Tx = append(Tx, Read_A)
	Producer := &pipelinepoc.Producer{
		Nodes: make([]*pipelinepoc.Node, 0),
	}
	Producer.TreeMaking(Tx)
	ch := make(chan *pipelinepoc.Node, 10)
	done := make(chan bool, 2)

	Processor1 := &pipelinepoc.Processor{Name: "p1"}
	Processor2 := &pipelinepoc.Processor{Name: "p2"}
	go Processor1.Process(ch, done)
	go Processor2.Process(ch, done)
	go Producer.IntoChan(ch)
	<-done
	<-done
}
