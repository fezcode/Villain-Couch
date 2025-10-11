package storage

import _ "embed"

//go:embed queries/createTables.sql
var queryCreateTables string

//go:embed queries/setMediaFile.sql
var querySetMediaFile string

//go:embed queries/getMediaFile.sql
var queryGetMediaFile string

//go:embed queries/getLatestMediaFile.sql
var queryGetLatestMediaFile string

//go:embed queries/setWorkspace.sql
var querySetWorkspace string

//go:embed queries/getWorkspace.sql
var queryGetWorkspace string
