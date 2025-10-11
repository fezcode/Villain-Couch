package app_info

// VersionInfo holds the application version string.
// this value should be overwritten during the build process using linker flags (ldflags).
// But not for now
//
// Example build command:
// go build -ldflags="-X 'app_info.VersionInfo=v1.2.3'"
var VersionInfo = "0.0.1"
