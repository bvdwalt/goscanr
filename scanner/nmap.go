package scanner

import (
	"encoding/xml"
	"os/exec"
	"strconv"
	"strings"
)

type PortResult struct {
	IP      string
	Port    string
	Proto   string
	State   string
	Service string
}

type nmapXML struct {
	Hosts []nmapHost `xml:"host"`
}

type nmapHost struct {
	Address nmapAddress `xml:"address"`
	Ports   []nmapPort  `xml:"ports>port"`
}

type nmapAddress struct {
	Addr string `xml:"addr,attr"`
}

type nmapPort struct {
	Protocol string      `xml:"protocol,attr"`
	PortID   int         `xml:"portid,attr"`
	State    nmapState   `xml:"state"`
	Service  nmapService `xml:"service"`
}

type nmapState struct {
	State string `xml:"state,attr"`
}

type nmapService struct {
	Name string `xml:"name,attr"`
}

func NmapAvailable() bool {
	_, err := exec.LookPath("nmap")
	return err == nil
}

func buildNmapArgs(target string, ports []int) []string {
	portList := make([]string, len(ports))
	for i, p := range ports {
		portList[i] = strconv.Itoa(p)
	}
	return []string{"-oX", "-", "-p", strings.Join(portList, ","), target}
}

func parseNmapXML(data []byte) ([]PortResult, error) {
	var result nmapXML
	if err := xml.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	var portResults []PortResult
	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			portResults = append(portResults, PortResult{
				IP:      host.Address.Addr,
				Port:    strconv.Itoa(port.PortID),
				Proto:   port.Protocol,
				State:   port.State.State,
				Service: port.Service.Name,
			})
		}
	}
	return portResults, nil
}

func RunNmap(target string, ports []int) ([]PortResult, error) {
	args := buildNmapArgs(target, ports)
	cmd := exec.Command("nmap", args...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseNmapXML(out)
}
