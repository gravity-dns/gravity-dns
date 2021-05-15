package dns

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

const addIP = "0.0.0.0"

func (s *gravityDNS) ParseAdFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		if text[0] == '#' {
			continue
		}
		split := strings.Split(text, " ")
		if len(split) != 2 {
			return errors.New("Invalid format " + text)
		}

		s.domains[split[1]] = split[0]
	}

	return nil
}

func IsAdIP(ip string) bool {
	return ip == addIP
}
