package dns

import (
	"bufio"
	"errors"
	"net"
	"os"
	"strings"
)

type Sinkhole interface {
	ParseAdFile(string) error
}

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

		s.AddNewEntry(AEntry, split[1]+".", net.IPv4(0, 0, 0, 0))
	}

	return nil
}
