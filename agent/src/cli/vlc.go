package cli

import "os/exec"

type VLCRunnerArguments struct {
	VLCPath      string
	MediaFile    string
	StartTime    string
	ExtraIntf    string
	HttpPort     string
	HttpPassword string
}

func PrepareRunnerArguments(VLCPath, MediaFile, StartTime, ExtraIntf, HttpPort, HttpPassword string) VLCRunnerArguments {
	return VLCRunnerArguments{
		VLCPath:      VLCPath,
		MediaFile:    MediaFile,
		StartTime:    StartTime,
		ExtraIntf:    ExtraIntf,
		HttpPort:     HttpPort,
		HttpPassword: HttpPassword,
	}
}

func PrepareVLCCommand(args VLCRunnerArguments) *exec.Cmd {
	arr := []string{
		args.MediaFile,
		"--extraintf", args.ExtraIntf,
		"--http-port", args.HttpPort,
		"--http-password", args.HttpPassword,
		"--start-time", args.StartTime,
	}

	cmd := exec.Command(args.VLCPath, arr...)
	return cmd
}
