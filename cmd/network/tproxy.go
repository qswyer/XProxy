package network

import (
    log "github.com/sirupsen/logrus"
    "strconv"
)

func loadV4TProxy(v4 Config, v4SysCidr []string) {
    log.Info("Setting up TProxy of IPv4")
    tableNum := strconv.Itoa(v4.RouteTable)
    v4Bypass := append(v4SysCidr, v4.Bypass...)
    run("ip", "-4", "rule", "add", "fwmark", "1", "table", tableNum)
    run("ip", "-4", "route", "add", "local", "0.0.0.0/0", "dev", "lo", "table", tableNum)
    run("iptables", "-t", "mangle", "-N", "XPROXY")
    log.Infof("Setting up IPv4 bypass CIDR -> %v", v4Bypass)
    for _, cidr := range v4Bypass {
        run("iptables", "-t", "mangle", "-A", "XPROXY", "-d", cidr, "-j", "RETURN")
    }
    run("iptables", "-t", "mangle", "-A", "XPROXY",
        "-p", "tcp", "-j", "TPROXY", "--on-port", strconv.Itoa(v4.TProxyPort), "--tproxy-mark", "1")
    run("iptables", "-t", "mangle", "-A", "XPROXY",
        "-p", "udp", "-j", "TPROXY", "--on-port", strconv.Itoa(v4.TProxyPort), "--tproxy-mark", "1")
    run("iptables", "-t", "mangle", "-A", "PREROUTING", "-j", "XPROXY")
}

func loadV6TProxy(v6 Config, v6SysCidr []string) {
    log.Info("Setting up TProxy of IPv6")
    tableNum := strconv.Itoa(v6.RouteTable)
    v6Bypass := append(v6SysCidr, v6.Bypass...)
    run("ip", "-6", "rule", "add", "fwmark", "1", "table", tableNum)
    run("ip", "-6", "route", "add", "local", "::/0", "dev", "lo", "table", tableNum)
    run("ip6tables", "-t", "mangle", "-N", "XPROXY6")
    log.Infof("Setting up IPv6 bypass CIDR -> %v", v6Bypass)
    for _, cidr := range v6Bypass {
        run("ip6tables", "-t", "mangle", "-A", "XPROXY6", "-d", cidr, "-j", "RETURN")
    }
    run("ip6tables", "-t", "mangle", "-A", "XPROXY6",
        "-p", "tcp", "-j", "TPROXY", "--on-port", strconv.Itoa(v6.TProxyPort), "--tproxy-mark", "1")
    run("ip6tables", "-t", "mangle", "-A", "XPROXY6",
        "-p", "udp", "-j", "TPROXY", "--on-port", strconv.Itoa(v6.TProxyPort), "--tproxy-mark", "1")
    run("ip6tables", "-t", "mangle", "-A", "PREROUTING", "-j", "XPROXY6")
}
