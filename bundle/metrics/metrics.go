package metrics

const (
	ExperimentalPythonWheelWrapperIsSet = "python_wheel_wrapper_is_set"
	ArtifactDynamicVersionIsSet         = "artifact_dynamic_version_is_set"
	ArtifactBuildCommandIsSet           = "artifact_build_command_is_set"
	ArtifactFilesIsSet                  = "artifact_files_is_set"
	PresetsNamePrefixIsSet              = "presets_name_prefix_is_set"
	AppLifecycleStarted                 = "app_lifecycle_started"
	ClusterLifecycleStarted             = "cluster_lifecycle_started"
	SqlWarehouseLifecycleStarted        = "sql_warehouse_lifecycle_started"
	SelectUsed                          = "select_used"

	// Whether workspace.state_path is under /Workspace/Shared.
	StatePathIsShared = "state_path_is_shared"

	// Whether this deploy is compatible with an automatic DMS migration. Migration
	// moves the deployment state behind a deployment object that carries only the
	// bundle's statically declared permissions, so it is seamless only when every
	// CAN_MANAGE (or IS_OWNER) permission on the state folder is statically declared
	// in the permissions section. A state folder in /Workspace/Shared is manageable by
	// all workspace users and so requires group_name: users CAN_MANAGE; a folder under
	// /Workspace/Users/<owner> always carries the owner's CAN_MANAGE.
	//
	// Exactly one of the three keys below is recorded per deploy:
	//   - definitely: every manager of the state folder is statically declared.
	//   - not:        the state folder has a manager that is not statically declared.
	//   - maybe:      no permissions section is set and the state folder is not in a
	//                 location with known managers, so we cannot tell without an
	//                 additional GetPermissions call.
	//
	// A deployment (a bundle target, identified by deployment ID) is auto-migratable
	// only if all of its deploys are compatible, so aggregate these by deployment ID.
	IsDefinitelyAutoMigrationCompatible = "is_definitely_auto_migration_compatible"
	IsMaybeAutoMigrationCompatible      = "is_maybe_auto_migration_compatible"
	IsNotAutoMigrationCompatible        = "is_not_auto_migration_compatible"
)
