package mutator

import (
	"context"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/config/resources"
	"github.com/databricks/cli/libs/diag"
	"github.com/databricks/cli/libs/dyn"
	"github.com/databricks/cli/libs/dyn/convert"
)

type initializeVolumePaths struct{}

// InitializeVolumePaths sets resources.volumes.*.volume_path from catalog, schema, and name.
// The path is only set when those fields are present and fully resolved (no ${...} references).
// This enables ${resources.volumes.<key>.volume_path} interpolation during initialize.
func InitializeVolumePaths() bundle.Mutator {
	return &initializeVolumePaths{}
}

func (m *initializeVolumePaths) Name() string {
	return "InitializeVolumePaths"
}

func (m *initializeVolumePaths) Apply(_ context.Context, b *bundle.Bundle) diag.Diagnostics {
	err := b.Config.Mutate(func(root dyn.Value) (dyn.Value, error) {
		pattern := dyn.NewPattern(dyn.Key("resources"), dyn.Key("volumes"), dyn.AnyKey())
		return dyn.MapByPattern(root, pattern, func(_ dyn.Path, v dyn.Value) (dyn.Value, error) {
			var vol resources.Volume
			if err := convert.ToTyped(&vol, v); err != nil {
				return dyn.InvalidValue, err
			}
			return dyn.Set(v, "volume_path", dyn.V(vol.ComputeVolumePath()))
		})
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
