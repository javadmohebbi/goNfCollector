package common

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var ServicesFile = "/etc/services"

type Port struct {
	Proto string
	Name  string
}

func (p *Port) String() string {
	return fmt.Sprintf("{Proto:%v Name: %v}", p.Proto, p.Name)
}

type PortMap map[string]Port

func GetServices() (PortMap, error) {

	if runtime.GOOS == "windows" {
		ServicesFile = os.Getenv("systemroot") + "\\System32\\drivers\\etc\\services"
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		ServicesFile = "/etc/services"
	}

	file, err := os.Open(ServicesFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	port_map := make(PortMap)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// ignore comments
			continue
		}
		line = strings.TrimSpace(line)
		split := strings.SplitN(line, "#", 2)

		fields := strings.Fields(split[0])

		if len(fields) < 2 {
			continue
		}
		name := fields[0]
		portproto := strings.SplitN(fields[1], "/", 2)
		port, err := strconv.ParseInt(portproto[0], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		proto := strings.ToLower(portproto[1])
		port_map[proto+"/"+fmt.Sprintf("%v", port)] = Port{
			Name:  name,
			Proto: proto,
		}
	}
	return port_map, nil
}

// Return Protocol Name
func GetPortName(port string, protocol string) string {
	s, err := GetServices()
	if err != nil {
		return fmt.Sprintf("%v/%v", port, protocol)
	}
	protocol = strings.ToLower(protocol)
	for i, v := range s {
		if fmt.Sprintf("%v", i) == (protocol + "/" + port) {
			return v.Name
		}
	}
	return protocol + "/" + port
}
