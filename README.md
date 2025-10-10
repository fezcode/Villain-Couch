# VLC Tracker Agent

This project is a Windows agent to track VLC Media Player status. For more detailed information, please see the README in the `agent` directory.

## Building

This is a Windows application. To build the application, run the `build.ps1` PowerShell script.

```powershell
.\build.ps1
```

## Directory Structure

### agent

This directory contains the main source code for the agent application. It is responsible for all the core logic, including interacting with VLC, managing data, and handling user commands.

### bin

This directory holds the compiled binaries of the application after a successful build using the `build.ps1` script.

### common

This directory contains shared utility packages that are used by different parts of the application. This includes modules for logging, error handling, and other common functionalities.