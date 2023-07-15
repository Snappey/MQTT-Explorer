package ui

import (
    "fmt"
    "github.com/charmbracelet/bubbles/textarea"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/snappey/mqtt-explorer/internal"
    "strings"
)

type NodePublishModel struct {
    node         *internal.MessageNode
    textarea     textarea.Model
    windowWidth  int
    windowHeight int
}

func CreateNodePublishModel(selectedNode *internal.MessageNode) NodePublishModel {
    ta := textarea.New()
    ta.Placeholder = "Send a message..."
    ta.Focus()

    ta.Prompt = "┃ "
    ta.CharLimit = 4096

    // Remove cursor line styling
    ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
    ta.ShowLineNumbers = false

    return NodePublishModel{
        node:     selectedNode,
        textarea: ta,
    }
}

type PublishMessage struct {
    Topic    string
    Payload  []byte
    QoS      byte
    Retained bool
}

func (m NodePublishModel) Init() tea.Cmd {
    return nil
}

func (m NodePublishModel) Update(msg tea.Msg) (NodePublishModel, tea.Cmd) {
    var (
        cmd  tea.Cmd
        cmds []tea.Cmd
    )

    m.textarea, cmd = m.textarea.Update(msg)
    cmds = append(cmds, cmd)

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyTab:
            if m.textarea.Focused() {
                m.textarea.Blur()
            } else {
                m.textarea.Focus()
            }
        case tea.KeyEnter:
            if !m.textarea.Focused() || m.textarea.Value() == "" {
                break
            }

            cmds = append(cmds, func() tea.Msg {
                return PublishMessage{
                    Topic:    m.node.Topic,
                    Payload:  []byte(m.textarea.Value()),
                    QoS:      0,
                    Retained: false,
                }
            })
        }
    case SetSelectedNode:
        m.node = msg.node
    case tea.WindowSizeMsg:
        m.windowHeight = msg.Height/2 - 2
        m.windowWidth = msg.Width/2 - 4

        m.textarea.SetWidth(m.windowWidth - 6)
        m.textarea.SetHeight(6)

        m.textarea.Blur()
    }

    return m, tea.Batch(cmds...)
}

func (m NodePublishModel) View() string {
    sb := strings.Builder{}

    sb.WriteString(fmt.Sprintf("Publish to %s\n", m.node.Topic))
    sb.WriteString(m.textarea.View())

    return sb.String()
}
