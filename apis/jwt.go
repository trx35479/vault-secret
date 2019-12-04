package apis

import (
	"bufio"
	"log"
	"os"
)

func GetJwt(path string) ([]string, error) {

	// Path secrets usually in
	// "/var/run/secrets/kubernetes.io/serviceaccount/token"
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var text []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	return text, scanner.Err()
}
