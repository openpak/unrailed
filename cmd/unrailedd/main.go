// Command unrailedd answers Unrailed!'s own web API (unrailed-online.com) for the
// Switch build (010042C011476000, update v1.0.12) when a console is redirected to OpenPak.
//
// Multiplayer itself is Photon Realtime (ns.exitgames.com, served by
// platforms/switch/games/photon); this service is only the developer's small HTTP API
// the client calls around it:
//
//	GET  /nintendoHasSubscription?ticket=<id token>  online gate; body "NoOnlineService" refuses
//	GET  /getKick?userid=&ticket=&buildNumber=      public-match block list ("id:name" lines)
//	GET  /notifyKick?id=&name=&userid=&ticket=&...  a host reports a player it kicked
//	POST /uploadAvatar                              multipart: ticket, userid, avatar (jpeg)
//	POST /uploadReplay                              multipart: userid, ticket, datane, replay, week
//
// The client dials https://unrailed-online.com:8443 or http://unrailed-online.com:8080, ports it
// names itself. The service listens on plain HTTP (-addr) behind whatever terminates TLS, and can
// terminate TLS itself (-tls-addr with -cert/-key) when nothing in front can take port 8443.
//
// The contract was read out of the Switch executable; see docs/protocol.md.
package main

import (
	"flag"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"openpak/unrailed/redact"
)

var (
	addr       = flag.String("addr", ":23500", "plain HTTP listen address (the client's :8080, or behind a TLS terminator)")
	tlsAddr    = flag.String("tls-addr", "", "TLS listen address (the client's :8443); empty disables")
	certFile   = flag.String("cert", "", "TLS certificate chain (leaf for unrailed-online.com + OpenPak CA)")
	keyFile    = flag.String("key", "", "TLS private key")
	healthAddr = flag.String("health-addr", ":23501", "health listen address")
)

func main() {
	flag.Parse()
	s := &server{ids: newIdentityClient()}

	go func() {
		h := http.NewServeMux()
		h.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { io.WriteString(w, "ok\n") })
		log.Printf("[unrailed] health on %s", *healthAddr)
		log.Fatal(http.ListenAndServe(*healthAddr, h))
	}()
	if *tlsAddr != "" {
		go func() {
			log.Printf("[unrailed] serving TLS on %s", *tlsAddr)
			log.Fatal(http.ListenAndServeTLS(*tlsAddr, *certFile, *keyFile, s))
		}()
	}
	log.Printf("[unrailed] serving HTTP on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, s))
}

type server struct {
	ids identityChecker
}

// noOnlineService is the one body the client compares against: equal means "this player has
// no online access", anything else lets it go online.
const noOnlineService = "NoOnlineService"

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/nintendoHasSubscription":
		s.subscription(w, r)
	case "/getKick":
		// No public-match blocks: the block list the client reads is "id:name" lines, and an
		// empty body blocks nobody. Reports are not tallied into blocks (see notifyKick).
		log.Printf("[unrailed] getKick %s", logURL(r.URL))
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	case "/notifyKick":
		log.Printf("[unrailed] notifyKick %s", logURL(r.URL))
		w.WriteHeader(http.StatusOK)
	case "/uploadAvatar", "/uploadReplay":
		s.upload(w, r)
	default:
		log.Printf("[unrailed] !! unhandled %s %s", r.Method, logURL(r.URL))
		http.NotFound(w, r)
	}
}

// subscription is the online gate. The ticket is the player's Nintendo id token, which on an
// OpenPak console OpenPak itself issued, so nx-baas is the one to ask. A rejected or banned
// token gets the refusal body; an nx-baas outage lets the player through, since the gate only
// decides whether online play is offered and Photon will still refuse a banned account.
func (s *server) subscription(w http.ResponseWriter, r *http.Request) {
	tok := r.URL.Query().Get("ticket")
	verdict := "OK"
	switch err := s.ids.check(tok); {
	case err == nil:
	case tok == "", isRefusal(err):
		verdict = noOnlineService
	default:
		log.Printf("[unrailed] subscription: identity lookup failed (%s), letting through", redact.Err(err))
	}
	log.Printf("[unrailed] subscription ticket=%s -> %s", tagOf(tok), verdict)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, verdict)
}

// upload accepts avatar and replay uploads and keeps nothing: OpenPak shows neither, and a
// replay is the player's run, not ours to store. The fields are logged by name and size only.
func (s *server) upload(w http.ResponseWriter, r *http.Request) {
	fields := []string{}
	if mt, params, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err == nil && strings.HasPrefix(mt, "multipart/") {
		mr := multipart.NewReader(io.LimitReader(r.Body, 64<<20), params["boundary"])
		for {
			p, err := mr.NextPart()
			if err != nil {
				break
			}
			n, _ := io.Copy(io.Discard, p)
			fields = append(fields, p.FormName()+"("+strconv.FormatInt(n, 10)+")")
		}
	} else {
		n, _ := io.Copy(io.Discard, io.LimitReader(r.Body, 64<<20))
		fields = append(fields, "body("+strconv.FormatInt(n, 10)+")")
	}
	log.Printf("[unrailed] %s %s fields=%s", r.Method, r.URL.Path, strings.Join(fields, ","))
	w.WriteHeader(http.StatusOK)
}

// logURL renders a request URL with every query value that names or authenticates a player
// replaced by a tag.
func logURL(u *url.URL) string {
	q := u.Query()
	for k, vs := range q {
		switch strings.ToLower(k) {
		case "buildnumber":
			continue
		}
		for i, v := range vs {
			vs[i] = tagOf(v)
		}
		q[k] = vs
	}
	if len(q) == 0 {
		return u.Path
	}
	return u.Path + "?" + q.Encode()
}

func tagOf(v string) string {
	if v == "" {
		return "-"
	}
	return redact.Tag(v)
}
