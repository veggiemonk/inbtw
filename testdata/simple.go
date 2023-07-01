package server

import "fmt"

// [START bla]
var Bla = "bla"

// [END bla]

// [START morebla]

type PubArg struct {
	arg     []byte
	pacache []byte
	origin  []byte
	account []byte
	// [START inner]
	subject []byte
	deliver []byte
	mapped  []byte
	// [END inner]
	// should be printed
	// [START complex-2345 ]
	reply  []byte
	szb    []byte
	hdb    []byte
	queues [][]byte
	size   int
	hdr    int
}

// [END morebla]

// [START qwer]
const intertwine = 123456

// [END complex-2345 ]

// Oh is a sounds people make when they realize
// they have to maintain documentation
func Oh() error {
	// this is nice
	a := 5
	b := 8
	c := a + b*2
	fmt.Println("result", c)
	for i := range make([]string, 10) {
		// [START sup in]
		fmt.Println("i", i)
		// [END sup in]
	}
	return nil
}

// [END qwer]
