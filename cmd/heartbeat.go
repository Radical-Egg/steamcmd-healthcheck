package cmd

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Radical-Egg/steamcmd-healthcheck/internal"
	"github.com/spf13/cobra"
)

var (
	hbHost         string
	hbPort         int
	hbTimeoutSec   int
	hbRecvResponse bool
	hbPacketHex    string
)

var defaultPacket = []byte{
	0xff, 0xff, 0xff, 0xff,
	'T', 'S', 'o', 'u', 'r', 'c', 'e', ' ',
	'E', 'n', 'g', 'i', 'n', 'e', ' ',
	'Q', 'u', 'e', 'r', 'y',
	0x00,
}

var heartbeatCmd = &cobra.Command{
	Use:   "heartbeat",
	Short: "Send a UDP heartbeat packet",
	Long: `Send a UDP packet to a target host/port, 
by default this program will heartbeat with 127.0.0.1:27015`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.New(io.Discard, "", 0)
		if verbose {
			logger.SetOutput(os.Stderr)
		}

		packet := append([]byte(nil), defaultPacket...)

		if strings.TrimSpace(hbPacketHex) != "" {
			p, err := decodeHexPacket(hbPacketHex)
			logger.Printf("%s [DEBUG] decoding packet override from --packet-hex", time.Now().UTC().Format(time.RFC3339))
			if err != nil {
				logger.Printf("%s [ERROR] invalid --packet-hex: %v", time.Now().UTC().Format(time.RFC3339), err)
				return fmt.Errorf("--packet-hex invalid hex packet: %w", err)
			}
			packet = p
		}

		if strings.TrimSpace(hbHost) == "" {
			return fmt.Errorf("host cannot be empty")
		}
		if hbPort < 1 || hbPort > 65535 {
			return fmt.Errorf("port must be between 1 and 65535")
		}
		if hbTimeoutSec <= 0 {
			return fmt.Errorf("timeout must be > 0 seconds")
		}
		if len(packet) == 0 {
			return fmt.Errorf("packet cannot be empty")
		}

		hb := internal.Heartbeat{
			Packet:       packet,
			Port:         hbPort,
			Host:         hbHost,
			RecvResponse: hbRecvResponse,
			Timeout:      time.Duration(hbTimeoutSec) * time.Second,
		}

		logger.Printf(
			"%s [INFO] heartbeat starting host=%s port=%d timeout=%s recv_response=%t packet_len=%d",
			time.Now().UTC().Format(time.RFC3339),
			hb.Host,
			hb.Port,
			hb.Timeout,
			hb.RecvResponse,
			len(hb.Packet),
		)

		if err := internal.Send(hb); err != nil {
			logger.Printf(
				"%s [ERROR] heartbeat failed host=%s port=%d err=%v",
				time.Now().UTC().Format(time.RFC3339),
				hb.Host,
				hb.Port,
				err,
			)
			return err
		}

		logger.Printf(
			"%s [INFO] heartbeat succeeded host=%s port=%d",
			time.Now().UTC().Format(time.RFC3339),
			hb.Host,
			hb.Port,
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(heartbeatCmd)

	heartbeatCmd.Flags().StringVar(&hbHost, "host", "127.0.0.1", "Target host")
	heartbeatCmd.Flags().IntVar(&hbPort, "port", 27015, "Target UDP port")
	heartbeatCmd.Flags().IntVar(&hbTimeoutSec, "timeout", 2, "Timeout in seconds")
	heartbeatCmd.Flags().BoolVar(&hbRecvResponse, "recv-response", true, "Wait for a UDP response")
	heartbeatCmd.Flags().StringVar(&hbPacketHex, "packet-hex", "", "Packet bytes as hex string (optional override)")

	heartbeatCmd.Example = strings.TrimSpace(`
./program heartbeat
./program heartbeat --host=127.0.0.1 --port=27015 --timeout=10
./program heartbeat --recv-response=false
./program heartbeat --packet-hex=ffffffff54536f7572636520456e67696e6520517565727900`)
}

func decodeHexPacket(s string) ([]byte, error) {
	clean := strings.TrimSpace(s)
	clean = strings.TrimPrefix(clean, "0x")
	clean = strings.ReplaceAll(clean, " ", "")

	if len(clean)%2 != 0 {
		return nil, fmt.Errorf("hex string must have even length")
	}

	b, err := hex.DecodeString(clean)
	if err != nil {
		return nil, err
	}

	return b, nil
}
