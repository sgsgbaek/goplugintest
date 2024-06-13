package main

import (
	"fmt"
	"goplugintest/internal/plugin"
	"strings"

	rolloutsPlugin "github.com/argoproj/argo-rollouts/rollout/trafficrouting/plugin/rpc"
	goPlugin "github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
)

var handshakeConfig = goPlugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ARGO_ROLLOUTS_RPC_PLUGIN",
	MagicCookieValue: "trafficrouter",
}

func main() {
	result := plugin.DoAction()
	fmt.Print(result)

	logCtx := log.WithFields(log.Fields{"plugin": "goplugintest"})

	setLogLevel("debug")
	log.SetFormatter(createFormatter("text"))

	rpcPluginImp := &plugin.RpcPlugin{
		LogCtx: logCtx,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]goPlugin.Plugin{
		"RpcTrafficRouterPlugin": &rolloutsPlugin.RpcTrafficRouterPlugin{Impl: rpcPluginImp},
	}

	logCtx.Debug("message from plugin", "foo", "bar")

	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

func createFormatter(logFormat string) log.Formatter {
	var formatType log.Formatter
	switch strings.ToLower(logFormat) {
	case "json":
		formatType = &log.JSONFormatter{}
	case "text":
		formatType = &log.TextFormatter{
			FullTimestamp: true,
		}
	default:
		log.Infof("Unknown format: %s. Using text logformat", logFormat)
		formatType = &log.TextFormatter{
			FullTimestamp: true,
		}
	}

	return formatType
}

func setLogLevel(logLevel string) {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}
