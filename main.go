package main

import (
	"fmt"
	"net/http"

	"github.com/unbot2313/distributed-cache/pkg/hash"
	"github.com/unbot2313/distributed-cache/pkg/ring"
)

func main() {
	fmt.Println("Hello World!")

	hasher := hash.NewXXH3Hasher()

	// aprox con 100 nodos tendra una desviacion de 10% como maximo, con 200 un 5%
	r := ring.NewRing(hasher, 100)

	r.AddNode("server:1")
	r.AddNode("server:3")
	r.AddNode("server:2")

	http.ListenAndServe(":3211", nil)

}
