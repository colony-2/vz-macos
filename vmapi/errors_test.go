package vmapi_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/tmc/vz-macos/vmapi"
)

func TestErrorInspection(t *testing.T) {
	tests := []struct {
		name  string
		input vmapi.ErrorCode
		want  error
	}{
		{"invalid", vmapi.CodeInvalidArgument, vmapi.ErrInvalidArgument},
		{"missing", vmapi.CodeNotFound, vmapi.ErrNotFound},
		{"conflict", vmapi.CodeConflict, vmapi.ErrConflict},
		{"unsupported", vmapi.CodeUnsupported, vmapi.ErrUnsupported},
		{"unavailable", vmapi.CodeUnavailable, vmapi.ErrUnavailable},
		{"uncertain", vmapi.CodeUnknownOutcome, vmapi.ErrUnknownOutcome},
		{"expired", vmapi.CodeCursorExpired, vmapi.ErrCursorExpired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &vmapi.Error{Code: tt.input, RequestID: "request"}
			wrapped := fmt.Errorf("submit: %w", input)
			for _, other := range tests {
				if got := errors.Is(wrapped, other.want); got != (other.input == tt.input) {
					t.Errorf("Is(%v) = %v", other.want, got)
				}
			}
			if errors.Is(wrapped, errors.New(input.Error())) {
				t.Fatal("matched unrelated target")
			}
			var detail *vmapi.Error
			if !errors.As(wrapped, &detail) || detail.RequestID != "request" {
				t.Fatalf("lost correlation: %v", detail)
			}
			if input.Error() != tt.want.Error() {
				t.Fatalf("message = %q", input.Error())
			}
		})
	}
	for _, code := range []vmapi.ErrorCode{"", "future"} {
		input := &vmapi.Error{Code: code}
		for _, tt := range tests {
			if errors.Is(input, tt.want) {
				t.Errorf("unknown %q matched %v", code, tt.want)
			}
		}
		if input.Is(nil) {
			t.Fatal("unknown matched nil")
		}
	}
	if got := (&vmapi.Error{Message: "safe message"}).Error(); got != "safe message" {
		t.Fatal(got)
	}
	var nilError *vmapi.Error
	if nilError.Is(vmapi.ErrUnavailable) || nilError.Error() != "vmapi error" {
		t.Fatal("nil error inspection")
	}
}

func TestErrorFields(t *testing.T) {
	typ := reflect.TypeOf(vmapi.Error{})
	want := map[string]reflect.Type{
		"Code":      reflect.TypeOf(vmapi.ErrorCode("")),
		"RequestID": reflect.TypeOf(vmapi.RequestID("")),
		"Message":   reflect.TypeOf(""),
	}
	if typ.NumField() != len(want) {
		t.Fatal("unexpected error fields")
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if want[field.Name] != field.Type {
			t.Errorf("unexpected field %s: %v", field.Name, field.Type)
		}
	}
}
