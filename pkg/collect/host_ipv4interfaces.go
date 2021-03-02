package collect

import (
	"encoding/json"
	"net"

	"github.com/pkg/errors"
	troubleshootv1beta2 "github.com/replicatedhq/troubleshoot/pkg/apis/troubleshoot/v1beta2"
)

type CollectHostIPV4Interfaces struct {
	hostCollector *troubleshootv1beta2.IPV4Interfaces
}

func (c *CollectHostIPV4Interfaces) Title() string {
	return hostCollectorTitleOrDefault(c.hostCollector.HostCollectorMeta, "IPv4 Interfaces")
}

func (c *CollectHostIPV4Interfaces) IsExcluded() (bool, error) {
	return isExcluded(c.hostCollector.Exclude)
}

func (c *CollectHostIPV4Interfaces) Collect(progressChan chan<- interface{}) (map[string][]byte, error) {
	var ipv4Interfaces []net.Interface

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.Wrap(err, "list host network interfaces")
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		ip, _ := getIPv4FromInterface(&iface)
		if ip == nil {
			continue
		}
		ipv4Interfaces = append(ipv4Interfaces, iface)
	}

	b, err := json.Marshal(ipv4Interfaces)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal network interfaces")
	}

	return map[string][]byte{
		"system/ipv4Interfaces.json": b,
	}, nil
}
