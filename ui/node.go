package ui

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/snappey/mqtt-explorer/internal"
    "strings"
)

type NodeModel struct {
    node     *internal.MessageNode
    cursor   internal.MessageNodeCursor
    height   int
    expanded bool
}

var selectedStyle = func() lipgloss.Style {
    return lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FAFAFA")).
        Background(lipgloss.Color("#7D56F4"))
}

var selectedTopicStyle = func() lipgloss.Style {
    return lipgloss.NewStyle().
        Italic(true).
        Foreground(lipgloss.Color("#BABABA"))
}

func CreateNodeModel(node *internal.MessageNode) NodeModel {
    return NodeModel{
        node:     node,
        cursor:   node.CreateCursor(),
        height:   15,
        expanded: true,
    }
}

func (m NodeModel) Init() tea.Cmd {
    return nil
}

func (m NodeModel) Update(msg tea.Msg) (NodeModel, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up":
            m.cursor.Previous()
        case "down":
            m.cursor.Next()
        case "left":
            m.cursor.Up()
        case "right":
            m.cursor.Down()
        }

        cmds = append(cmds, func() tea.Msg {
            return SetSelectedNode{node: m.cursor.SelectedNode}
        })
    }

    return m, tea.Batch(cmds...)
}

func (m NodeModel) RenderRootNode() string {
    var messageCount int
    if m.cursor.SelectedNode.Parent == nil {
        messageCount = m.cursor.SelectedNode.MessageCount
    } else {
        messageCount = m.cursor.SelectedNode.Parent.MessageCount
    }

    parentTopics := strings.Split(m.cursor.SelectedNode.Topic, "/")
    return fmt.Sprintf("%s (%d messages)", strings.Join(parentTopics[:len(parentTopics)-1], "/"), messageCount)
}

func (m NodeModel) RenderNodes() string {
    sb := strings.Builder{}

    if m.cursor.SelectedNode == nil {
        sb.WriteString("|- Empty...")

        return sb.String()
    }

    i := m.cursor.SelectedNode.CreateSiblingIterator()
    i.SkipUntil(m.cursor.SelectedNode.Segment)
    i.Rewind(5)

    drawn := 0
    for i.Next() {
        topic, child, exists := i.Pair()
        if !exists {
            continue
        }

        if drawn > m.height {
            break
        }

        msg := strings.Builder{}
        msg.WriteString(fmt.Sprintf("-> %s %s", topic, child.GetDetailsString()))

        if child.Topic == m.cursor.SelectedNode.Topic {
            sb.WriteString(selectedStyle().Render(msg.String()))

            if child.Children.Length() > 0 {
                iterator := child.CreateChildrenIterator()

                sb.WriteRune('\n')
                for j := 0; j < 3 && iterator.Next(); j++ {
                    sb.WriteString("   |-> ")
                    if subtopic, exists := iterator.Value(); exists {
                        sb.WriteString(selectedTopicStyle().Render(fmt.Sprintf("%s %s", subtopic.Segment, subtopic.GetDetailsString())))
                    } else {
                        sb.WriteString("<MISSING NODE>")
                    }

                    if iterator.HasNext() {
                        sb.WriteRune('\n')
                    }
                }

                if iterator.HasNext() {
                    sb.WriteString("   |-> ")
                    sb.WriteString(selectedTopicStyle().Render(fmt.Sprintf("... %d Hidden Topic(s)", iterator.Remaining())))
                }
            }
        } else {
            sb.WriteString(msg.String())
        }

        if drawn < m.height {
            sb.WriteRune('\n')
        }

        drawn += 1
    }

    return sb.String()
}

func (m NodeModel) View() string {
    sb := strings.Builder{}

    sb.WriteString(m.RenderRootNode())

    if m.expanded {
        sb.WriteString("\n")

        sb.WriteString(m.RenderNodes())
    }

    return sb.String()
}
