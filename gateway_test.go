package kenan

import (
	"encoding/json"
	"testing"
)

// TestSetSessionAcceptsStructPointer 验证 gateway 可直接传入 JWT Claims 等结构体指针作为会话。
func TestSetSessionAcceptsStructPointer(t *testing.T) {
	type testSession struct {
		AccountName string `json:"account_name"`
	}

	args := new(Args)
	setSession(args, &testSession{AccountName: "admin"})

	sessionData, ok := args.SessionData.([]byte)
	if !ok || !json.Valid(sessionData) {
		t.Fatalf("SessionData should contain valid JSON bytes, got %#v", args.SessionData)
	}
	if args.Session["account_name"] != "admin" {
		t.Fatalf("Session should retain decoded claims, got %#v", args.Session)
	}
}
