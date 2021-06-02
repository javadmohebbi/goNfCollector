package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/goNfCollector/common"
	"github.com/goNfCollector/configurations"
	"github.com/goNfCollector/database/model"
	"github.com/sirupsen/logrus"
)

// write metrics to db
func (p *Postgres) write(metrics []common.Metric) {

	if p.closed {
		return
	}

	successWrites := 0

	p.WaitGroup.Add(1)
	defer p.WaitGroup.Done()

	// define device ID default value
	var deviceID uint = 0
	var err error

	// check if metrics length > 0
	if len(metrics) > 0 {
		// get first array device
		// because all of them are the same in the
		// further loop
		deviceID, err = p.writeDevice(metrics[0].Device)

		// if err not null
		// return with log
		if err != nil {
			p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (%v)",
				configurations.ERROR_CAN_T_INSERT_METRICS_TO_POSTGRES_DB.Int(),
				configurations.ERROR_CAN_T_INSERT_METRICS_TO_POSTGRES_DB, err),
				logrus.ErrorLevel,
			)
			return
		}
	} else {
		// no metrics
		p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: (metrics length: %v)",
			configurations.ERROR_NO_METRICS_IN_THE_ARRAY.Int(),
			configurations.ERROR_NO_METRICS_IN_THE_ARRAY, len(metrics)),
			logrus.ErrorLevel,
		)
		return
	}

	// arrFlows defines for bach insert
	// on database to decrease the chance of
	// max_connection error
	var arrFlows []model.Flow

	// loop through metrics and write
	for _, m := range metrics {

		// extract version
		verID, err := p._getVersion(m.FlowVersion)
		if err != nil {
			continue
		}

		// extract protocol
		protoID, err := p.writeProtocol(m.Protocol, m.ProtoName)
		if err != nil {
			continue
		}

		// extract Source ASN
		_, srcAsnID, err := p.writeAutonomous(m.SrcIP)
		if err != nil {
			continue
		}

		// extract Destination ASN
		_, dstAsnID, err := p.writeAutonomous(m.DstIP)
		if err != nil {
			continue
		}

		// extract src host
		srcHostID, err := p.writeHost(m.SrcIP)
		if err != nil {
			continue
		}
		// extract dst host
		dstHostID, err := p.writeHost(m.DstIP)
		if err != nil {
			continue
		}

		// extract src port
		srcPortID, err := p.writePort(m.SrcPortName, m.ProtoName, m.SrcPort)
		if err != nil {
			continue
		}
		// extract dst port
		dstPortID, err := p.writePort(m.DstPortName, m.ProtoName, m.DstPort)
		if err != nil {
			continue
		}

		// extract src geo
		srcGeoID, err := p.writeGeo(m.SrcIP)
		if err != nil {
			continue
		}
		// extract dst geo
		dstGeoID, err := p.writeGeo(m.DstIP)
		if err != nil {
			continue
		}

		// extract next hop host
		nextHopHostID, err := p.writeHost(m.NextHop)
		if err != nil {
			continue
		}

		// extract next hop geo
		nextHopGeoID, err := p.writeGeo(m.NextHop)
		if err != nil {
			continue
		}

		// extract flags
		flagsID, fin, syn, rst, psh, ack, urg, ece, cwr, err := p.writeFlag(m.TCPFlags)
		if err != nil {
			continue
		}

		// // log.Println(m.Device, m.FlowVersion, "===", deviceID, verID)
		// fmt.Println()
		// fmt.Printf("device: %v ID: %v\n   version: %v ID: %v\n   proto: %v (%v) ID: %v\n    SRC Host: %v ID: %v\n    SRC ASN: %v ID: %v\n    DST Host: %v ID: %v\n    DST ASN: %v ID: %v\n    SRC PORT: %v ID: %v\n    DST PORT: %v ID: %v\n    SRC GEO ID: %v\n    DST GEO ID: %v\n    TCP Flags: %v\n    Next Hop: %v NxtHopID: %v    GEO ID: %v",
		// 	m.Device, deviceID,
		// 	m.FlowVersion, verID,
		// 	m.Protocol, m.ProtoName, protoID,
		// 	m.SrcIP, srcHostID,
		// 	srcASN, srcAsnID,
		// 	m.DstIP, dstHostID,
		// 	dstASN, dstAsnID,
		// 	m.SrcPortName, srcPortID,
		// 	m.DstPortName, dstPortID,
		// 	srcGeoID, dstGeoID,
		// 	flagsID,
		// 	m.NextHop, nextHopHostID, nextHopGeoID,
		// )
		// fmt.Println()

		t := time.Now().Add(-time.Duration(m.Time.Second()))

		by, _ := strconv.Atoi(m.Bytes)
		pa, _ := strconv.Atoi(m.Packets)

		// check for ip reputation
		srcThreatID, srcIsThreat, _ := p.writeThreat(m.SrcIP, srcHostID)
		dstThreatID, dstIsThreat, _ := p.writeThreat(m.DstIP, dstHostID)

		// if next hop is valid, check that
		var nxtHopThreatID uint
		var nxtHopIsThreat bool
		if m.NextHop != "0.0.0.0" {
			nxtHopThreatID, nxtHopIsThreat, _ = p.writeThreat(m.NextHop, nextHopHostID)
		}

		flow := model.Flow{
			DeviceID:  deviceID,
			VersionID: verID,

			ProtocolID: protoID,

			SrcASNID:  srcAsnID,
			SrcHostID: srcHostID,
			SrcPortID: srcPortID,
			SrcGeoID:  srcGeoID,

			DstASNID:  dstAsnID,
			DstHostID: dstHostID,
			DstPortID: dstPortID,
			DstGeoID:  dstGeoID,

			DstIsThreat: dstIsThreat,

			NextHopID:    nextHopHostID,
			NextHopGeoID: nextHopGeoID,

			NextHopIsThreat: nxtHopIsThreat,

			FlagID: flagsID,

			FlagFin: fin,
			FlagSyn: syn,
			FlagRst: rst,
			FlagPsh: psh,
			FlagAck: ack,
			FlagUrg: urg,
			FlagEce: ece,
			FlagCwr: cwr,

			Byte:   uint(by),
			Packet: uint(pa),
		}

		// src threats
		if srcIsThreat && srcThreatID != 0 {
			flow.SrcIsThreat = srcIsThreat
			flow.SrcThreatID = &srcThreatID
		}
		// dst threats
		if dstIsThreat && dstThreatID != 0 {
			flow.DstIsThreat = dstIsThreat
			flow.DstThreatID = &dstThreatID
		}
		// next hop threats
		if nxtHopIsThreat && nxtHopThreatID != 0 {
			flow.NextHopIsThreat = nxtHopIsThreat
			flow.NextHopThreatID = &nxtHopThreatID
		}

		// set flow real date/time
		flow.CreatedAt = t
		flow.UpdatedAt = t

		// append flow to arrays
		arrFlows = append(arrFlows, flow)

	}

	result := p.db.CreateInBatches(arrFlows, len(arrFlows))

	if result.Error != nil {
		p.Debuuger.Verbose(fmt.Sprintf("[%d]-%s: ([FLOW] %v)",
			configurations.ERROR_CAN_T_INSERT_METRICS_TO_POSTGRES_DB.Int(),
			configurations.ERROR_CAN_T_INSERT_METRICS_TO_POSTGRES_DB, result.Error),
			logrus.ErrorLevel,
		)
		if len(metrics)-int(result.RowsAffected) > 0 {
			p.Debuuger.Verbose(fmt.Sprintf("'%v' has not been insterted due to error.", len(metrics)-successWrites), logrus.DebugLevel)
		}
	}

	if result.RowsAffected > 0 {
		p.Debuuger.Verbose(fmt.Sprintf("'%v' out of '%v' has been inserted to db.", int(result.RowsAffected), len(metrics)), logrus.DebugLevel)
	}

}
