package ip

import (
	"fmt"
	"io/ioutil"
	"net"
	"sort"

	"github.com/pelletier/go-toml"
)

// AllowList ip列表
type AllowList struct {
	IPs    []net.IP
	IPNets []*net.IPNet
}

func (l *AllowList) sort() {
	var compare = func(ip1, ip2 net.IP) bool {
		ni, nj := len(ip1), len(ip2)
		switch {
		case ni < nj:
			return true
		case ni > nj:
			return false
		default:
			for k := 0; k < ni; k++ {
				if ip1[k] < ip2[k] {
					return true
				} else if ip1[k] > ip2[k] {
					return false
				}
			}
			return false
		}
	}
	sort.Slice(l.IPs, func(i, j int) bool {
		return compare(l.IPs[i], l.IPs[j])
	})
	sort.Slice(l.IPNets, func(i, j int) bool {
		return compare(l.IPNets[i].IP, l.IPNets[j].IP)

	})
}

// ContainsIP 判断是否包含指定ip
func (l *AllowList) ContainsIP(ip1 net.IP) bool {
	for _, ip2 := range l.IPs {
		if ip2.Equal(ip1) {
			return true
		}
	}
	for _, v := range l.IPNets {
		if v.Contains(ip1) {
			return true
		}
	}
	return false
}

// ContainsString 是否包含目标ip
func (l *AllowList) ContainsString(s string) bool {
	fmt.Println(s)
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}
	return l.ContainsIP(ip)
}

type option struct {
	IPs []string `toml:"ips"`
}

func arrayToOption(a []string) option {
	return option{IPs: a}
}

func (o option) NewList() *AllowList {
	l := new(AllowList)
	ips, ipNets := make(map[string]net.IP), make(map[string]*net.IPNet)
	for _, v := range o.IPs {
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			allowIP := net.ParseIP(v)
			if allowIP == nil {
				continue
			}
			ips[allowIP.String()] = allowIP
		} else if ipNet != nil {
			ipNets[ipNet.String()] = ipNet
		}
	}
L1:
	for key, v1 := range ipNets {
		for j, v2 := range l.IPNets {
			if v1.Contains(v2.IP) || v2.Contains(v1.IP) {
				v1Ones, _ := v1.Mask.Size()
				v2Ones, _ := v2.Mask.Size()
				if v1Ones < v2Ones {
					l.IPNets[j] = v1
				}
				continue L1
			}
		}
		l.IPNets = append(l.IPNets, ipNets[key])
	}
	for key, v1 := range ips {
		if l.ContainsIP(v1) {
			continue
		}
		l.IPs = append(l.IPs, ips[key])
	}
	l.sort()
	return l
}

func loadFile(filename string) (list *AllowList, err error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		err = fmt.Errorf("ioutil.ReadFile(%s) %w", filename, err)
		return
	}
	return loadBytes(b)
}

func loadBytes(b []byte) (list *AllowList, err error) {
	if len(b) == 0 {
		return
	}
	var o option
	if err = toml.Unmarshal(b, &o); err != nil {
		err = fmt.Errorf("toml.Unmarshal() %w", err)
		return
	}
	list = o.NewList()
	return
}

// LoadArray 加载白名单, ip列表参数为字符串数组
func LoadArray(a []string) *AllowList {
	return arrayToOption(a).NewList()
}
