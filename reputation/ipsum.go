package reputation

import (
	"encoding/csv"
	"errors"
	"net"
	"os"
	"strconv"

	"github.com/goNfCollector/common"
	"github.com/goNfCollector/configurations"
)

type IPSum struct {
	data [][]string
}

func NewIPSum(ipSumPath string) (IPSum, error) {

	file, err := os.Open(ipSumPath)
	if err != nil {
		return IPSum{}, errors.New(configurations.ERROR_READ_CONFIG.String())
	}
	defer file.Close()
	parser := csv.NewReader(file)
	parser.Comma = '\t'
	parser.Comment = '#'
	records, err := parser.ReadAll()
	if err != nil {
		return IPSum{}, errors.New(configurations.ERROR_READ_CONFIG.String())
	}

	return IPSum{
		data: records,
	}, nil
}

func (i *IPSum) Get(ipv4 string) ReputationResponse {

	ip := net.ParseIP(ipv4)
	if common.IsPrivateIP(ip) {
		return ReputationResponse{}
	}

	for _, record := range i.data {
		// ipv4 = "171.25.193.20"
		if record[0] == ipv4 {
			i, _ := strconv.Atoi(record[1])
			return ReputationResponse{
				Current: uint(i),
			}
		}
	}

	return ReputationResponse{Current: 0}
}

func (i *IPSum) GetType() string {
	return "ipsum"
}

func (i *IPSum) GetKind() string {
	return "malific:host:signature"
}
