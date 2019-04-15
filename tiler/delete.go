package tiler

import (
	"context"

	"github.com/starkandwayne/om-tiler/steps"
)

func (t *Tiler) Delete(ctx context.Context) error {
	return steps.Run(ctx, []steps.Step{
		t.stepPollTillOnline(),
		t.stepConfigureAuthentication(),
		t.stepDeleteInstallation(),
	})
}
