package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/crypto/ssh"
)

func main() {
	fmt.Println("--- GOSH System Diognostic ---")

	fmt.Printf("[OS]		%s (%s)\n", runtime.GOOS, runtime.GOARCH)

	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Println("[ERR] FFmpeg not found in PATH")
	} else {
		out, _ := exec.Command(path, "--version").Output()
		version := strings.Split(string(out), "\n")[0]
		fmt.Printf("[OK] FFmpeg found: %s\n", path)
		fmt.Printf("		 %s\n", version)
	}

	checkAudio()

	checkSSH("127.0.0.1:22")

}

func checkAudio() {
	fmt.Println("[AUD]  Probing audio input devices...")
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("arecord", "-l")
	case "darwin":
		cmd = exec.Command("ffmpeg", "-f", "avfoundation", "-list_devices", "true", "-i", "")
	case "windows":
		cmd = exec.Command("ffmpeg", "-list_devices", "true", "-f", "dshow", "-i", "dummy")
	}

	if cmd != nil {
		out, _ := cmd.CombinedOutput()
		// We print the raw output so you can find your device index/name
		fmt.Println(string(out))
	}
}

func checkSSH(addr string) {
	fmt.Printf("[SSH]  Probing %s...\n", addr)
	// Just a simple TCP dial to see if the port is open before handshaking
	config := &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// We don't even need auth for a simple handshake test
	_, err := ssh.Dial("tcp", addr, config)

	if err != nil && strings.Contains(err.Error(), "handshake failed") {
		fmt.Printf("[OK]   SSH Port is open (Server responded, Handshake error expected due to no auth)\n")
	} else if err != nil {
		fmt.Printf("[ERR]  SSH Port unreachable: %v\n", err)
	}
}
