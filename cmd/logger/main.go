/*
Command nf-dump decodes NetFlow packets from UDP datagrams.

Usage:
		nf-dump [flags]

Flags:
		-addr string 	Listen address (default ":2055")
*/
package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/tehmaze/netflow"
	"github.com/tehmaze/netflow/ipfix"
	"github.com/tehmaze/netflow/netflow1"
	"github.com/tehmaze/netflow/netflow5"
	"github.com/tehmaze/netflow/netflow6"
	"github.com/tehmaze/netflow/netflow7"
	"github.com/tehmaze/netflow/netflow9"
	"github.com/tehmaze/netflow/read"
	"github.com/tehmaze/netflow/session"
)

// Safe default
var readSize = 2 << 16

func main() {
	listen := flag.String("addr", ":6859", "Listen address")

	logPath := flag.String("out", "/tmp/nfcollector-dump.log", "Path to collected logs")

	flag.Parse()

	var err error

	log.Println("Creating log file in this path: ", *logPath)
	f, err := os.Create(*logPath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var addr *net.UDPAddr
	if addr, err = net.ResolveUDPAddr("udp", *listen); err != nil {
		log.Fatal(err)
	}

	var server *net.UDPConn
	if server, err = net.ListenUDP("udp", addr); err != nil {
		log.Fatal(err)
	}

	if err = server.SetReadBuffer(readSize); err != nil {
		log.Fatal(err)
	}

	decoders := make(map[string]*netflow.Decoder)

	for {
		buf := make([]byte, 8192)
		var remote *net.UDPAddr
		var octets int
		if octets, remote, err = server.ReadFromUDP(buf); err != nil {
			log.Printf("error reading from %s: %v\n", remote, err)
			continue
		}

		log.Printf("received %d bytes from %s\n", octets, remote)

		d, found := decoders[remote.String()]
		if !found {
			s := session.New()
			d = netflow.NewDecoder(s)
			decoders[remote.String()] = d
		}

		m, err := d.Read(bytes.NewBuffer(buf[:octets]))
		if err != nil {
			log.Println("decoder error:", err)
			continue
		}

		dmp := ""

		switch p := m.(type) {
		case *netflow1.Packet:
			dmp = p1Dump(p)

		case *netflow5.Packet:
			dmp = p5Dump(p)

		case *netflow6.Packet:
			dmp = p6Dump(p)

		case *netflow7.Packet:
			dmp = p7Dump(p)

		case *netflow9.Packet:
			dmp = p9Dump(p)

		case *ipfix.Message:
			dmp = ipfixDump(p)
		}

		// log.Println(dmp)
		_, err = f.Write([]byte(dmp))
		if err != nil {
			log.Printf("failed to write to file '%v' due to error: '%v'...\n", *logPath, err)
		}

	}
}

func p1Dump(p *netflow1.Packet) string {
	s := fmt.Sprintln("NetFlow version 1 packet", p.Header)
	s += fmt.Sprintf("  %d flow records:\n", len(p.Records))
	for i, r := range p.Records {
		s += fmt.Sprintf("  record %d:\n", i)
		s += fmt.Sprintln("    srcAddr: ", r.SrcAddr)
		s += fmt.Sprintln("    srcPort: ", r.SrcPort)
		s += fmt.Sprintln("    dstAddr: ", r.DstAddr)
		s += fmt.Sprintln("    dstPort: ", r.DstPort)
		s += fmt.Sprintln("    nextHop: ", r.NextHop)
		s += fmt.Sprintln("    bytes:   ", r.Bytes)
		s += fmt.Sprintln("    packets: ", r.Packets)
		s += fmt.Sprintln("    first:   ", r.First)
		s += fmt.Sprintln("    last:    ", r.Last)
		s += fmt.Sprintln("    protocol:", r.Protocol, read.Protocol(r.Protocol))
		s += fmt.Sprintln("    tos:     ", r.ToS)
		s += fmt.Sprintln("    flags:   ", r.Flags, read.TCPFlags(r.Flags))
	}

	return s
}

func p5Dump(p *netflow5.Packet) string {
	s := fmt.Sprintln("NetFlow version 5 packet", p.Header)
	s += fmt.Sprintf("  %d flow records:\n", len(p.Records))
	for i, r := range p.Records {
		s += fmt.Sprintf("    record %d:\n", i)
		s += fmt.Sprintln("      srcAddr: ", r.SrcAddr)
		s += fmt.Sprintln("      srcPort: ", r.SrcPort)
		s += fmt.Sprintln("      dstAddr: ", r.DstAddr)
		s += fmt.Sprintln("      dstPort: ", r.DstPort)
		s += fmt.Sprintln("      nextHop: ", r.NextHop)
		s += fmt.Sprintln("      bytes:   ", r.Bytes)
		s += fmt.Sprintln("      packets: ", r.Packets)
		s += fmt.Sprintln("      first:   ", r.First)
		s += fmt.Sprintln("      last:    ", r.Last)
		s += fmt.Sprintln("      tcpflags:", r.TCPFlags, read.TCPFlags(r.TCPFlags))
		s += fmt.Sprintln("      protocol:", r.Protocol, read.Protocol(r.Protocol))
		s += fmt.Sprintln("      tos:     ", r.ToS)
		s += fmt.Sprintln("      srcAs:   ", r.SrcAS)
		s += fmt.Sprintln("      dstAs:   ", r.DstAS)
		s += fmt.Sprintln("      srcMask: ", r.SrcMask)
		s += fmt.Sprintln("      dstMask: ", r.DstMask)
	}

	return s
}

func p6Dump(p *netflow6.Packet) string {
	s := fmt.Sprintln("NetFlow version 6 packet", p.Header)
	s += fmt.Sprintf("  %d flow records:\n", len(p.Records))
	for i, r := range p.Records {
		s += fmt.Sprintf("    record %d:\n", i)
		s += fmt.Sprintln("      srcAddr: ", r.SrcAddr)
		s += fmt.Sprintln("      srcPort: ", r.SrcPort)
		s += fmt.Sprintln("      dstAddr: ", r.DstAddr)
		s += fmt.Sprintln("      dstPort: ", r.DstPort)
		s += fmt.Sprintln("      nextHop: ", r.NextHop)
		s += fmt.Sprintln("      bytes:   ", r.Bytes)
		s += fmt.Sprintln("      packets: ", r.Packets)
		s += fmt.Sprintln("      first:   ", r.First)
		s += fmt.Sprintln("      last:    ", r.Last)
		s += fmt.Sprintln("      tcpflags:", r.TCPFlags, read.TCPFlags(r.TCPFlags))
		s += fmt.Sprintln("      protocol:", r.Protocol, read.Protocol(r.Protocol))
		s += fmt.Sprintln("      tos:     ", r.ToS)
		s += fmt.Sprintln("      srcAs:   ", r.SrcAS)
		s += fmt.Sprintln("      dstAs:   ", r.DstAS)
		s += fmt.Sprintln("      srcMask: ", r.SrcMask)
		s += fmt.Sprintln("      dstMask: ", r.DstMask)
	}

	return s
}

func p7Dump(p *netflow7.Packet) string {
	s := fmt.Sprintln("NetFlow version 7 packet", p.Header)
	s += fmt.Sprintf("  %d flow records:\n", len(p.Records))
	for i, r := range p.Records {
		s += fmt.Sprintf("    record %d:\n", i)
		s += fmt.Sprintln("      srcAddr: ", r.SrcAddr)
		s += fmt.Sprintln("      srcPort: ", r.SrcPort)
		s += fmt.Sprintln("      dstAddr: ", r.DstAddr)
		s += fmt.Sprintln("      dstPort: ", r.DstPort)
		s += fmt.Sprintln("      nextHop: ", r.NextHop)
		s += fmt.Sprintln("      bytes:   ", r.Bytes)
		s += fmt.Sprintln("      packets: ", r.Packets)
		s += fmt.Sprintln("      first:   ", r.First)
		s += fmt.Sprintln("      last:    ", r.Last)
		s += fmt.Sprintln("      tcpflags:", r.TCPFlags, read.TCPFlags(r.TCPFlags))
		s += fmt.Sprintln("      protocol:", r.Protocol, read.Protocol(r.Protocol))
		s += fmt.Sprintln("      tos:     ", r.ToS)
		s += fmt.Sprintln("      srcAs:   ", r.SrcAS)
		s += fmt.Sprintln("      dstAs:   ", r.DstAS)
		s += fmt.Sprintln("      srcMask: ", r.SrcMask)
		s += fmt.Sprintln("      dstMask: ", r.DstMask)
		s += fmt.Sprintln("      flags:   ", r.Flags)
		s += fmt.Sprintln("      routerSC:", r.RouterSC)
	}

	return s
}

func p9Dump(p *netflow9.Packet) string {
	s := fmt.Sprintln("NetFlow version 9 packet")
	for _, ds := range p.DataFlowSets {
		s += fmt.Sprintf("  data set template %d, length: %d\n", ds.Header.ID, ds.Header.Length)
		if ds.Records == nil {
			s += fmt.Sprintf("    %d raw bytes:\n", len(ds.Bytes))
			s += fmt.Sprintln(hex.Dump(ds.Bytes))
			continue
		}
		s += fmt.Sprintf("    %d records:\n", len(ds.Records))
		for i, dr := range ds.Records {
			s += fmt.Sprintf("      record %d:\n", i)
			for _, f := range dr.Fields {
				if f.Translated != nil {
					if f.Translated.Name != "" {
						s += fmt.Sprintf("        %s: %v\n", f.Translated.Name, f.Translated.Value)
					} else {
						s += fmt.Sprintf("        %d: %v\n", f.Translated.Type, f.Bytes)
					}
				} else {
					s += fmt.Sprintf("        %d: %v (raw)\n", f.Type, f.Bytes)
				}
			}
		}
	}
	return s
}

func ipfixDump(m *ipfix.Message) string {
	s := fmt.Sprintln("IPFIX message")
	for _, ds := range m.DataSets {
		s += fmt.Sprintln("  data set")
		if ds.Records == nil {
			s += fmt.Sprintf("    %d raw bytes:\n", len(ds.Bytes))
			s += fmt.Sprintln(hex.Dump(ds.Bytes))
			continue
		}
		s += fmt.Sprintf("    %d records:\n", len(ds.Records))
		for i, dr := range ds.Records {
			s += fmt.Sprintf("      record %d:\n", i)
			for _, f := range dr.Fields {
				if f.Translated != nil {
					if f.Translated.Name != "" {
						s += fmt.Sprintf("        %s: %v\n", f.Translated.Name, f.Translated.Value)
					} else {
						s += fmt.Sprintf("        %d.%d: %v\n", f.Translated.EnterpriseNumber, f.Translated.InformationElementID, f.Bytes)
					}
				} else {
					s += fmt.Sprintf("        %v\n", f.Bytes)
				}
			}
		}
	}
	return s
}
