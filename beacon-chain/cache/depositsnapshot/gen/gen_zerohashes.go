package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/prysmaticlabs/prysm/v3/beacon-chain/cache/depositsnapshot"
)

// tpl is the template for the generated go file.
var tpl = `// Code generated by gen_zerohashes. DO NOT EDIT.
package {{ .Package }}

var Zerohashes = [][32]byte{
{{- range .Zerohashes }}
	// {{ . | printf "%x" }}
	{ {{- . | join -}} },
{{- end }}
}
`

// join is used to format the zerohashes into a list.
func join(a [32]byte) string {
	var buff strings.Builder
	buff.Grow(32 * (2 + 3))
	buff.WriteString(fmt.Sprintf("%d", a[0]))
	for _, b := range a[1:] {
		buff.WriteString(fmt.Sprintf(", %d", b))
	}
	return buff.String()
}

// calculate the zerohashes as specified in the EIP and write them
// to a generated go file as a [][32]byte
func main() {
	zerohashes := make([][32]byte, 1, depositsnapshot.DepositContractDepth)
	fmt.Printf("%x\n", zerohashes[0])
	for i := 1; i < depositsnapshot.DepositContractDepth; i++ {
		zerohashes = append(zerohashes, sha256.Sum256(append(zerohashes[i-1][:], zerohashes[i-1][:]...)))
		fmt.Printf("%x\n", zerohashes[i])
	}

	tmpl, err := template.New("zerohashes").Funcs(map[string]any{"join": join}).Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("zerohashes.gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	err = tmpl.Execute(file, struct {
		Package    string
		Zerohashes [][32]byte
	}{
		Package:    "depositsnapshot",
		Zerohashes: zerohashes,
	})
	if err != nil {
		log.Fatal(err)
	}
}
