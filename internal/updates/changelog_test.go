package updates

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	wupdater "github.com/wailsapp/wails/v3/pkg/updater"
)

func coordinatorWithDoer(d HTTPDoer) *Coordinator {
	return New(&fakeEngine{rel: (*wupdater.Release)(nil)}, d, false, zerolog.Nop())
}

func TestChangelog_ReturnsReleaseBody(t *testing.T) {
	c := coordinatorWithDoer(doerFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://api.github.com/repos/"+RepoSlug+"/releases/latest" {
			t.Fatalf("unexpected URL %s", r.URL)
		}
		return jsonResponse(200, `{"body":"## What changed\n- stuff"}`), nil
	}))

	got := c.Changelog(context.Background())
	if got != "## What changed\n- stuff" {
		t.Fatalf("Changelog = %q, want the release body", got)
	}
}

func TestChangelog_OfflineFallsBack(t *testing.T) {
	c := coordinatorWithDoer(doerFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("dial tcp: no route to host")
	}))

	if got := c.Changelog(context.Background()); got != ChangelogUnavailable {
		t.Fatalf("Changelog = %q, want %q", got, ChangelogUnavailable)
	}
}

func TestChangelog_HTTPErrorFallsBack(t *testing.T) {
	c := coordinatorWithDoer(doerFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(404, `{"message":"Not Found"}`), nil
	}))

	if got := c.Changelog(context.Background()); got != ChangelogUnavailable {
		t.Fatalf("Changelog = %q, want %q", got, ChangelogUnavailable)
	}
}

func TestChangelog_EmptyBodyGivesPlaceholder(t *testing.T) {
	c := coordinatorWithDoer(doerFunc(func(*http.Request) (*http.Response, error) {
		return jsonResponse(200, `{"body":""}`), nil
	}))

	got := c.Changelog(context.Background())
	if got == "" || got == ChangelogUnavailable {
		t.Fatalf("Changelog = %q, want a non-empty 'no notes' placeholder", got)
	}
}
