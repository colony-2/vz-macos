package vmapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tmc/vz-macos/vmapi"
)

func ExampleBackend() {
	var backend vmapi.Backend = orchardBackend{}
	ctx := context.Background()
	info, _ := backend.Info(ctx)
	journal := journalFor()
	registry := backendRegistry{"selected": {backend, "store-v1"}}
	instance, err := reconcile(ctx, registry, &journal)
	fmt.Println(instance.Ref.UID, err)
	events, _ := backend.Events(ctx, vmapi.EventRequest{Instance: instance.Ref, Limit: 2})
	fmt.Println(events.Items[0].Template.BootIdentity, events.Items[1].Kind)
	fmt.Println(selectConfiguration(info.Capabilities, false, true, true))
	// Output:
	// remote-uid <nil>
	// boot delivery-succeeded
	// <nil>
}

func ExampleBackend_restart() {
	ctx := context.Background()
	journal := journalFor()
	registry := backendRegistry{"selected": {orchardBackend{}, "store-v1"}}
	_, _ = reconcile(ctx, registry, &journal)
	data, _ := json.Marshal(journal)
	var reopened launchJournal
	_ = json.Unmarshal(data, &reopened)
	// Reopen the same configured domain with a new object; never redispatch Create.
	registry["selected"] = backendBinding{orchardBackend{}, "store-v1"}
	instance, err := reconcile(ctx, registry, &reopened)
	fmt.Println(instance.Ref.UID, err)
	registry["selected"] = backendBinding{orchardBackend{}, "different-store"}
	_, err = reconcile(ctx, registry, &reopened)
	fmt.Println(errors.Is(err, vmapi.ErrConflict))
	// Output:
	// remote-uid <nil>
	// true
}
