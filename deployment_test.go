package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Exercise real API decoding, cache initialization and page rendering without
// depending on the availability of the public upstream API.
func TestDeploymentRoutes(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/artists" {
			fmt.Fprint(w, `[{"id":1,"name":"Queen","members":["Freddie Mercury"],"creationDate":1970,"firstAlbum":"13-07-1973"}]`)
		} else if r.URL.Path == "/api/relation" {
			fmt.Fprint(w, `{"index":[{"id":1,"datesLocations":{"london-uk":["01-01-2026"]}}]}`)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer upstream.Close()
	oldClient := apiClient
	oldArtists, oldRelations := ArtistCache, RelationCache
	defer func() { apiClient = oldClient; ArtistCache = oldArtists; RelationCache = oldRelations }()
	apiClient = &http.Client{Transport: redirectTransport{base: upstream.URL}}
	if err := InitCache(); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/artist", ArtistHandler)
	mux.HandleFunc("/search", SearchHandler)
	for _, tc := range []struct {
		method, path string
		status       int
		contains     string
	}{
		{"GET", "/", 200, "Queen"},
		{"GET", "/artist?id=1", 200, "Queen"},
		{"GET", "/search?q=freddie", 200, `"name":"Queen"`},
		{"GET", "/search?q=nomatch", 200, "[]"},
		{"GET", "/artist?id=bad", 400, "Bad Request"},
		{"GET", "/artist?id=53", 400, "Bad Request"},
		{"GET", "/missing", 404, "404 Not Found"},
		{"POST", "/", 405, "Method Not Allowed"},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != tc.status || !strings.Contains(rec.Body.String(), tc.contains) {
				t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
			}
		})
	}
}

type redirectTransport struct{ base string }

func (rt redirectTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	req, err := http.NewRequestWithContext(r.Context(), r.Method, rt.base+r.URL.Path, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultTransport.RoundTrip(req)
}
