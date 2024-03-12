package main

import (
	"fmt"
	"github.com/bnb-chain/go-sdk/client/rpc"
	ctypes "github.com/bnb-chain/go-sdk/common/types"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const nodeAddr = "tcp://dataseed1.bnbchain.org:80"

func main() {
	result := getTokenBindStatus()
	fmt.Println(result)

	updateReadme(result)
}

func updateReadme(result string) {
	file, err := os.Open("README.md")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	original := string(content)

	var re = regexp.MustCompile(`/<!-- AUTO_UPDATE_START -->([\s\S]*?)<!-- AUTO_UPDATE_END -->/`)
	current := re.ReplaceAllString(original, fmt.Sprintf(`<!-- AUTO_UPDATE_START -->\n%s\n<!-- AUTO_UPDATE_END -->`, result))

	fmt.Println("Original", original)
	fmt.Println("Current", current)

	err = ioutil.WriteFile("README.md", []byte(current), 0644)
	if err != nil {
		panic(err)
	}
}

func getTokenBindStatus() string {
	client := rpc.NewRPCClient(nodeAddr, ctypes.ProdNetwork)
	tokens, err := client.ListAllTokens(0, 10000)
	if err != nil {
		panic(err)
	}
	result := ""
	for _, token := range tokens {
		if token.ContractAddress != "" && token.Symbol != "BNB" {
			splits := strings.Split(token.Symbol, "-")
			line := fmt.Sprintf("| %s | %s | %s | |\n", splits[0], token.Symbol, token.ContractAddress)
			result = result + line
		}
	}
	for _, token := range tokens {
		if token.ContractAddress == "" {
			splits := strings.Split(token.Symbol, "-")
			line := fmt.Sprintf("| %s | %s | | |\n", splits[0], token.Symbol)
			result = result + line
		}
	}
	return result
}
