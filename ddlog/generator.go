package ddlog

import (
	"bytes"
	"math/rand"
	"net/url"
	"time"
	"unicode"

	"github.com/gee-go/ddlog/ddlog/randutil"
)

const (
	ExampleLogLine        = `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
	ExampleLogLineUnicode = `127.0.0.1 - 日a本b語ç日ð本Ê語þ日¥本¼語i日© [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
)

// A Generator is used to create random messages.
type Generator struct {
	c *Config
	r *randutil.Rand

	// If true, Generate fields with full range of unicode chars.
	// Otherwise only use ascii.
	UseUnicode  bool
	PageChoices []string
}

// NewGenerator inits a Generator with the given config.
// If config is nil, it uses the default config.
func NewGenerator(c *Config) *Generator {
	if c == nil {
		c = NewConfig()
	}

	return &Generator{
		c:          c,
		r:          randutil.R,
		UseUnicode: true,
	}
}

func (g *Generator) SeedPageChoice(n int) {
	for i := 0; i < n; i++ {
		g.PageChoices = append(g.PageChoices, g.r.Alpha(g.r.IntRange(1, 20)))
	}

}

func (g *Generator) randURI() string {
	r := g.r
	var u bytes.Buffer

	// preset pages as opposed to random.
	if len(g.PageChoices) > 0 {

		u.WriteString("/")
		u.WriteString(g.r.SelectString(g.PageChoices...))
	}

	// 1 to 5 path components
	for i := 0; i < r.IntRange(1, 6); i++ {
		u.WriteString("/")
		u.WriteString(r.Alpha(r.IntRange(1, 3)))
	}

	// 0 - 5 random url params
	v := &url.Values{}
	for i := 0; i < r.Rand.Intn(6); i++ {
		v.Set(r.Alpha(r.IntRange(1, 5)), r.Alpha(r.IntRange(1, 5)))
	}

	if p := v.Encode(); len(p) > 0 {
		u.WriteString("?")
		u.WriteString(p)
	}

	return u.String()
}

func (g *Generator) randAuth(maxLen int) string {
	var auth string
	authLen := rand.Intn(maxLen)
	if g.UseUnicode {
		auth = string(g.r.Unicode(unicode.L, authLen))
	} else {
		auth = g.r.Alpha(authLen)
	}

	if len(auth) == 0 {
		auth = "-"
	}

	return auth
}

func (g *Generator) time(t time.Time) time.Time {
	// parse and format to have correct accuracy
	out, _ := time.Parse(g.c.TimeFormat, t.Format(g.c.TimeFormat))
	return out
}

// TestMsg always returns a message equivalent to ExampleLogLine.
func (g *Generator) TestMsg() *Message {
	t, _ := time.Parse(DefaultTimeFormat, "10/Oct/2000:13:55:36 -0700")
	return &Message{
		Remote: "127.0.0.1",
		Ident:  "-",
		Auth:   "frank",
		Time:   g.time(t),
		Method: "GET",
		URI:    "/apache_pb.gif",
		Proto:  "HTTP/1.0",
		Status: 200,
		Size:   2326,
	}
}

func (g *Generator) MsgWithPage(page string) *Message {
	m := g.TestMsg()
	m.URI = page
	return m
}

// RandMsg creates a random Message
func (g *Generator) RandMsg() *Message {

	return &Message{
		Remote: g.r.IPv4().String(),
		Ident:  "-",
		Auth:   g.randAuth(20),
		Time:   g.time(time.Now()),
		Method: g.r.SelectString("GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"),
		URI:    g.randURI(),
		Proto:  "HTTP/1.0",
		Status: g.r.SelectInt(200, 400, 201, 304, 401, 404, 500),
		Size:   g.r.IntRange(1<<8, 1<<26),
	}
}
