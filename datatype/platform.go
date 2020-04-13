package datatype

import (
	"net"
	"time"
)

const (
	MIN_MASK_LEN      = 0
	STANDARD_MASK_LEN = 16
	MAX_MASK_LEN      = 32
	MAX_MASK6_LEN     = 128
	MASK_LEN_NUM      = MAX_MASK_LEN + 1

	IF_TYPE_WAN = 3

	DATA_VALID_TIME = 1 * time.Minute
	ARP_VALID_TIME  = 1 * time.Minute
)

type IpNet struct {
	RawIp    net.IP
	Netmask  uint32
	SubnetId uint32
}

type PlatformData struct {
	Mac    uint64
	Ips    []*IpNet
	EpcId  int32
	IfType uint8
}
