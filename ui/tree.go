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

    ctx                 context.Context
    ready               bool
    windowWidth         int
    windowHeight        int
    incoming            <-chan mqtt.Message
    mqttClient          mqtt.Client
    messages            internal.MessageTree
    rootNode            NodeModel
    selectedNodeDetails NodeDetailsModel
    selectedNodePublish NodePublishModel
}

var (
    treeStyle = func(width int, height int) lipgloss.Style {
        return lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).
            Width(width / 2).
            Height(height - 2).
            MaxHeight(height).
            AlignHorizontal(lipgloss.Left)
    }

    detailStyle = func(width int, height int) lipgloss.Style {
        border := lipgloss.NormalBorder()
        border.Left = ""
        border.Right = ""

        return lipgloss.NewStyle().BorderStyle(border).
            Width(width/2-4).
            Height(height/2-2).
            Padding(0, 2, 0, 2).
            MaxWidth(width/2 - 2).
            MaxHeight(height / 2).
            AlignHorizontal(lipgloss.Left)
    }
)

func CreateTreeModel(ctx context.Context, url *url.URL, subscriptions []string, mqttClient mqtt.Client, incomingMessages <-chan mqtt.Message) TreeModel {
    messageTree := internal.CreateMessageTree(url.String())

    return TreeModel{
        Url:           url,
        Subscriptions: subscriptions,

        ctx:                 ctx,
        incoming:            incomingMessages,
        messages:            messageTree,
        mqttClient:          mqttClient,
        rootNode:            CreateNodeModel(messageTree.Root),
        selectedNodeDetails: CreateNodeDetailsModel(messageTree.Root),
        selectedNodePublish: CreateNodePublishModel(messageTree.Root),
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
        case "R":
            return m, tea.ClearScreen
        }
    case TickMsg:
        cmds = append(cmds, m.doTick())
    case tea.WindowSizeMsg:
        m.windowHeight = msg.Height
        m.windowWidth = msg.Width
    case PublishMessage:
        m.mqttClient.Publish(msg.Topic, msg.QoS, msg.Retained, msg.Payload) // TODO: Convert this to tea.Cmd to report back when token has completed (or errored)
    }

    m.selectedNodeDetails, cmd = m.selectedNodeDetails.Update(msg)
    cmds = append(cmds, cmd)

    m.selectedNodePublish, cmd = m.selectedNodePublish.Update(msg)
    cmds = append(cmds, cmd)

    m.rootNode, cmd = m.rootNode.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m TreeModel) View() string {
    return fmt.Sprintf("%s",
        lipgloss.JoinHorizontal(lipgloss.Left,
            treeStyle(m.windowWidth, m.windowHeight).Render(m.rootNode.View()),
            lipgloss.JoinVertical(lipgloss.Center,
                detailStyle(m.windowWidth, m.windowHeight).Render(m.selectedNodeDetails.View()),
                detailStyle(m.windowWidth, m.windowHeight).Render(m.selectedNodePublish.View()),
            ),
        ),
    )
}
