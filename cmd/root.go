package cmd

import (
    "context"
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/snappey/mqtt-explorer/ui"
    "log"
    "net/url"
    "os"
    "time"

    "github.com/spf13/cobra"
)

var (
    scheme   string
    hostname string
    port     int

    topic string
)

var rootCmd = &cobra.Command{
    Use:   "mqtt-explorer",
    Short: "Connect to an MQTT Broker",
    Long:  `Connect to an MQTT Broker`,
    Run: func(cmd *cobra.Command, args []string) {
        brokerUrl, err := url.Parse(fmt.Sprintf("%s://%s:%d", scheme, hostname, port))
        if err != nil {
            log.Fatalf("failed to parse options to valid url: %s err=%v", fmt.Sprintf("%s://%s:%d", scheme, hostname, port), err)
        }

        client := mqtt.NewClient(mqtt.NewClientOptions().
            SetClientID("mqtt-explorer").
            SetKeepAlive(time.Second * 30).
            SetPingTimeout(time.Second * 30).
            AddBroker(brokerUrl.String()))

        connectToken := client.Connect()

        if !connectToken.WaitTimeout(time.Second * 30) {
            log.Fatalf("failed to connect to %s... (Server Connection TimedOut)", brokerUrl)
        }

        log.Printf("connected to %s..", brokerUrl)

        messages := make(chan mqtt.Message)
        subToken := client.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
            messages <- message
        })

        if !subToken.WaitTimeout(time.Second * 5) {
            log.Fatalf("failed to connect to %s timeout exceeded", brokerUrl)
        }

        log.Printf("subcribed to %s", topic)

        go func() {
            for {
                if !client.IsConnected() {
                    log.Fatalf("mqtt broker has disconnected")
                }

                time.Sleep(time.Second * 2)
            }
        }()

        if _, err := tea.NewProgram(ui.CreateTreeModel(context.TODO(), brokerUrl, []string{topic}, client, messages), tea.WithAltScreen()).Run(); err != nil {
            log.Fatalf("error processing tree model err=%v", err)
        }

        client.Disconnect(0)

        log.Printf("exiting..")
    },
}

func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().StringVar(&scheme, "scheme", "tcp", "protocol to use tcp, ssl, ws or wss")
    rootCmd.PersistentFlags().StringVar(&hostname, "hostname", "test.mosquitto.org", "hostname of the broker")
    rootCmd.PersistentFlags().IntVar(&port, "port", 1883, "port the broker is running on, typically tcp: 1883, ssl: 8883, ws: 8083 or wss: 8084")

    rootCmd.Flags().StringVar(&topic, "topic", "#", "topic to subscribe to on the broker")
}
