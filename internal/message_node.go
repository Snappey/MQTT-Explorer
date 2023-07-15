package internal

import (
    "fmt"
    "strconv"
    "strings"
    "time"
)

type MessageNode struct {
    Children     OrderedMap[*MessageNode]
    Payloads     [][]byte
    ReceivedAt   time.Time
    MessageCount int
    Topic        string
    Segment      string
    Depth        uint
    Parent       *MessageNode
}

func (n MessageNode) GetAllDescendantsBFS() []MessageNode {
    return n.GetAllDescendants(BFS)
}

func (n MessageNode) GetAllDescendantsDFS() []MessageNode {
    return n.GetAllDescendants(DFS)
}

func (n MessageNode) GetAllDescendants(method TreeSearchMethod) []MessageNode {
    iterator := n.CreateIterator(method, 0)
    var res []MessageNode

    for iterator.Next() {
        res = append(res, *iterator.Value())
    }

    return res
}

func (n MessageNode) GetSiblings() []*MessageNode {
    var res []*MessageNode
    if n.Parent == nil {
        return res
    }

    iterator := n.Parent.Children.CreateIterator()
    for iterator.Next() {
        if val, exists := iterator.Value(); exists && val.Segment != n.Segment {
            res = append(res, val)
        }
    }

    return res
}

func (n MessageNode) CreateChildrenIterator() OrderedMapIterator[*MessageNode] {
    return n.Children.CreateIterator()
}

func (n MessageNode) CreateSiblingIterator() OrderedMapIterator[*MessageNode] {
    if n.Parent == nil {
        rootMap := CreateOrderedMap[*MessageNode]()
        rootMap.Set(n.Segment, &n)

        return rootMap.CreateIterator()
    }

    return n.Parent.Children.CreateIterator()
}

func (n MessageNode) Length() int {
    return n.Children.Length()
}

func (n MessageNode) GetDetailsString() string {
    msg := strings.Builder{}

    totalMessages := n.MessageCount
    if n.Children.Length() > 0 {
        if totalMessages > 0 {
            msg.WriteString(fmt.Sprintf("(%d topics, %d messages)", n.Children.Length(), totalMessages))
        } else {
            msg.WriteString(fmt.Sprintf("(%d topics)", n.Children.Length()))
        }
    } else {
        if totalMessages > 0 {
            msg.WriteString(fmt.Sprintf("= %s", strconv.Quote(string(n.Payloads[0]))))
        } else {
            msg.WriteString("= <empty>")
        }
    }

    return msg.String()
}
