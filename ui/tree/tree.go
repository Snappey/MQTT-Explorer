package tree

import (
    "context"
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "mqtt-explorer/internal"
    "net/url"
    "strings"
    "time"
)

type Model struct {
    Url           *url.URL
    Subscriptions []string

    ctx           context.Context
    ready         bool
    incoming      <-chan mqtt.Message
    messages      internal.MessageTree
    rootNode      NodeModel
    selectedTopic string
}

var (
    titleStyle = func() lipgloss.Style {
        b := lipgloss.RoundedBorder()
        b.Right = "├"
        return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
    }()

    infoStyle = func() lipgloss.Style {
        b := lipgloss.RoundedBorder()
        b.Left = "┤"
        return titleStyle.Copy().BorderStyle(b)
    }()

    contentStyle = func() lipgloss.Style {
        return lipgloss.NewStyle().Bold(true).Blink(true)
    }()
)

func CreateTreeModel(ctx context.Context, url *url.URL, subscriptions []string, incomingMessages <-chan mqtt.Message) Model {
    messageTree := internal.CreateMessageTree(url.String())

    return Model{
        Url:           url,
        Subscriptions: subscriptions,

        ctx:           ctx,
        incoming:      incomingMessages,
        messages:      messageTree,
        rootNode:      CreateNodeModel(messageTree.Root),
        selectedTopic: "",
    }
}

func (m Model) waitForMessage() tea.Cmd {
    return func() tea.Msg {
        return <-m.incoming
    }
}

func (m Model) processMessages() {
    for {
        select {
        case msg := <-m.incoming:
            m.messages.AddMessage(msg)
        case <-m.ctx.Done():
            return
        }
    }
}

type TickMsg time.Time

func (m Model) doTick() tea.Cmd {
    return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

func (m Model) Init() tea.Cmd {
    go m.processMessages()

    return tea.Batch(m.doTick(), tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var (
        cmd  tea.Cmd
        cmds []tea.Cmd
    )

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "crtl+c", "q", "esc":
            return m, tea.Quit
        }
    case TickMsg:
        cmds = append(cmds, m.doTick())
    }

    m.rootNode, cmd = m.rootNode.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m Model) headerView() string {
    return titleStyle.Render(fmt.Sprintf("MQTT Explorer (%s)", m.messages.Root.Segment))
}

func (m Model) View() string {
    return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.rootNode.View(), strings.Repeat("-", 60))
}
