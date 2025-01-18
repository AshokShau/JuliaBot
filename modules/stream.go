package modules

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/amarnathcjd/gogram/telegram"
	tg "github.com/amarnathcjd/gogram/telegram"
)

func StreamHandler(m *tg.NewMessage) error {
	if m.IsPrivate() {
		m.Reply("This command can only be used in groups")
		return nil
	}

	if !m.IsReply() {
		m.Reply("Reply to a video/audio file to stream it")
		return nil
	}

	reply, err := m.GetReplyMessage()
	if err != nil {
		m.Reply("Failed to get reply message")
		return err
	}

	if !reply.IsMedia() {
		m.Reply("The replied message is not a media file")
		return nil
	}

	if reply.Audio() == nil && reply.Video() == nil && reply.Document() == nil {
		m.Reply("The replied message is not a video/audio file")
		return nil
	}

	m.Reply("Starting real-time streaming...")
	Stream(reply)
	return nil
}

func Stream(m *telegram.NewMessage) {
	fmt.Println("Starting real-time streaming...")
	var chunkSize int64 = 1024 * 1024
	fileSize := m.File.Size

	cmd := exec.Command("ffmpeg",
		"-stream_loop", "-1",
		"-re",
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-b:v", "2500k",
		"-maxrate", "2500k",
		"-bufsize", "5000k",
		"-pix_fmt", "yuv420p",
		"-g", "50",
		"-c:a", "aac",
		"-b:a", "128k",
		"-ac", "2",
		"-ar", "44100",
		"-f", "flv",
		os.Getenv("STREAM_URL"),
	)

	ffmpegIn, err := cmd.StdinPipe()
	if err != nil {
		m.Reply("Failed to initialize ffmpeg stdin pipe")
		return
	}

	if err := cmd.Start(); err != nil {
		m.Reply("Failed to start ffmpeg")
		return
	}

	go func() {
		defer ffmpegIn.Close()

		for i := int64(0); i < fileSize; i += chunkSize {
			chunk, _, err := m.Client.DownloadChunk(m.Media(), int(i), int(i+chunkSize), int(chunkSize))
			if err != nil {
				m.Reply("Failed to get file chunk")
				return
			}

			_, writeErr := ffmpegIn.Write(chunk)
			if writeErr != nil {
				fmt.Println("Failed to write chunk to ffmpeg:", writeErr)
				return
			}

		}
	}()

	if err := cmd.Wait(); err != nil {
		fmt.Println("ffmpeg process ended with error:", err)
		m.Reply("Streaming stopped with an error")
	} else {
		m.Reply("Streaming completed successfully")
	}
}
