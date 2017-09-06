package stack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	vault "github.com/hashicorp/vault/api"
)

func TestMultipleVaultClients_DoNotMerge(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintln(w, `{"initialized":true,"sealed":false,"standby":true,"server_time_utc":1504683364,"version":"0.7.3","cluster_name":"vault-test","cluster_id":"test"}`)
	}))
	defer ts1.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"initialized":true,"sealed":false,"standby":false,"server_time_utc":1504683364,"version":"0.7.3","cluster_name":"vault-test","cluster_id":"test"}`)
	}))
	defer ts2.Close()

	vault1, err := vault.NewClient(nil)
	if err != nil {
		t.Error("error initialising vault1: ", err)
	}
	vault1.SetAddress(ts1.URL)
	health1, err := vault1.Sys().Health()
	if err == nil {
		t.Error("expecting error getting vault1's status")
	}

	vault2, err := vault.NewClient(nil)
	if err != nil {
		t.Error("error initialising vault2: ", err)
	}
	vault2.SetAddress(ts2.URL)
	health2, err := vault2.Sys().Health()
	if err != nil {
		t.Error("error getting vault2's status: ", err)
	}

	t.Logf("status health vault1=%v vault2=%v", health1, health2)

}
