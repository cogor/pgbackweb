package webhooks

import (
	"github.com/eduardolat/pgbackweb/internal/util/pathutil"
	"github.com/eduardolat/pgbackweb/internal/view/web/component"
	"github.com/eduardolat/pgbackweb/internal/view/web/respondhtmx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	nodx "github.com/nodxdev/nodxgo"
	htmx "github.com/nodxdev/nodxgo-htmx"
	lucide "github.com/nodxdev/nodxgo-lucide"
)

func (h *handlers) deleteWebhookHandler(c echo.Context) error {
	ctx := c.Request().Context()

	webhookID, err := uuid.Parse(c.Param("webhookID"))
	if err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	if err = h.servs.WebhooksService.DeleteWebhook(ctx, webhookID); err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	return respondhtmx.Refresh(c)
}

func deleteWebhookButton(webhookID uuid.UUID) nodx.Node {
	return component.OptionsDropdownButton(
		htmx.HxDelete(pathutil.BuildPath("/dashboard/webhooks/"+webhookID.String())),
		htmx.HxConfirm("Are you sure you want to delete this webhook?"),
		lucide.Trash(),
		component.SpanText("Delete webhook"),
	)
}
