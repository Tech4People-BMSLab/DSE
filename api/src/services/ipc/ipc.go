package ipc

import (
	"fmt"
	"os"
	"time"

	"bms.dse/src/common"
	"bms.dse/src/utils/gatekeeper"
	"bms.dse/src/utils/logutil"
	"github.com/nats-io/nats.go"
)

// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type Client = common.Client

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var logger  = logutil.NewLogger("IPC")
var gkReady = gatekeeper.NewGateKeeper(true)
var nc *nats.Conn
// ------------------------------------------------------------
// : Wrappers
// ------------------------------------------------------------
func Request(topic string, data []byte, d time.Duration) (*nats.Msg, error) {
	gkReady.Wait()
	logger.Info().Str("topic", topic).Msg("Request")
	return nc.Request(topic, data, d)
}

func Publish(topic string, data []byte) error {
	gkReady.Wait()
	logger.Info().Str("topic", topic).Msg("Publish")
	return nc.Publish(topic, data)
}

func Subscribe(topic string, handler func(*nats.Msg)) (*nats.Subscription, error) {
	gkReady.Wait()
	logger.Info().Str("topic", topic).Msg("Subscribe")
	return nc.Subscribe(topic, handler)
}

// ------------------------------------------------------------
// : IPC (Inter-Process Communication)
// ------------------------------------------------------------
func Init() {
	var err error

	logger.Info().
		Str("url", os.Getenv("NATS_URI")).
		Msg("Connecting")

	/// Connect
	nc, err = nats.Connect(
		os.Getenv("NATS_URI"),

		nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
			switch (err) {
				case nats.ErrSlowConsumer:{
					pending, _, err := s.Pending()
					if err != nil {
						fmt.Printf("couldn't get pending messages: %v", err)
						return
					}
	
					logger.Warn().
						Int("pending", pending).
						Str("subject", s.Subject).
						Msg("Slow consumer")
				}

				default:{
					logger.Error().Err(err).Msg("Error")
				}
			}
		}),

		nats.DisconnectHandler(func(nc *nats.Conn) {
			logger.Warn().Msg("Disconnected") 
		}),
		nats.ConnectHandler(func(nc *nats.Conn) {
			logger.Info().Msg("Connected")
			gkReady.Unlock()
		}),
	)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to NATS")
		return
	}

	go func() {
		// NOTE: This has been moved to crawler.go

		// Subscribe("bms.dse.users.updated"     , OnUpdate)
		// Subscribe("bms.dse.users.connected"   , OnConnect)
		// Subscribe("bms.dse.users.disconnected", OnDisconnect)
		// Subscribe("bms.dse.users.*.pong"      , OnPong)
		// Subscribe("bms.dse.users.pong"        , OnPong)

		// Subscribe("bms.dse.users.*.updated"     , OnUpdate)
		// Subscribe("bms.dse.users.*.connected"   , OnConnect)
	}()
}
