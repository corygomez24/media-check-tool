package media_processing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"github.com/satori/go.uuid"
)


// Media describes file information.
type Media struct {
	// Format common file format information.
	Format map[string]interface{} `json:"format"`

	// Streams item information structure.
	Streams []*Stream `json:"streams"`

	// Duration is duration of storage file in seconds.
	Duration float64

	// Size of file in bytes.
	Size uint64
}

type TranscodedMedia struct {
	FileName string `json:"file_name"`
}

type FFMpegResults struct {
	FileNameMP3   string  `json:"file_name_mp3"`
	ErrorMP3      error  `json:"error_mp3"`
	FileNameWebm  string  `json:"file_name_webm"`
	ErrorWebm     error  `json:"error_webm"`
	FileNameMp4   string  `json:"file_name_mp4"`
	ErrorMP4      error  `json:"error_mp4"`

}

func (info *Media) Read(p []byte) (n int, err error) {
	panic("implement me")
}

// Strings packs info into json.
func (info *Media) String() string {
	body, _ := json.MarshalIndent(info, "", "  ")
	return string(body)
}

// Stream information about single stream.
type Stream struct {
	Index         int               `json:"index"`
	Type          string            `json:"codec_type"`
	CodecName     string            `json:"codec_name"`
	CodecLongName string            `json:"codec_long_name"`
	Bitrate       string            `json:"bitRate"`
	Duration      float64           `json:"duration,string"`
	Width         int64             `json:"width"`
	Height        int64             `json:"height"`
	Tags          map[string]string `json:"tags,omitempty"`
}


func extractError(e error, s string) error {
	return e
}

// FFProbe runs default FFprobe (OTF)
func FFProbe(url string)(*Media, error){
	cmd := exec.Command(
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-show_error",
		"-show_chapters",
		url,
	)

	info, nil := ExecuteFFProbeCommand(cmd)
	return info, nil
}

// FFProbe runs FFprobe with custom parameters (OTF)
func FFProbeCustom(url string, custom_parameters string)(*Media, error){

	cmd := exec.Command(
		"ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-show_error",
		"-show_chapters",
		custom_parameters,
		url,
	)
	info, nil := ExecuteFFProbeCommand(cmd)
	return info, nil
}

// FFMpegAudioOnly (FFMpeg Audio Only OTF)
func FFMpegAudioOnly(url string, duration float64)(*FFMpegResults){

	tmp_file_name := uuid.Must(uuid.NewV4(), nil)

	file_name_mp3, err_mp3 :=TranscodeToMP3(url, tmp_file_name)
	file_name_webm, err_webm := TranscodeToWebmAudioOnly(url, tmp_file_name, duration)
	file_name_mp4, err_mp4 := TranscodeToMP4AudioOnly(url, tmp_file_name, duration)

	// Command to run MP4Box (Not sure if using yet)
	//mp4box_arguments := "-hint -tmp /tmp/ -inter 0"
	//mp4box_cmd := exec.Command("MP4Box", mp4box_arguments)

	results := &FFMpegResults{}
	results.FileNameMP3 = file_name_mp3.FileName
	results.ErrorMP3 = err_mp3
	results.FileNameMp4 = file_name_mp4.FileName
	results.ErrorMP4 = err_mp4
	results.FileNameWebm = file_name_webm.FileName
	results.ErrorWebm = err_webm

	return results
}

//Transcodes a video url into mp3 format, returns file name or error
func TranscodeToMP3(url string, tmp_file_name uuid.UUID)(*TranscodedMedia, error){

	// Command to Transcode into mp3
	mp3_post_arguments := fmt.Sprintf("-max_muxing_queue_size 50000 -y -f mp3 -ab 192k -ac 2 -sws_flags lanczos -v error %v.mp3", tmp_file_name)

	mp3_cmd := exec.Command(
		"ffmpeg",
		"-sn",
		"-i", url,
		mp3_post_arguments,
	)
	return ExecuteFFMpegCommand(mp3_cmd, fmt.Sprintf("%v.mp3", tmp_file_name))
}

//Transcodes a video url into webm format, returns file name or error
func TranscodeToWebmAudioOnly(url string, tmp_file_name uuid.UUID, duration float64)(*TranscodedMedia, error){

	// Command to Transcode into webm (audio only with splash_screen)
	webm_post_arguments := fmt.Sprintf("-loop 1 -i /app/assets/AudioOnly.png -max_muxing_queue_size 50000 -t %f -y -f webm -acodec libvorbis -ac 2 -sws_flags lanczos -vcodec libvpx -cpu-used 3 -vb 256k -vf \"scale=min(920\\,iw):trunc(ow/a/2)*2,scale=trunc(oh*a/2)*2:min(524\\,ih)\" -r 1 -v error %v.webm",  duration, tmp_file_name)
	webm_cmd := exec.Command(
		"ffmpeg",
		"-sn",
		"-i", url,
		webm_post_arguments,
	)
	return ExecuteFFMpegCommand(webm_cmd, fmt.Sprintf("%v.webm", tmp_file_name))
}

//Transcodes a video url into mp4 format, returns file name or error
func TranscodeToMP4AudioOnly(url string, tmp_file_name uuid.UUID, duration float64)(*TranscodedMedia, error){

	// Command to Transcode into mp4 (audio only with splash_screen)
	mp4_post_arguments := fmt.Sprintf("-loop 1 -i /app/assets/AudioOnly.png -max_muxing_queue_size 50000 -t %f -y -f mp4 -acodec aac -ab 96k -ac 2 -sws_flags lanczos -vcodec libx264 -pix_fmt yuv420p -vb 256k -vf \"scale=min(920\\,iw):trunc(ow/a/2)*2,scale=trunc(oh*a/2)*2:min(524\\,ih)\" -r 1 -v error %v.mp4", duration, tmp_file_name)

	mp4_cmd := exec.Command(
		"ffmpeg",
		"-sn",
		"-i", url,
		mp4_post_arguments,
	)
	return ExecuteFFMpegCommand(mp4_cmd, fmt.Sprint("%v.mp4", tmp_file_name))
}

// FFMpeg Snippetizing OTF
func FFMpegSnippetize(){}

// FFMpeg Segmentize OTF
func FFMpegSegmentize(){}

// MediaInfo runs media through the MediaInfo tool
func MediaInfo(){}

// CallCheckMedia makes a call to the check media service
func CallCheckMedia(){}

// CallBlankAudio makes a call to the blank audio service
func CallBlankAudio(){}

// RunServer starts the local check media suite
func RunServer(){}

func ExecuteFFProbeCommand(cmd *exec.Cmd)(*Media, error){
	var cmdOut, cmdErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &cmdOut, &cmdErr

	if err := cmd.Run(); err != nil {
		return nil, extractError(err, cmdErr.String())
	}

	info := &Media{}
	if err := json.Unmarshal([]byte(cmdOut.String()), &info); err != nil {
		return nil, err
	}

	// post processing
	streams := make([]*Stream, 0)
	for _, stream := range info.Streams {
		if stream.Type == "video" {
			if stream.Width == 0 || stream.Height == 0 {
				continue
			}
		}

		if info.Duration < stream.Duration {
			info.Duration = stream.Duration
		}

		streams = append(streams, stream)
	}

	info.Streams = streams

	// fallback duration using container information
	if info.Duration == 0 {
		if _, ok := info.Format["duration"]; ok {
			if d, err := strconv.ParseFloat(info.Format["duration"].(string), 64); err != nil {
				info.Duration = float64(d)
			}
		}
	}

	for _, stream := range info.Streams {
		if stream.Duration == 0 {
			stream.Duration = info.Duration
		}

		if stream.Bitrate == "" {
			stream.Bitrate = "0"
		}

		// flipping width and height
		if angle, rotate := stream.Tags["rotate"]; rotate {
			if angle == "90" || angle == "-90" {
				// please do not rotate by more than 180
				stream.Width, stream.Height = stream.Height, stream.Width
			}
		}
	}

	return info, nil
}

func ExecuteFFMpegCommand(cmd *exec.Cmd, FileName string)(*TranscodedMedia, error){
	var cmdOut, cmdErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &cmdOut, &cmdErr

	if err := cmd.Run(); err != nil {
		return nil, extractError(err, cmdErr.String())
	}

	results := &TranscodedMedia{}
	results.FileName = FileName
	return results, nil
}
