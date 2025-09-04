package utils

import "fmt"

// util para obtener el nombre del virtual id
func GetVirtualId(server string, virtualNode int) string {
	virtualId := fmt.Sprintf("%s:%d", server, virtualNode)
	return virtualId
}