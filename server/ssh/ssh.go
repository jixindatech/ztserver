package ssh

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"time"
	"zt-server/pkg/core/golog"
	"zt-server/settings"
)

var closeDoor int = 0
var openDoor int = 1
var clearFw int = 2
var getFwEntry int = 3

type sshConfig struct {
	Servers  map[string]int
	Fw       string
	Ipset    string
	Host     string
	User     string
	Password string
}

var config []*sshConfig

func Init(ssh []settings.Ssh) error {
	for _, v := range ssh {
		item := &sshConfig{}
		item.Host = v.Host
		item.Fw = v.Fw
		item.Ipset = v.Ipset
		item.User = v.User
		item.Password = v.Password
		err := sshClearFw(&v)
		if err != nil {
			return err
		}

		item.Servers = make(map[string]int)
		for _, v1 := range v.Server {
			item.Servers[v1] = 1
		}
		config = append(config, item)
	}
	return nil
}

func sshClearFw(config *settings.Ssh) error {
	clientConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", config.Host, clientConfig)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}

	ipset := config.Ipset
	cmdRun := fmt.Sprintf(firewallCmd[getFwEntry], ipset)

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmdRun); err != nil {
		return err
	}
	session.Close()

	ips := strings.Split(string(b.Bytes()), "\n")

	cmd := ""
	for _, ip := range ips {
		if len(ip) ==0 {
			continue
		}

		cmdRun = fmt.Sprintf(firewallCmd[closeDoor], ipset, ip)
		cmd = cmd + cmdRun + ";"
	}

	if len(cmd) > 0 {
		session, err = client.NewSession()
		if err != nil {
			return err
		}
		_, err := session.CombinedOutput(cmd)
		if err != nil {
			return err
		}

		session.Close()
	}

	return nil
}

func runCmd(host *sshConfig, cmd int, ip string) error {
	sshConfig := &ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", host.Host, sshConfig)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var cmdStr string
	switch host.Fw {
	case firewall:
		cmdStr = firewallCmd[cmd]
	default:
		return fmt.Errorf("%s", "unknown fw type")
	}

	cmdRun := fmt.Sprintf(cmdStr, host.Ipset, ip)
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmdRun); err != nil {
		return err
	}

	return nil
}

func execCmd(hosts []*sshConfig, cmd int, ip string) []string {
	var ips []string
	for _, host := range hosts {
		err := runCmd(host, cmd, ip)
		if err != nil {
			golog.Error("ssh", zap.String("runcmd", err.Error()))
		}
		ipArray := strings.Split(host.Host, ":")
		if len(ipArray) == 2 {
			ips = append(ips, ipArray[0])
		} else {
			golog.Error("ssh", zap.String("runcmd",
				fmt.Errorf("%s", "invalid ssh Host format, ':' is missing").Error()))
		}
	}
	return ips
}

func getHostsSshConfig(servers map[string]int) []*sshConfig {
	var hosts []*sshConfig
	for k := range servers {
		for _, v := range config {
			if _, ok := v.Servers[k]; ok {
				/* skip the same item */
				found := false
				for _, host := range hosts {
					if host == v {
						found = true
					}
				}
				if !found {
					hosts = append(hosts, v)
				}
			}
		}
	}

	return hosts
}

func OpenConnection(ip string, servers map[string]int) []string {
	hostsSshConfig := getHostsSshConfig(servers)
	return execCmd(hostsSshConfig, openDoor, ip)
}

func CloseConnection(ip string, servers map[string]int) []string {
	hostsSshConfig := getHostsSshConfig(servers)
	return execCmd(hostsSshConfig, closeDoor, ip)
}
