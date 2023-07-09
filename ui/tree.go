package ui

import (
    "context"
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/snappey/mqtt-explorer/internal"
    "net/url"
    "time"
)

type TreeModel struct {
    Url           *url.URL
    Subscriptions []string

    ctx          context.Context
    ready        bool
    width        int
    height       int
    incoming     <-chan mqtt.Message
    messages     internal.MessageTree
    rootNode     NodeModel
    selectedNode NodeDetailsModel
}

var (
    titleStyle = func() lipgloss.Style {
        b := lipgloss.RoundedBorder()
        b.Right = "├"
        return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
    }()

    treeStyle = func(width int, height int) lipgloss.Style {
        return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).
            Width(width / 2).
            Height(height - 2).
            MaxHeight(height).
            AlignHorizontal(lipgloss.Left)
    }

    detailStyle = func(width int, height int) lipgloss.Style {
        return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).
            Width(width/2 - 4).
            Height(height - 2).
            MaxHeight(height).
            AlignHorizontal(lipgloss.Left)
    }
)

func CreateTreeModel(ctx context.Context, url *url.URL, subscriptions []string, incomingMessages <-chan mqtt.Message) TreeModel {
    messageTree := internal.CreateMessageTree(url.String())

    return TreeModel{
        Url:           url,
        Subscriptions: subscriptions,

        ctx:          ctx,
        incoming:     incomingMessages,
        messages:     messageTree,
        rootNode:     CreateNodeModel(messageTree.Root),
        selectedNode: CreateNodeDetailsModel(messageTree.Root),
    }
}

func (m TreeModel) waitForMessage() tea.Cmd {
    return func() tea.Msg {
        return <-m.incoming
    }
}

func (m TreeModel) processMessages() {
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

func (m TreeModel) doTick() tea.Cmd {
    return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
        return TickMsg(t)
    })
}

func (m TreeModel) Init() tea.Cmd {
    go m.processMessages()

    return tea.Batch(m.doTick(), tea.EnterAltScreen)
}

func (m TreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
    case tea.WindowSizeMsg:
        m.height = msg.Height
        m.width = msg.Width
    }

    m.selectedNode, cmd = m.selectedNode.Update(msg)
    cmds = append(cmds, cmd)

    m.rootNode, cmd = m.rootNode.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m TreeModel) headerView() string {
    return titleStyle.Render(fmt.Sprintf("MQTT Explorer (%s)", m.messages.Root.Segment))
}

func (m TreeModel) View() string {
    return fmt.Sprintf("%s",
        lipgloss.JoinHorizontal(lipgloss.Left,
            treeStyle(m.width, m.height).Render(m.rootNode.View()),
            detailStyle(m.width, m.height).Render(m.selectedNode.View()),
        ),
    )
}