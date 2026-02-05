package ffmpeg

import (
	"io"
	"net"
	"os"
	"testing"

	"golang.org/x/crypto/ssh"
)

func TestAudioStreamOverSSH(t *testing.T) {
	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	privateBytes, err := os.ReadFile("/home/sophos/.ssh/id_ed25519")
	if err != nil {
		t.Errorf("Failed to load private: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		t.Errorf("Unable to parst private key: %v", err)
	}

	config.AddHostKey(private)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, _ := listener.Accept()
		_, chans, reqs, _ := ssh.NewServerConn(conn, config)
		go ssh.DiscardRequests(reqs)
		for newChan := range chans {
			channel, requests, _ := newChan.Accept()
			go func(in <-chan *ssh.Request) {
				for req := range in {
					req.Reply(true, nil)
				}
			}(requests)

			buf := make([]byte, 1024)
			n, _ := channel.Read(buf)
			if n > 0 {
				t.Logf("Received")
			}
		}
	}()

	stream, cmd, err := CaptureStream()
	if err != nil {
		t.Skip(err)
	}

	defer cmd.Process.Kill()

	client, err := ssh.Dial("tcp", listener.Addr().String(), &ssh.ClientConfig{
		User: "sophos",
		Auth: []ssh.AuthMethod{
			ssh.Password("1"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		t.Errorf("Failed to make client (on %s): %v", listener.Addr().String(), err)
	}

	session, err := client.NewSession()
	if err != nil {
		t.Errorf("Failed to make a session: %v", err)
	}
	sshStdin, _ := session.StdinPipe()

	session.Shell()

	_, err = io.CopyN(sshStdin, stream, 5000)
	if err != nil {
		t.Errorf("Failed to stream audio: %v", err)
	}

}
