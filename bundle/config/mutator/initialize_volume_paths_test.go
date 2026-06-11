package mutator

import (
	"testing"

	"github.com/databricks/cli/bundle"
	"github.com/databricks/cli/bundle/config"
	"github.com/databricks/cli/bundle/config/resources"
	"github.com/databricks/cli/libs/dyn"
	"github.com/databricks/databricks-sdk-go/service/catalog"
	"github.com/stretchr/testify/require"
)

func TestInitializeVolumePaths(t *testing.T) {
	b := &bundle.Bundle{
		Config: config.Root{
			Resources: config.Resources{
				Volumes: map[string]*resources.Volume{
					"bar": {
						CreateVolumeRequestContent: catalog.CreateVolumeRequestContent{
							CatalogName: "main",
							SchemaName:  "myschema",
							Name:        "volbar",
						},
					},
					"foo": {
						CreateVolumeRequestContent: catalog.CreateVolumeRequestContent{
							CatalogName: "main",
							SchemaName:  "${resources.schemas.my.name}",
							Name:        "volfoo",
						},
					},
				},
			},
		},
	}

	diags := bundle.Apply(t.Context(), b, InitializeVolumePaths())
	require.NoError(t, diags.Error())
	require.Equal(t, "/Volumes/main/myschema/volbar", b.Config.Resources.Volumes["bar"].VolumePath)
	require.Empty(t, b.Config.Resources.Volumes["foo"].VolumePath)
}

// TestVolumePathPipeline_ResolveThenCompute mirrors applyVolumePathMutators.
func TestVolumePathPipeline_ResolveThenCompute(t *testing.T) {
	b := &bundle.Bundle{
		Config: config.Root{
			Resources: config.Resources{
				Schemas: map[string]*resources.Schema{
					"my": {
						CreateSchema: catalog.CreateSchema{
							CatalogName: "main",
							Name:        "myschema",
						},
					},
				},
				Volumes: map[string]*resources.Volume{
					"bar": {
						CreateVolumeRequestContent: catalog.CreateVolumeRequestContent{
							CatalogName: "main",
							SchemaName:  "${resources.schemas.my.name}",
							Name:        "volbar",
						},
					},
				},
			},
		},
	}

	diags := bundle.ApplySeq(
		t.Context(),
		b,
		ResolveVariableReferencesByPattern(
			dyn.NewPattern(dyn.Key("resources"), dyn.Key("volumes"), dyn.AnyKey(), dyn.Key("catalog_name")),
			"resources",
		),
		ResolveVariableReferencesByPattern(
			dyn.NewPattern(dyn.Key("resources"), dyn.Key("volumes"), dyn.AnyKey(), dyn.Key("schema_name")),
			"resources",
		),
		ResolveVariableReferencesByPattern(
			dyn.NewPattern(dyn.Key("resources"), dyn.Key("volumes"), dyn.AnyKey(), dyn.Key("name")),
			"resources",
		),
		InitializeVolumePaths(),
	)
	require.NoError(t, diags.Error())
	require.Equal(t, "myschema", b.Config.Resources.Volumes["bar"].SchemaName)
	require.Equal(t, "/Volumes/main/myschema/volbar", b.Config.Resources.Volumes["bar"].VolumePath)
}

func TestVolumePathPipeline_ResolvesCrossVolumeReference(t *testing.T) {
	b := &bundle.Bundle{
		Config: config.Root{
			Resources: config.Resources{
				Volumes: map[string]*resources.Volume{
					"bar": {
						CreateVolumeRequestContent: catalog.CreateVolumeRequestContent{
							CatalogName: "main",
							SchemaName:  "myschema",
							Name:        "volbar",
						},
					},
					"foo": {
						CreateVolumeRequestContent: catalog.CreateVolumeRequestContent{
							CatalogName: "main",
							SchemaName:  "myschema",
							Name:        "volfoo",
							Comment:     "${resources.volumes.bar.volume_path}",
						},
					},
				},
			},
		},
	}

	diags := bundle.ApplySeq(
		t.Context(),
		b,
		ResolveVariableReferencesByPattern(
			dyn.NewPattern(dyn.Key("resources"), dyn.Key("volumes"), dyn.AnyKey(), dyn.Key("catalog_name")),
			"resources",
		),
		ResolveVariableReferencesByPattern(
			dyn.NewPattern(dyn.Key("resources"), dyn.Key("volumes"), dyn.AnyKey(), dyn.Key("schema_name")),
			"resources",
		),
		ResolveVariableReferencesByPattern(
			dyn.NewPattern(dyn.Key("resources"), dyn.Key("volumes"), dyn.AnyKey(), dyn.Key("name")),
			"resources",
		),
		InitializeVolumePaths(),
		ResolveVariableReferencesOnlyResources("resources"),
	)
	require.NoError(t, diags.Error())
	require.Equal(t, "/Volumes/main/myschema/volbar", b.Config.Resources.Volumes["bar"].VolumePath)
	require.Equal(t, "/Volumes/main/myschema/volbar", b.Config.Resources.Volumes["foo"].Comment)
}
