package ffmpeg

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
)

func CaptureStream() (io.ReadCloser, *exec.Cmd, error) {
	var inputFormat, inputDevice string
	switch runtime.GOOS {
	case "linux":
		inputFormat = "alsa"
		inputDevice = "default"
	case "darwin":
		inputFormat = "avfoundation"
		inputDevice = ":0"
	case "windows":
		inputFormat = "dshow"
		inputDevice = "audio=Microphone"
	default:
		return nil, nil, fmt.Errorf("Unsupported Platform")
	}

	cmd := exec.Command("ffmpeg",
		"-loglevel", "quiet",
		"-f", inputFormat,
		"-i", inputDevice,
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "2",
		"pipe:1",
	)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	return stdOut, cmd, nil

}
