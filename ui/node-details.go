package ui

import (
    "fmt"
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/snappey/mqtt-explorer/internal"
    "strconv"
    "strings"
    "time"
)

type NodeDetailsModel struct {
    node         *internal.MessageNode
    viewport     viewport.Model
    ready        bool
    windowWidth  int
    windowHeight int
}

func CreateNodeDetailsModel(selectedNode *internal.MessageNode) NodeDetailsModel {
    return NodeDetailsModel{
        node: selectedNode,
    }
}

var viewportStyle = lipgloss.NewStyle().
    Blink(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Border(lipgloss.RoundedBorder()).
    Padding(1, 1, 1, 1)

type SetSelectedNode struct {
    node *internal.MessageNode
}

func (m NodeDetailsModel) Init() tea.Cmd {
    return nil
}

func (m NodeDetailsModel) Update(msg tea.Msg) (NodeDetailsModel, tea.Cmd) {
    var (
        cmd  tea.Cmd
        cmds []tea.Cmd
    )

    switch msg := msg.(type) {
    case SetSelectedNode:
        m.node = msg.node
        m.viewport.SetContent(m.ViewPayload())
    case tea.WindowSizeMsg:
        m.windowHeight = msg.Height/2 - 2
        m.windowWidth = msg.Width/2 - 4

        if !m.ready {
            m.viewport = viewport.New(m.windowWidth, 5)
            m.viewport.HighPerformanceRendering = false
            m.viewport.KeyMap = viewport.KeyMap{}
            m.viewport.SetContent(m.ViewPayload())
            m.ready = true
        } else {
            m.viewport.Width = m.windowWidth
            m.viewport.Height = m.windowHeight
        }

        if m.viewport.HighPerformanceRendering {
            cmds = append(cmds, viewport.Sync(m.viewport))
        }
    }

    // Handle keyboard and mouse events in the viewport
    m.viewport, cmd = m.viewport.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m NodeDetailsModel) ViewPayload() string {
    sb := strings.Builder{}

    if len(m.node.Payloads) > 0 {
        sb.WriteString(strconv.Quote(string(m.node.Payloads[0])))
    }

    return sb.String()
}

func (m NodeDetailsModel) View() string {
    sb := strings.Builder{}

    if m.node.Parent == nil {
        sb.WriteString(fmt.Sprintf("%s\n", m.node.Segment))
    } else {
        sb.WriteString(fmt.Sprintf("%s\n", m.node.Topic))
    }

    sb.WriteString(
        viewportStyle.Width(m.windowWidth - 6).Render(m.viewport.View()),
    )
    sb.WriteRune('\n')

    if len(m.node.Payloads) > 0 {
        sb.WriteString(lipgloss.NewStyle().
            AlignHorizontal(lipgloss.Right).
            Width(m.windowWidth - 4).
            Italic(true).
            Render(
                fmt.Sprintf("Last Message: %s", m.node.ReceivedAt.Format(time.StampMilli)),
            ),
        )
    }

    return sb.String()
}
