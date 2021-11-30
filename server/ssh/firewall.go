package ssh

var firewall string = "firewall"

var firewallCmd = map[int]string{
	closeDoor: "firewall-cmd --ipset=%s --remove-entry=%s",
	openDoor:  "firewall-cmd --ipset=%s --add-entry=%s",
	clearFw:   "firewall-cmd --reload",
	getFwEntry: "firewall-cmd --ipset=%s --get-entries",
}
