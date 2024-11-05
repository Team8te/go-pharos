package ds

import "net/url"

type TaskMove struct {
	ID    int
	URL   *url.URL
	Files []string
}
