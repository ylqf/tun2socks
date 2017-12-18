package tun2socks

import (
	"log"
	"net"
	"sync"

	"github.com/FlowerWrong/netstack/tcpip/stack"
	"github.com/FlowerWrong/tun2socks/configure"
	"github.com/FlowerWrong/tun2socks/dns"
	"github.com/FlowerWrong/tun2socks/util"
	"github.com/FlowerWrong/water"
)

type App struct {
	FakeDNS  *dns.DNS
	Cfg      *configure.AppConfig
	Proxies  *configure.Proxies
	S        *stack.Stack
	Ifce     *water.Interface
	HookPort uint16
	WG       sync.WaitGroup
}

func (app *App) NewTun() *App {
	NewTun(app)
	return app
}

func (app *App) AddRoutes() *App {
	name := app.Ifce.Name()
	for _, val := range app.Cfg.Route.V {
		_, subnet, _ := net.ParseCIDR(val)
		if subnet != nil {
			util.AddNetRoute(name, subnet)
		} else {
			util.AddHostRoute(name, val)
		}
	}
	return app
}

func (app *App) Config(configFile string) *App {
	// parse config
	app.Cfg = new(configure.AppConfig)
	err := app.Cfg.Parse(configFile)
	if err != nil {
		log.Fatal("Get default proxy failed", err)
	}

	if app.Cfg.DNS.DNSMode == "fake" {
		app.FakeDNS, err = dns.NewFakeDNSServer(app.Cfg)
		if err != nil {
			log.Fatal("New fake dns server failed", err)
		}
	}

	app.Proxies, err = configure.NewProxies(app.Cfg.Proxy)
	if err != nil {
		log.Fatalln("New proxies failed", err)
	}

	return app
}
