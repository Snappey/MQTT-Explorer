package internal

import (
    "fmt"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/muesli/reflow/indent"
    "strings"
    "time"
)

const PayloadHistoryMax = 5

type MessageTree struct {
    Root *MessageNode
}

func CreateMessageTree(rootNode string) MessageTree {
    return MessageTree{
        Root: &MessageNode{
            Topic:    "",
            Segment:  rootNode,
            Depth:    0,
            Children: CreateOrderedMap[*MessageNode](),
            Payloads: [][]byte{},
            Parent:   nil,
        },
    }
}

func (t *MessageTree) AddMessage(message mqtt.Message) {
    topic, payload := message.Topic(), message.Payload()
    path := strings.Split(topic, "/")

    var node = t.Root
    for i, segment := range path {
        node.MessageCount += 1

        if _, exists := node.Children.Get(segment); !exists {
            topic := segment
            if node.Parent != nil {
                topic = fmt.Sprintf("%s/%s", node.Topic, segment)
            }

            node.Children.Set(segment, &MessageNode{
                Topic:    topic,
                Segment:  segment,
                Depth:    uint(i),
                Children: CreateOrderedMap[*MessageNode](),
                Payloads: [][]byte{},
                Parent:   node,
            })
        }

        if i == len(path)-1 {
            child, _ := node.Children.Get(segment)

            child.Payloads = append(child.Payloads, payload)
            if len(child.Payloads) > PayloadHistoryMax {
                child.Payloads = child.Payloads[1:]
            }

            child.ReceivedAt = time.Now()
            child.MessageCount += 1

            node.Children.Set(segment, child)
        } else {
            node, _ = node.Children.Get(segment)
        }
    }
}

func (t *MessageTree) Render() string {
    s := strings.Builder{}
    var node *MessageNode
    toSearch := []*MessageNode{t.Root}
    for len(toSearch) > 0 {
        node, toSearch = toSearch[len(toSearch)-1], toSearch[:len(toSearch)-1]

        if len(node.Payloads) > 0 {
            payload := node.Payloads[len(node.Payloads)-1]
            s.WriteString(fmt.Sprintf("%s: %s\n", indent.String(node.Segment, node.Depth*4), trimByteArray(payload, 120)))
        } else {
            s.WriteString(fmt.Sprintf("%s\n", indent.String(fmt.Sprintf("|- %s", node.Segment), node.Depth*4)))
        }

        iterator := node.Children.CreateIterator()
        for iterator.Next() {
            if child, exists := iterator.Value(); exists {
                toSearch = append(toSearch, child)
            }
        }
    }
    return s.String()
}

func (t *MessageTree) GetNode(path string) (*MessageNode, error) {
    segments := strings.Split(path, "/")

    node := t.Root
    var exists bool
    for _, segment := range segments {
        node, exists = node.Children.Get(segment)
        if !exists {
            return nil, fmt.Errorf("failed to get node missing segment path=%s segment=%s", path, segment)
        }
    }

    return node, nil
}

func trimByteArray(bytes []byte, cap int) []byte {
    if len(bytes) > cap {
        return bytes[0:cap]
    }
    return bytes
}
