package util

import (
	"sort"
	"sync"

	"github.com/bwmarrin/snowflake"
)

//Pair -
type Pair struct {
	Key   string
	Value *EventMessage
}

//PairList -
type PairList []Pair

//Len -
func (p PairList) Len() int { return len(p) }

//Less -
func (p PairList) Less(i, j int) bool { return p[i].Value.Status > p[j].Value.Status }

//Swap -
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

//MakeList -
func MakeList(hash map[string]*EventMessage) PairList {
	pl := make(PairList, len(hash))
	i := 0
	for k, v := range hash {
		pl[i] = Pair{k, v}
		i++
	}
	if i > 1 {
		sort.Sort(sort.Reverse(pl))
	}

	return pl
}

//SortList -
func SortList(pl PairList) PairList {
	if len(pl) > 1 {
		sort.Sort(sort.Reverse(pl))
	}

	return pl
}

//AgentMap -
type AgentMap struct {
	sync.Mutex
	_hash map[string]*EventMessage
	node  *snowflake.Node
}

//NewAgentMap -
func NewAgentMap() *AgentMap {
	n, err := snowflake.NewNode(1)
	if err != nil {
		return nil
	}
	return &AgentMap{
		_hash: make(map[string]*EventMessage),
		node:  n,
	}
}

//GenerateID -
func (m *AgentMap) GenerateID() int64 {
	if m.node == nil {
		return 0
	}
	var newid int64
	m.Lock()
	newid = m.node.Generate().Int64()
	m.Unlock()

	return newid
}

//Exists -
func (m *AgentMap) Exists(key string) bool {
	var found bool = false
	m.Lock()
	if _, ok := m._hash[key]; ok {
		found = true
	}
	m.Unlock()
	return found
}

//Set -
func (m *AgentMap) Set(key string, e *EventMessage) {
	m.Lock()
	if _, ok := m._hash[key]; !ok {
		m._hash[key] = &EventMessage{
			Mode:      e.Mode,
			IP:        e.IP,
			Port:      e.Port,
			TimeStamp: e.TimeStamp,
			Status:    0,
		}
	}
	m.Unlock()
}

//Get -
func (m *AgentMap) Get(key string) *EventMessage {
	var val *EventMessage
	m.Lock()
	if v, ok := m._hash[key]; ok {
		val = v
	}
	m.Unlock()
	return val
}

//UpdateStatus -
func (m *AgentMap) UpdateStatus(key string, status int) {
	m.Lock()
	if s, ok := m._hash[key]; ok {
		s.Status = status
	}
	m.Unlock()
}

//Delete -
func (m *AgentMap) Delete(key string) {
	m.Lock()
	delete(m._hash, key)
	m.Unlock()
}

//DeleteList -
func (m *AgentMap) DeleteList(keys []string) {
	m.Lock()
	for _, key := range keys {
		delete(m._hash, key)
	}
	m.Unlock()
}

//Keys -
func (m *AgentMap) Keys() []string {
	keys := make([]string, 0)
	m.Lock()
	for k := range m._hash {
		keys = append(keys, k)
	}
	m.Unlock()
	return keys
}

//List -
func (m *AgentMap) List() []*EventMessage {
	result := make([]*EventMessage, 0)
	m.Lock()
	for _, v := range m._hash {
		result = append(result, &EventMessage{
			Mode:      v.Mode,
			IP:        v.IP,
			Port:      v.Port,
			TimeStamp: v.TimeStamp,
			Status:    v.Status,
		})
	}
	m.Unlock()

	return result
}

//MinKey - search min status
func (m *AgentMap) MinKey() string {
	var key string
	m.Lock()
	list := MakeList(m._hash)
	if len(list) > 0 {
		key = list[0].Key
	}
	m.Unlock()

	return key
}
