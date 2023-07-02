package internal

type TreeSearchMethod = byte

const (
    BFS TreeSearchMethod = iota
    DFS
)

type MessageNodeIterator struct {
    toSearch   []*MessageNode
    startDepth uint
    maxDepth   uint
    method     TreeSearchMethod
}

func (n *MessageNode) CreateIterator(method TreeSearchMethod, depth uint) MessageNodeIterator {
    return MessageNodeIterator{
        toSearch:   []*MessageNode{n},
        startDepth: n.Depth,
        maxDepth:   n.Depth + depth,
        method:     method,
    }
}

func (i *MessageNodeIterator) Next() bool {
    return len(i.toSearch) > 0
}

func (i *MessageNodeIterator) Value() *MessageNode {
    var node *MessageNode
    switch i.method {
    case BFS:
        node, i.toSearch = i.toSearch[0], i.toSearch[1:]
    case DFS:
        node, i.toSearch = i.toSearch[len(i.toSearch)-1], i.toSearch[:len(i.toSearch)-1]
    }

    if node.Depth < i.maxDepth || i.maxDepth == 0 {
        childrenIterator := node.Children.CreateIterator()
        for childrenIterator.Next() {
            child, exists := childrenIterator.Value()
            if !exists {
                continue
            }

            i.toSearch = append(i.toSearch, child)
        }
    }

    return node
}
