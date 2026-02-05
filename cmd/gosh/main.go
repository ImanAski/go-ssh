package main

import (
	"gosh/internal/ffmpeg"
	"io"
	"log"

	"golang.org/x/crypto/ssh"
)

func main() {
	sshCfg := &ssh.ClientConfig{
		User:            "sophos",
		Auth:            []ssh.AuthMethod{ssh.Password("1")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := "localhost:22"

	client, err := ssh.Dial("tcp", addr, sshCfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	audioStream, cmd, err := ffmpeg.CaptureStream()
	if err != nil {
		panic(err)
	}
	defer audioStream.Close()

	sshStdin, _ := session.StdinPipe()

	session.Start("aplay -f cd")

	go func() {
		io.Copy(sshStdin, audioStream)
		sshStdin.Close()
	}()

	cmd.Wait()
	session.Wait()
}
