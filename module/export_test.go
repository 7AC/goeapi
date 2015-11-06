//
// Copyright (c) 2015, Arista Networks, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//   * Redistributions of source code must retain the above copyright notice,
//   this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright
//   notice, this list of conditions and the following disclaimer in the
//   documentation and/or other materials provided with the distribution.
//
//   * Neither the name of Arista Networks nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL ARISTA NETWORKS
// BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR
// BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN
// IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package module

import (
	"fmt"
	"os"
	"testing"

	"github.com/aristanetworks/goeapi"
)

// Acl
var AclParseEntries = (*AclEntity).parseEntries

// IpInterface
var IPIntfParseAddress = (*IPInterfaceEntity).parseAddress
var IPIntfParseMtu = (*IPInterfaceEntity).parseMtu

// Vlans
var VlanParseName = (*VlanEntity).parseName
var VlanParseState = (*VlanEntity).parseState
var VlanParseTrunkGroups = (*VlanEntity).parseTrunkGroups

// Mlag
var MlagParseDomainID = (*MlagEntity).parseDomainID
var MlagParseLocalInterface = (*MlagEntity).parseLocalInterface
var MlagParsePeerAddress = (*MlagEntity).parsePeerAddress
var MlagParsePeerLink = (*MlagEntity).parsePeerLink
var MlagParseShutdown = (*MlagEntity).parseShutdown

// SwitchPort
var SwPortParseMode = (*SwitchPortEntity).parseMode
var SwPortParseTrunkGroups = (*SwitchPortEntity).parseTrunkGroups
var SwPortParseAccessVlan = (*SwitchPortEntity).parseAccessVlan
var SwPortParseTrunkNativeVlan = (*SwitchPortEntity).parseTrunkNativeVlan
var SwPortParseTrunkAllowedVlans = (*SwitchPortEntity).parseTrunkAllowedVlans

// BGP
var BgpParseAS = (*BGPEntity).parseAS
var BgpParseRouterID = (*BGPEntity).parseRouterID
var BgpParseMaxPaths = (*BGPEntity).parseMaxPaths
var BgpParseShutdown = (*BGPEntity).parseShutdown
var BgpParseNetworks = (*BGPEntity).parseNetworks
var BgpNeighborParsePeerGroup = (*BgpNeighborsEntity).parsePeerGroup
var BgpNeighborParseRemoteAS = (*BgpNeighborsEntity).parseRemoteAS
var BgpNeighborParseSendCommunity = (*BgpNeighborsEntity).parseSendCommunity
var BgpNeighborParseShutdown = (*BgpNeighborsEntity).parseShutdown
var BgpNeighborParseDescription = (*BgpNeighborsEntity).parseDescription
var BgpNeighborParseNextHopSelf = (*BgpNeighborsEntity).parseNextHopSelf
var BgpNeighborParseRouteMapIn = (*BgpNeighborsEntity).parseRouteMapIn
var BgpNeighborParseRouteMapOut = (*BgpNeighborsEntity).parseRouteMapOut

// STP
var STPIntfParseBPDUGuard = (*STPInterfaceEntity).parseBPDUGuard
var STPIntfParsePortFast = (*STPInterfaceEntity).parsePortfast
var STPIntfParsePortFastType = (*STPInterfaceEntity).parsePortfastType

// Interface

var BaseInterfaceParseShutdown = (*BaseInterfaceEntity).parseShutdown
var BaseInterfaceParseDescription = (*BaseInterfaceEntity).parseDescription

var EthernetParseSflow = (*EthernetInterfaceEntity).parseSflow
var EthernetParseFlowControlSend = (*EthernetInterfaceEntity).parseFlowControlSend
var EthernetParseFlowControlReceive = (*EthernetInterfaceEntity).parseFlowControlReceive

var PortChannelParseMinimumLinks = (*PortChannelInterfaceEntity).parseMinimumLinks
var PortChannelGetMembers = (*PortChannelInterfaceEntity).getMembers
var PortChannelGetLacpMode = (*PortChannelInterfaceEntity).getLacpMode

var VxlanParseSourceInterface = (*VxlanInterfaceEntity).parseSourceInterface
var VxlanParseMulticastGroup = (*VxlanInterfaceEntity).parseMulticastGroup
var VxlanParseUDPPort = (*VxlanInterfaceEntity).parseUDPPort
var VxlanParseVlans = (*VxlanInterfaceEntity).parseVlans
var VxlanParseFloodList = (*VxlanInterfaceEntity).parseFloodList

const checkMark = "\u2713"
const xMark = "\u2717"

/*
 ****************************************************************************
 *
 * DummyEapiConnection is a Dummy connection object that adheres to the
 * EapiConnection Inteface. The Execute() method below (currently) returns
 * a non-error with allocated JSONRPCResponse so the upper layer API can
 * be tested. Commands received by this DummyConnection are cached and retreived
 * to compare to what would be sent.
 *
 * Note:
 *		Execute() clears the the previous cached list of commands and replaces
 *		with current command list.
 *
 ****************************************************************************
 */
type DummyEapiConnection struct {
	goeapi.EapiConnection
	Commands []interface{}
}

func NewDummyEapiConnection(transport string, host string, username string,
	password string, port int) *DummyEapiConnection {
	conn := goeapi.EapiConnection{}
	return &DummyEapiConnection{EapiConnection: conn}
}

func (conn *DummyEapiConnection) Execute(commands []interface{},
	encoding string) (*goeapi.JSONRPCResponse, error) {
	if conn == nil {
		return &goeapi.JSONRPCResponse{}, fmt.Errorf("No connection")
	}
	conn.Commands = nil
	conn.Commands = append(conn.Commands, commands...)
	if encoding == "json" {
		return &goeapi.JSONRPCResponse{Result: make([]map[string]interface{}, len(commands))}, nil
	} else {
		resp := &goeapi.JSONRPCResponse{Result: make([]map[string]interface{}, len(commands))}
		for idx := range resp.Result {
			resp.Result[idx] = make(map[string]interface{})
			resp.Result[idx]["output"] = ""
		}
		return resp, nil
	}
}

// Retreive the cached list of commands from the connection.
func (conn *DummyEapiConnection) GetCommands() []interface{} {
	return conn.Commands
}

func initFixture() {
	if initFile {
		return
	}
	initFile = true
	dummyNode.SetRunningConfig(LoadFixtureFile("running_config.text"))
	return
}

var runConf string
var duts []*goeapi.Node
var dummyNode *goeapi.Node
var dummyConnection *DummyEapiConnection

func TestMain(m *testing.M) {
	fmt.Println("Export_test.go")
	runConf = GetFixture("running_config.text")
	goeapi.LoadConfig(GetFixture("dut.conf"))
	conns := goeapi.Connections()
	fmt.Println("Connections: ", conns)
	for _, name := range conns {
		if name != "localhost" {
			node, _ := goeapi.ConnectTo(name)
			duts = append(duts, node)
		}
	}

	// Create a Node with a DummyConnection to be used in
	// UnitTests.
	dummyConnection = NewDummyEapiConnection("", "", "", "", 0)
	dummyNode = &goeapi.Node{}
	dummyNode.SetAutoRefresh(false)
	dummyNode.SetConnection(dummyConnection)

	//
	os.Exit(m.Run())
}