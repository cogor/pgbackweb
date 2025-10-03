package backups

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

func (h *handlers) duplicateBackupHandler(c echo.Context) error {
	ctx := c.Request().Context()

	backupID, err := uuid.Parse(c.Param("backupID"))
	if err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	if _, err = h.servs.BackupsService.DuplicateBackup(ctx, backupID); err != nil {
		return respondhtmx.ToastError(c, err.Error())
	}

	return respondhtmx.Refresh(c)
}

func duplicateBackupButton(backupID uuid.UUID) nodx.Node {
	return component.OptionsDropdownButton(
		htmx.HxPost(pathutil.BuildPath("/dashboard/backups/"+backupID.String()+"/duplicate")),
		htmx.HxConfirm("Are you sure you want to duplicate this backup task?"),
		lucide.CopyPlus(),
		component.SpanText("Duplicate backup task"),
	)
}
