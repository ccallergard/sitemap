package sitemap

import (
	"encoding/xml"
	"errors"
	"fmt"
	link "github.com/gophercises/link/students/ccallergard"
	"net/http"
	"net/url"
	"strings"
)

// Map travels given http/https site, returning set of links found within the same site,
// as absolute URLs
func Map(rawurl string) ([]string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("URL not http/https")
	}

	siteset := make(map[string]bool)
	siteset[u.String()] = true

	linkQueue := []string{rawurl}
	for len(linkQueue) > 0 {
		nextVisit := linkQueue[0]
		linkQueue = linkQueue[1:]

		foundLinks, err := visit(nextVisit)
		if err != nil {
			fmt.Println(err)
		}

		nextLinks := checkFound(foundLinks, siteset)
		linkQueue = append(linkQueue, nextLinks...)
	}

	var links []string
	for l := range siteset {
		links = append(links, l)
	}

	return links, nil
}

// visit GETs given URL, and returns hrefs found in the response body.
// Links to outside the same host are ignored, and links are absolutized
func visit(rawurl string) ([]string, error) {
	//fmt.Println(rawurl)
	resp, err := http.Get(rawurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	links, err := link.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	// Getting this from last request since redirects may have been followed
	from := resp.Request.URL

	var hrefs []string
	for _, l := range links {
		u, err := url.Parse(l.Href)
		if err != nil {
			continue
		}
		if n, ok := filterAndNormalize(u, from); ok {
			hrefs = append(hrefs, n)
		}
	}

	return hrefs, nil
}

// filterAndNormalize returns u in absolute form as string,
// if it's within the same site as. Otherwise returns ("", false).
func filterAndNormalize(u *url.URL, from *url.URL) (string, bool) {
	if u.IsAbs() {
		// Skip links not within host
		if u.Host != from.Host || u.Scheme != from.Scheme {
			return "", false
		}
	} else {
		if len(u.Path) == 0 {
			return "", false
		}
		// Absolutize
		u.Host = from.Host
		// No leading slash: path is relative to current directory
		if u.Path[0] != '/' {
			if dirEnd := strings.LastIndex(from.Path, "/"); dirEnd != -1 {
				u.Path = from.Path[:dirEnd+1] + u.Path
			}
		}
	}

	// Ensure same scheme (for absolute links too perhaps)
	u.Scheme = from.Scheme
	// Ignore anchors
	u.Fragment = ""

	return u.String(), true
}

// Returns list of links of ls not in the visited set, while adding them to the set
func checkFound(links []string, visited map[string]bool) (next []string) {
	for _, x := range links {
		if !visited[x] {
			visited[x] = true
			next = append(next, x)
		}
	}
	return next
}

var (
	urlsetTag = xml.StartElement{Name: xml.Name{Local: "urlset"},
		Attr: []xml.Attr{xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: "http://www.sitemaps.org/schemas/sitemap/0.9"}}}
	urlTag = xml.StartElement{Name: xml.Name{Local: "url"}}
	locTag = xml.StartElement{Name: xml.Name{Local: "loc"}}
)

// Returns urls as XML sitemap string
func XMLSiteMap(urls []string) string {
	var sb strings.Builder
	enc := xml.NewEncoder(&sb)
	enc.Indent("", "  ")
	sb.WriteString(xml.Header)
	enc.EncodeToken(urlsetTag)
	for _, u := range urls {
		enc.EncodeToken(urlTag)
		enc.EncodeElement(u, locTag)
		enc.EncodeToken(urlTag.End())

	}
	enc.EncodeToken(urlsetTag.End())
	enc.Flush()
	return sb.String()
}
