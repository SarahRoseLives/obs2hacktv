package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	// Prepare ffmpeg command: listen for RTMP, transcode to low-res, low-fps yuv4mpegpipe
	ffmpeg := exec.Command("ffmpeg",
		"-listen", "1",
		"-i", "rtmp://0.0.0.0:1935/live/stream",
		"-vf", "scale=360:240",
		"-r", "15",
		"-pix_fmt", "yuv420p",
		"-an", // no audio
		"-f", "yuv4mpegpipe", "-",
	)
	ffmpeg.Stderr = os.Stderr

	// Prepare hacktv command: read video from stdin, output NTSC analog signal
	hacktv := exec.Command("./hacktv",
		"-m", "m",
		"-f", "471250000",
		"-s", "13500000",
		"-g", "40",
		"--noaudio",
		"-", // read input from stdin
	)
	hacktv.Stdout = os.Stdout
	hacktv.Stderr = os.Stderr

	// Pipe ffmpeg's stdout to hacktv's stdin
	ffmpegOut, err := ffmpeg.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to set up ffmpeg stdout pipe: %v", err)
	}
	hacktv.Stdin = ffmpegOut

	log.Println("Waiting for RTMP stream at rtmp://localhost:1935/live/stream...")

	// Start hacktv (must be running to accept pipe input)
	if err := hacktv.Start(); err != nil {
		log.Fatalf("Failed to start hacktv: %v", err)
	}
	// Start ffmpeg
	if err := ffmpeg.Start(); err != nil {
		log.Fatalf("Failed to start ffmpeg: %v", err)
	}
	// Wait for ffmpeg to finish
	if err := ffmpeg.Wait(); err != nil {
		log.Printf("ffmpeg exited with error: %v", err)
	}
	// Wait for hacktv to finish
	if err := hacktv.Wait(); err != nil {
		log.Printf("hacktv exited with error: %v", err)
	}
}
