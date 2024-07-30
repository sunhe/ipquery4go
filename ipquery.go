package ipquery4go

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	IPV4 = net.IPv4len
	IPV6 = net.IPv6len
)

type IPQuery struct {
	len  int
	tree *radixTreeNode
}

func New(len int) *IPQuery {
	ipq := new(IPQuery)
	ipq.len = len
	return ipq
}

func CreateIPV4FromYamlFile(file string, sep string, part int) *IPQuery {
	ipq := New(IPV4)
	out, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	err = ipq.BuildFromYaml(out, sep, part)
	if err != nil {
		return nil
	}
	return ipq
}

func (ipq *IPQuery) BuildFromYaml(in []byte, sep string, part int) error {
	tree := new(radixTreeNode)
	db := make(map[string][]string)
	err := yaml.Unmarshal(in, &db)
	if err != nil {
		return err
	}
	for k, v := range db {
		ar := strings.Split(k, sep)
		if len(ar) != part {
			return fmt.Errorf("%s part is not %d", k, part)
		}
		for i := range v {
			if !strings.Contains(v[i], "/") {
				switch ipq.len {
				case IPV4:
					v[i] += "/32"
				case IPV6:
					//to do
				}
			}
			_, ipnet, err := net.ParseCIDR(v[i])
			if err != nil {
				return err
			}
			tree.insert(ipnet.IP.To4(), ipnet.Mask, ipq.len, ar)
		}
	}
	ipq.tree = tree
	return nil
}

func (ipq *IPQuery) Query(ipstr string) []string {
	ip := net.ParseIP(ipstr).To4()
	value := ipq.tree.query(ip, ipq.len)
	if value == nil {
		return nil
	}
	return value.([]string)
}

func (ipq *IPQuery) Destroy() {
	ipq.tree.destroy()
}
