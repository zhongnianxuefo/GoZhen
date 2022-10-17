package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

type DoOperator struct {
	Bracket int

	Priority int
	Value    string
	Operator string
	Value2   string
	NextDo   *DoOperator
	t1       int
	t2       int
	t3       int
	leftvar  bool
	Rightvar bool
}

func test3() {
	words := strings.Split("1 + 2 + 3 * 4 / 5 + 6 * 7 + 8", " ")
	var ops []DoOperator
	//blnFirst := true
	//var NowOperator *DoOperator
	//t := 0
	for i := 0; i < len(words)-1; i += 2 {
		Value := words[i]
		Operator := words[i+1]
		value2 := words[i+2]
		Bracket := 0
		Priority := 0
		switch Operator {
		case "+", "-":
			Priority = 10
		case "*", "/":
			Priority = 20
		}
		op := &DoOperator{Bracket: Bracket, Operator: Operator, Priority: Priority, Value: Value, Value2: value2}
		//op.t1 = t + 1
		//op.t2 = t + 2
		ops = append(ops, *op)
		//t = op.t1
	}
	var opss []DoOperator
	var addOp func(nowNo int)
	//tNo := make(map[int]int)
	tt := 1
	lastTT := 0
	leftvar := false
	Rightvar := false
	appendOp := func(nowOp DoOperator) {
		//t1 := nowOp.t1
		//t2 := nowOp.t2
		//t, ok := tNo[t1]
		//if ok {
		//	nowOp.t1 = t
		//}
		//t, ok = tNo[t2]
		//if ok {
		//	nowOp.t2 = t
		//}
		if tt > lastTT {
			nowOp.t1 = tt
			nowOp.t2 = tt + 1
			nowOp.t3 = tt
		} else {
			nowOp.t1 = tt
			nowOp.t2 = tt + 1
			nowOp.t3 = tt
		}

		nowOp.leftvar = leftvar

		nowOp.Rightvar = Rightvar

		lastTT = tt
		opss = append(opss, nowOp)
		//tNo[t2] = nowOp.t1
	}

	addOp = func(nowNo int) {
		nowOp := ops[nowNo]
		if nowNo < len(ops)-1 {

			nextOp := ops[nowNo+1]
			if nextOp.Priority > nowOp.Priority {
				tt += 1
				leftvar = false
				Rightvar = false
				addOp(nowNo + 1)

				tt -= 1
				//nowOp.t1 = t + 1
				//nowOp.t2 = t + 2
				//opss = append(opss, nowOp)
				Rightvar = true
				appendOp(nowOp)
				leftvar = true
				//t = nowOp.t1
			} else {
				//nowOp.t1 = t + 1
				//nowOp.t2 = t + 2
				Rightvar = false
				appendOp(nowOp)
				leftvar = true
				//t = nowOp.t1
				addOp(nowNo + 1)
			}
		} else {
			//nowOp.t1 = t + 1
			//nowOp.t2 = t + 2
			appendOp(nowOp)
			//t = nowOp.t1
		}

	}
	addOp(0)
	for n, op := range ops {
		fmt.Println("a", n, op)
	}
	fmt.Println(words)
	for n, op := range opss {
		fmt.Println("b", n, op)
	}

}

type S1 struct {
	SS    string
	par   interface{}
	Child interface{}
	Next  interface{}
}

type S2 struct {
	SS    string
	par   interface{}
	Child interface{}
	Next  interface{}
}
type S3 struct {
	SS    string
	par   interface{}
	Child interface{}
	Next  interface{}
}

func (s *S1) ToGobFile(gobFile string) (err error) {
	file, err := os.Create(gobFile)
	if err != nil {
		return
	}
	defer file.Close()
	gob.Register(S1{})
	gob.Register(S2{})
	gob.Register(S3{})
	enc := gob.NewEncoder(file)
	err = enc.Encode(s)
	if err != nil {
		return
	}
	return
}
func NewFromGobFile(gobFile string) (s *S1, err error) {
	file, err := os.Open(gobFile)
	if err != nil {
		return
	}
	defer file.Close()
	dec := gob.NewDecoder(file)
	s = &S1{}
	err = dec.Decode(&s)
	if err != nil {
		return
	}
	return

}
func test4(gobFile string) {
	sRoot := &S1{SS: "root"}
	s1 := &S1{SS: "1"}
	s2 := &S2{SS: "2"}
	s3 := &S3{SS: "3"}
	sRoot.Child = s1
	s1.par = sRoot
	s1.Next = s2
	s2.par = sRoot
	s2.Next = s3
	s3.par = sRoot

	err := sRoot.ToGobFile(gobFile)
	if err != nil {
		return
	}

	ss, err := NewFromGobFile(gobFile)
	if err != nil {
		return
	}
	fmt.Println(ss)
}
