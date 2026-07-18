package schema_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/schema"
)

func TestButtonOmitsFieldsForOtherActions(t *testing.T) {
	t.Parallel()

	button := schema.Button{Label: "전체 시간표", Action: "webLink", WebLinkURL: "https://hyuabot.app/shuttle"}
	body, err := json.Marshal(button)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	for _, unexpected := range []string{"messageText", "phoneNumber", "blockId", "extra"} {
		if strings.Contains(string(body), unexpected) {
			t.Errorf("button JSON contains %q: %s", unexpected, body)
		}
	}
}
