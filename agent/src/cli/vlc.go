package cli

import "os/exec"

type VLCRunnerArguments struct {
	VLCPath      string
	MediaFile    string
	ExtraIntf    string
	HttpPort     string
	HttpPassword string
}

func PrepareRunnerArguments(VLCPath, MediaFile, ExtraIntf, HttpPort, HttpPassword string) VLCRunnerArguments {
	return VLCRunnerArguments{
		VLCPath:      VLCPath,
		MediaFile:    MediaFile,
		ExtraIntf:    ExtraIntf,
		HttpPort:     HttpPort,
		HttpPassword: HttpPassword,
	}
}

func PrepareVLCCommand(args VLCRunnerArguments) *exec.Cmd {
	cmd := exec.Command(
		args.VLCPath,
		args.MediaFile,
		"--extraintf", args.ExtraIntf,
		"--http-port", args.HttpPort,
		"--http-password", args.HttpPassword,
	)
	return cmd
}
