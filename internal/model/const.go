package model

type (
	PORT     int
	PROTOCOL string
	STATUS   string
)

const (
	PORT80  PORT = 80
	PORT443 PORT = 443

	TCP  PROTOCOL = "TCP"
	UDP  PROTOCOL = "UDP"
	ICMP PROTOCOL = "ICMP"
	ALL  PROTOCOL = "ALL"

	ACTIVE  STATUS = "Active"
	DISABLE STATUS = "Disable"
	UP      STATUS = "Up"
	DOWN    STATUS = "Down"
)
