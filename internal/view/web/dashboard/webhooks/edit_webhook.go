package webhooks

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/eduardolat/pgbackweb/internal/database/dbgen"
	"github.com/eduardolat/pgbackweb/internal/util/echoutil"
	"github.com/eduardolat/pgbackweb/internal/util/pathutil"
	"github.com/eduardolat/pgbackweb/internal/validate"
	"github.com/eduardolat/pgbackweb/internal/view/web/component"
	"github.com/eduardolat/pgbackweb/internal/view/web/respondhtmx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	nodx "github.com/nodxdev/nodxgo"
	htmx "github.com/nodxdev/nodxgo-htmx"
	lucide "github.com/nodxdev/nodxgo-lucide"
)

type editWebhookDTO struct {
	Name      string      `form:"name" validate:"required"`
	EventType string      `form:"event_type" validate:"required"`
	TargetIds []uuid.UUID `form:"target_ids" validate:"required,gt=0"`
	IsActive  string      `form:"is_active" validate:"required,oneof=true false"`
	Url       string      `form:"url" validate:"required,url"`
	Method    string      `form:"method" validate:"required,oneof=GET POST"`
	Headers   string      `form:"headers" validate:"omitempty,json"`
	Body      string      `form:"body" validate:"omitempty,json"`
}

func (h *handlers) editWebhookHandler(c echo.Context) error {
	ctx := c.Request().Context()
	webhookID, err := uuid.Parse(c.Param("webhookID"))
	if err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	var formData editWebhookDTO
	if err := c.Bind(&formData); err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}
	if err := validate.Struct(&formData); err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	_, err = h.servs.WebhooksService.UpdateWebhook(
		ctx, dbgen.WebhooksServiceUpdateWebhookParams{
			WebhookID: webhookID,
			Name:      sql.NullString{String: formData.Name, Valid: true},
			EventType: sql.NullString{String: formData.EventType, Valid: true},
			TargetIds: formData.TargetIds,
			IsActive:  sql.NullBool{Bool: formData.IsActive == "true", Valid: true},
			Url:       sql.NullString{String: formData.Url, Valid: true},
			Method:    sql.NullString{String: formData.Method, Valid: true},
			Headers:   sql.NullString{String: formData.Headers, Valid: true},
			Body:      sql.NullString{String: formData.Body, Valid: true},
		},
	)
	if err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	return respondhtmx.AlertWithRefresh(c, "Webhook updated")
}

func (h *handlers) editWebhookFormHandler(c echo.Context) error {
	ctx := c.Request().Context()
	webhookID, err := uuid.Parse(c.Param("webhookID"))
	if err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	webhook, err := h.servs.WebhooksService.GetWebhook(ctx, webhookID)
	if err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	databases, err := h.servs.DatabasesService.GetAllDatabases(ctx)
	if err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	destinations, err := h.servs.DestinationsService.GetAllDestinations(ctx)
	if err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	backups, err := h.servs.BackupsService.GetAllBackups(ctx)
	if err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	return echoutil.RenderNodx(c, http.StatusOK, editWebhookForm(
		webhook, databases, destinations, backups,
	))
}

func editWebhookForm(
	webhook dbgen.Webhook,
	databases []dbgen.DatabasesServiceGetAllDatabasesRow,
	destinations []dbgen.DestinationsServiceGetAllDestinationsRow,
	backups []dbgen.Backup,
) nodx.Node {
	return nodx.FormEl(
		htmx.HxPost(pathutil.BuildPath(fmt.Sprintf("/dashboard/webhooks/%s/edit", webhook.ID))),
		htmx.HxDisabledELT("find button[type='submit']"),
		nodx.Class("space-y-2"),

		createAndUpdateWebhookForm(databases, destinations, backups, webhook),

		nodx.Div(
			nodx.Class("flex justify-end items-center space-x-2 pt-2"),
			component.HxLoadingMd(),
			nodx.Button(
				nodx.Class("btn btn-primary"),
				nodx.Type("submit"),
				component.SpanText("Save"),
				lucide.Save(),
			),
		),
	)
}

func editWebhookButton(webhookID uuid.UUID) nodx.Node {
	mo := component.Modal(component.ModalParams{
		Size:  component.SizeLg,
		Title: "Edit webhook",
		Content: []nodx.Node{
			nodx.Div(
				htmx.HxGet(pathutil.BuildPath(fmt.Sprintf("/dashboard/webhooks/%s/edit", webhookID))),
				htmx.HxSwap("outerHTML"),
				htmx.HxTrigger("intersect once"),
				nodx.Class("p-10 flex justify-center"),
				component.HxLoadingMd(),
			),
		},
	})

	return nodx.Div(
		mo.HTML,
		component.OptionsDropdownButton(
			mo.OpenerAttr,
			lucide.Pencil(),
			component.SpanText("Edit webhook"),
		),
	)
}
