package backups

import (
	"fmt"

	"github.com/eduardolat/pgbackweb/internal/util/pathutil"
	"github.com/eduardolat/pgbackweb/internal/view/web/component"
	"github.com/eduardolat/pgbackweb/internal/view/web/respondhtmx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	nodx "github.com/nodxdev/nodxgo"
	htmx "github.com/nodxdev/nodxgo-htmx"
	lucide "github.com/nodxdev/nodxgo-lucide"
)

func (h *handlers) deleteBackupHandler(c echo.Context) error {
	ctx := c.Request().Context()

	backupID, err := uuid.Parse(c.Param("backupID"))
	if err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	if err = h.servs.BackupsService.DeleteBackup(ctx, backupID); err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	return respondhtmx.Refresh(c)
}

func deleteBackupButton(backupID uuid.UUID) nodx.Node {
	return component.OptionsDropdownButton(
		htmx.HxDelete(pathutil.BuildPath(fmt.Sprintf("/dashboard/backups/%s", backupID))),
		htmx.HxConfirm("Are you sure you want to delete this backup task?"),
		lucide.Trash(),
		component.SpanText("Delete backup task"),
	)
}
