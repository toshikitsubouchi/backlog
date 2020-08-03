package backlog

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/pkg/errors"
)

func getTestResolution() *Resolution {
	return &Resolution{
		ID:   Int(0),
		Name: String("対応済み"),
	}
}

func TestGetResolutions(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/resolutions", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `[{"id": 0, "name": "対応済み"}]`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetResolutions()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := []*Resolution{getTestResolution()}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse, errors.New(pretty.Compare(want, expected)))
	}
}

func TestGetResolutionsFailed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/resolutions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	if _, err := client.GetResolutions(); err == nil {
		t.Fatal("expected an error but got none")
	}
}
