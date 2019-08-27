package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	//nodeBits uint8 = 10
	nodeBits uint8 = 6
	//stepBits uint8 = 12
	stepBits  uint8 = 6
	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	stepMax   int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
)

// epoch from 2019-08-01 00:00:00
var Epoch int64 = 1564588800000


type Node struct {
	mu        sync.Mutex
	timestamp int64
	node      int64
	step      int64
}

func (n *Node) Generate() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	if n.timestamp == now {
		n.step++
		if n.step > stepMax {
			for now <= n.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}

	} else {
		n.step = 0
	}

	n.timestamp = now

	result := int64((now-Epoch)<<timeShift | (n.node << nodeShift) | (n.step))

	return result
}

func NewNode(node int64) (*Node, error) {

	if node < 0 || node > nodeMax {
		return nil, errors.New("Node number must be between 0 and 64")
	}

	return &Node{
		timestamp: 0,
		node:      node,
		step:      0,
	}, nil
}
