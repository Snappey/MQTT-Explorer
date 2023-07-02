package internal

type MessageNodeCursor struct {
    Root         *MessageNode
    SelectedNode *MessageNode
    iterator     OrderedMapIterator[*MessageNode]
}

func (n *MessageNode) CreateCursor() MessageNodeCursor {
    rootMap := CreateOrderedMap[*MessageNode]()
    rootMap.Set(n.Segment, n)

    return MessageNodeCursor{
        Root:         n,
        SelectedNode: n,
        iterator:     rootMap.CreateIterator(), // Iterator needs to start on the SelectedNode level
    }
}

func (n *MessageNodeCursor) setSelectedNode() {
    if val, exists := n.iterator.Value(); exists {
        n.SelectedNode = val
    }
}

func (n *MessageNodeCursor) Next() {
    if n.iterator.Next() {
        n.setSelectedNode()
    }
}

func (n *MessageNodeCursor) Previous() {
    if n.iterator.Previous() {
        n.setSelectedNode()
    }
}

func (n *MessageNodeCursor) Top() {
    n.iterator.Reset()
    n.setSelectedNode()
}

func (n *MessageNodeCursor) Bottom() {
    n.iterator.End()
    n.setSelectedNode()
}

func (n *MessageNodeCursor) Down() {
    if n.SelectedNode.Children.Length() == 0 {
        return
    }

    n.iterator = n.SelectedNode.CreateChildrenIterator()
    n.iterator.Next()

    n.setSelectedNode()
}

func (n *MessageNodeCursor) Up() {
    if n.SelectedNode.Parent == nil {
        return
    }

    n.iterator = n.SelectedNode.Parent.CreateSiblingIterator()
    n.iterator.Next()

    n.setSelectedNode()
}

func (n *MessageNodeCursor) GetSelectedNodeSiblings() []*MessageNode {
    return n.SelectedNode.GetSiblings()
}

func (n *MessageNodeCursor) GetParents() []*MessageNode {
    var res []*MessageNode

    node := n.SelectedNode
    for node != nil {
        res = append(res, node.Parent)
        node = node.Parent
    }

    return res
}
