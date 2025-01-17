// pkg/check/korcen.go

package check

import (
	"container/list"
	"encoding/xml"
	"errors"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/fluffy-melli/korcen-go"
)

type Header struct {
	XMLName xml.Name `json:"-" xml:"header"`
	Input   string   `json:"input" xml:"input"`
	Start   string   `json:"replace-front" xml:"replace-front"`
	End     string   `json:"replace-end" xml:"replace-end"`
}

type Respond struct {
	XMLName   xml.Name `json:"-" xml:"respond"`
	Detect    bool     `json:"detect" xml:"detect"`
	Swear     string   `json:"swear" xml:"swear"`
	String    string   `json:"input" xml:"input"`
	NewString string   `json:"output" xml:"output"`
}

type KorcenResult = korcen.CheckInfo

var korcenPool = sync.Pool{
	New: func() interface{} {
		return &KorcenResult{}
	},
}

func newKorcenResult() *KorcenResult {
	return korcenPool.Get().(*KorcenResult)
}

func freeKorcenResult(result *KorcenResult) {
	*result = KorcenResult{}
	korcenPool.Put(result)
}

type Node struct {
	key   string
	value *KorcenResult
	elem  *list.Element
}

var nodePool = sync.Pool{
	New: func() interface{} {
		return &Node{}
	},
}

func newNode(key string, value *KorcenResult) *Node {
	n := nodePool.Get().(*Node)
	n.key = key
	n.value = value
	n.elem = &list.Element{Value: n}
	return n
}

func freeNode(n *Node) {
	n.key = ""
	n.value = nil
	n.elem = nil
	nodePool.Put(n)
}

type LRUCache struct {
	mutex    sync.Mutex
	capacity int
	ll       *list.List
	items    map[string]*list.Element
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		ll:       list.New(),
		items:    make(map[string]*list.Element),
	}
}

func (c *LRUCache) Get(key string) (*KorcenResult, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		c.ll.MoveToFront(elem)
		node := elem.Value.(*Node)
		return node.value, true
	}
	return nil, false
}

func (c *LRUCache) Set(key string, value *KorcenResult) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		c.ll.MoveToFront(elem)
		node := elem.Value.(*Node)
		node.value = value
		return
	}

	n := newNode(key, value)
	elem := c.ll.PushFront(n)
	n.elem = elem
	c.items[key] = elem

	if c.ll.Len() > c.capacity {
		old := c.ll.Back()
		if old != nil {
			c.removeElement(old)
		}
	}
}

func (c *LRUCache) removeElement(elem *list.Element) {
	c.ll.Remove(elem)
	node := elem.Value.(*Node)
	delete(c.items, node.key)

	freeKorcenResult(node.value)
	freeNode(node)
}

var globalLRU = NewLRUCache(1000)

func Korcen(header *Header) *Respond {
	if info, ok := globalLRU.Get(header.Input); ok {
		return buildRespond(header, info)
	}

	info := newKorcenResult()
	*info = *korcen.Check(header.Input)

	globalLRU.Set(header.Input, info)

	return buildRespond(header, info)
}

func buildRespond(header *Header, info *KorcenResult) *Respond {
	if !info.Detect {
		return &Respond{
			Detect:    false,
			Swear:     "",
			String:    header.Input,
			NewString: header.Input,
		}
	}

	firstSwear := info.Swear[0]
	formattedMessage := formatMessage(info.NewText, firstSwear.Start, firstSwear.End, header.Start, header.End)

	return &Respond{
		Detect:    true,
		Swear:     firstSwear.Swear,
		String:    header.Input,
		NewString: formattedMessage,
	}
}

func formatMessage(text string, start, end int, prefix, suffix string) string {
	if start < 0 || end > len(text) || start > end {
		return text
	}

	switch {
	case prefix != "" && suffix != "":
		return text[:start] + prefix + text[start:end] + suffix + text[end:]
	case prefix != "":
		return text[:start] + prefix + text[start:end] + text[end:]
	case suffix != "":
		return text[:start] + text[start:end] + suffix + text[end:]
	default:
		return text
	}
}

// ---------------------------------------------------------------------
// Actor
// ---------------------------------------------------------------------

type KorcenRequest struct {
	Header *Header
}

type KorcenResponse struct {
	Respond *Respond
	Err     error
}

type KorcenActor struct{}

func (k *KorcenActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {

	case *KorcenRequest:
		if msg.Header == nil {
			context.Respond(&KorcenResponse{
				Respond: nil,
				Err:     errors.New("KorcenActor: Header is nil"),
			})
			return
		}
		resp := Korcen(msg.Header)

		context.Respond(&KorcenResponse{
			Respond: resp,
			Err:     nil,
		})

	default:
	}
}
