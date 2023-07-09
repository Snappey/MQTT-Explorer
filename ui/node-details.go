package ui

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/glamour"
    "github.com/snappey/mqtt-explorer/internal"
    "strings"
)

type NodeDetailsModel struct {
    node *internal.MessageNode
}

var markdownRenderer, _ = glamour.NewTermRenderer(
    glamour.WithAutoStyle(),
    glamour.WithPreservedNewLines(),
)

func CreateNodeDetailsModel(selectedNode *internal.MessageNode) NodeDetailsModel {
    return NodeDetailsModel{
        node: selectedNode,
    }
}

type SetSelectedNode struct {
    node *internal.MessageNode
}

func (m NodeDetailsModel) Init() tea.Cmd {
    return nil
}

func (m NodeDetailsModel) Update(msg tea.Msg) (NodeDetailsModel, tea.Cmd) {
    switch msg := msg.(type) {
    case SetSelectedNode:
        m.node = msg.node
    }

    return m, nil
}

func (m NodeDetailsModel) View() string {
    sb := strings.Builder{}

    sb.WriteString(fmt.Sprintf("%s\n", m.node.Topic))

    if m.node.Children.Length() > 0 {
        iterator := m.node.Children.CreateIterator()
        for iterator.Next() {
            if subtopic, exists := iterator.Value(); exists {
                sb.WriteString(fmt.Sprintf("|-- %s", subtopic.Segment))
            }

            if iterator.HasNext() {
                sb.WriteRune('\n')
            }
        }
    }

    sb.WriteRune('\n')

    if len(m.node.Payloads) > 0 {
        sb.WriteString("Payloads:")
        for _, payload := range m.node.Payloads {
            res, _ := markdownRenderer.Render(
                fmt.Sprintf("`%s`", payload),
            )

            sb.WriteString(res)
            sb.WriteRune('\n')
        }
    }

    return sb.String()
}
