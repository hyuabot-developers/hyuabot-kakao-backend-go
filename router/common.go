package router

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/backend"
	"github.com/hyuabot-developers/hyuabot-kakao-backend-go/schema"
)

const (
	BackendClientLocal = "backendClient"
	queryTimeout       = 3500 * time.Millisecond
	serviceTimezone    = "Asia/Seoul"
	serviceUTCOffset   = 9 * 60 * 60
	shortTimeLength    = 5
	noArrivalText      = "현재 도착 정보가 없어요"
)

func backendClient(ctx fiber.Ctx) (backend.Client, error) {
	client, ok := ctx.Locals(BackendClientLocal).(backend.Client)
	if !ok {
		return nil, errors.New("backend client not found")
	}
	return client, nil
}

func queryContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), queryTimeout)
}

func bindSkillPayload(ctx fiber.Ctx) (*schema.SkillPayload, error) {
	payload := new(schema.SkillPayload)
	if err := ctx.Bind().JSON(payload); err != nil {
		return nil, fmt.Errorf("decode skill payload: %w", err)
	}
	return payload, nil
}

func badRequest(ctx fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
}

func backendFailure(ctx fiber.Ctx, feature string, retryMessage string, err error) error {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Error("backend query failed", "feature", feature, "error", err)
	return ctx.JSON(schema.TextResponse(
		"지금은 정보를 불러오지 못했어요. 잠시 후 다시 시도해 주세요.",
		navigationQuickReplies(retryMessage),
	))
}

func navigationQuickReplies(retryMessage string) []schema.QuickReply {
	messages := []string{retryMessage, "셔틀", "버스", "지하철", "학식"}
	replies := make([]schema.QuickReply, 0, len(messages))
	for index, message := range messages {
		label := message
		if index == 0 {
			label = "새로고침"
		}
		replies = append(replies, schema.QuickReply{
			Label:       label,
			Action:      "message",
			MessageText: message,
		})
	}
	return replies
}

func appButton(path string, label string) schema.Button {
	return schema.Button{
		Label:      label,
		Action:     "webLink",
		WebLinkURL: "https://hyuabot.app" + path,
	}
}

func currentServiceTime() time.Time {
	location, err := time.LoadLocation(serviceTimezone)
	if err != nil {
		location = time.FixedZone("KST", serviceUTCOffset)
	}
	return time.Now().In(location)
}

func shortTime(value string) string {
	if len(value) >= shortTimeLength {
		return value[:shortTimeLength]
	}
	return value
}
