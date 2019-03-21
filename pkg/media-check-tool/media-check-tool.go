package media_check_tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/satori/go.uuid"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

type MediaError struct {
	Error  map[string]interface{}  `json:"error"`
}

type TranscodedMedia struct {
	FileName string `json:"file_name"`
}

type MediaInfoJson struct {
	Media  map[string]interface{}  `json:"media"`
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
	ExtraData     string            `json:"extradata"`
}


func extractError(e error, s string) error {
	return e
}

// FFProbe runs default FFprobe (OTF)
func FFProbe(url string)(*Media, *MediaError){
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

	info, err := ExecuteFFProbeCommand(cmd)
	return info, err
}

// FFProbe runs FFprobe with custom parameters (OTF)
func FFProbeCustom(url string, custom_parameters []string)(*Media, *MediaError){

	// Had to stringify the Array of custom parameters so I can append url to the end of the []string
	custom_parameters_w_url_string := strings.Join(custom_parameters, ",") + fmt.Sprintf(",%v", url)
	custom_parameters_w_url := strings.Split(custom_parameters_w_url_string, ",")

	cmd := exec.Command(
		"ffprobe",
		custom_parameters_w_url...,
	)
	info, err := ExecuteFFProbeCommand(cmd)
	return info, err
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
	webm_post_arguments := fmt.Sprintf("-loop 1 -i /pkg/media-check-tool/AudioOnly.png -max_muxing_queue_size 50000 -t %f -y -f webm -acodec libvorbis -ac 2 -sws_flags lanczos -vcodec libvpx -cpu-used 3 -vb 256k -vf \"scale=min(920\\,iw):trunc(ow/a/2)*2,scale=trunc(oh*a/2)*2:min(524\\,ih)\" -r 1 -v error %v.webm",  duration, tmp_file_name)
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
	mp4_post_arguments := fmt.Sprintf("-loop 1 -i /pkg/media-check-tool/AudioOnly.png -max_muxing_queue_size 50000 -t %f -y -f mp4 -acodec aac -ab 96k -ac 2 -sws_flags lanczos -vcodec libx264 -pix_fmt yuv420p -vb 256k -vf \"scale=min(920\\,iw):trunc(ow/a/2)*2,scale=trunc(oh*a/2)*2:min(524\\,ih)\" -r 1 -v error %v.mp4", duration, tmp_file_name)

	mp4_cmd := exec.Command(
		"ffmpeg",
		"-sn",
		"-i", url,
		mp4_post_arguments,
	)
	return ExecuteFFMpegCommand(mp4_cmd, fmt.Sprint("%v.mp4", tmp_file_name))
}

// FFMpeg Snippetizing OTF
func FFMpegSnippetize(url string){
	// Need to pull useful data:
	// audio only
	// how many slices
	// how long slices are
	// whether it should be input or determined
}


// FFMpeg Segmentize OTF
func FFMpegSegmentize(){}

// MediaInfo runs media through the MediaInfo tool
func MediaInfo(url string)(*MediaInfoJson, error){
	cmd := exec.Command(
		"mediainfo",
		"--Output=JSON",
		url,
	)

	data, err := ExecuteMediaInfoCommand(cmd)
	return data, err
}

// CallCheckMedia makes a call to the check media service
func CallCheckMedia(url string, job_id string)(string, error){

	checkMediaUrl := "http://check-media-service.default.svc.cluster.local"
	requestBody := map[string]interface{}{
		"url": url,
		"job_id": job_id,
		"created": fmt.Sprintf(time.Now().Format("2006-01-02 12:01:01")),
	}
	bytesRequestBody, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}

	req, err := http.Post(checkMediaUrl, "application/json", bytes.NewBuffer(bytesRequestBody))

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", req.Status), nil
}

// CallBlankAudio makes a call to the blank audio service
func CallBlankAudio(url string, job_id string, maximum_silence_percent float64, silence_threshold_db float64)(string, error){

	blankAudioUrl := "http://blank-audio-service.cielo24.co"
	requestBody := map[string]interface{}{
		"url": url, "job_id": job_id,
		"max_silence_percent": maximum_silence_percent,
		"silence_threshhold": silence_threshold_db,
		"created": fmt.Sprintf(time.Now().Format("2006-01-02 12:01:01")),
	}
	bytesRequestBody, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}

	req, err := http.Post(blankAudioUrl, "application/json", bytes.NewBuffer(bytesRequestBody))

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", req.Status), nil

}

// RunServer starts the local check media suite (cmd/main.go should be where the server is
//func RunServer(){}

func TranscodeIndividualSlice(url string, slice_start int, slice_length int, audio_only bool)(*FFMpegResults) {
	var pre_seek int
	var accurate_seek int
	var webm_cmd_args []string
	var mp4_cmd_args []string

	tmp_file_name := uuid.Must(uuid.NewV4(), nil)

	if slice_start > 30000 {
		pre_seek = slice_start - 30000 // slice_start is in milliseconds
		accurate_seek = 30000
	}else{
		pre_seek = 0
		accurate_seek = slice_start
	}

	if audio_only {
		webm_cmd_args = []string{
			"-ss", fmt.Sprintf("%d", pre_seek),
			"-i", url,
			"-sn",
			"-loop", "1",
			"-i", "/pkg/media-check-tool/AudioOnly.png",
			"-max_muxing_queue_size", "50000",
			"-ss", fmt.Sprintf("%d", accurate_seek),
			"-t", fmt.Sprintf("%d", slice_length),
			"-y",
			"-f", "webm",
			"-acodec", "libvorbis",
			"-ac", "2",
			"-sws_flags", "lanczos",
			"-vcodec", "libvpx",
			"-cpu-used", "3",
			"-vb", "256k",
			"-vf", "\"scale=min(920\\,iw):trunc(ow/a/2)*2,scale=trunc(oh*a/2)*2:min(524\\,ih)\"",
			"-r", "24",
		}
		mp4_cmd_args = []string{
			"-ss", fmt.Sprintf("%d", pre_seek),
			"-i", url,
			"-sn",
			"-loop", "1",
			"-i", "/pkg/media-check-tool/AudioOnly.png",
			"-max_muxing_queue_size", "50000",
			"-ss", fmt.Sprintf("%d", accurate_seek),
			"-t", fmt.Sprintf("%d", slice_length),
			"-y",
			"-f", "mp4",
			"-acodec", "aac",
			"-ab", "96k",
			"-ac", "2",
			"-sws_flags", "lanczos",
			"-vcodec", "libx264",
			"-pix_fmt", "yuv420p",
			"-vb", "256k",
			"-vf", "\"scale=min(920\\,iw):trunc(ow/a/2)*2,scale=trunc(oh*a/2)*2:min(524\\,ih)\"",
			"-r", "24",
		}
	}else{
		webm_cmd_args = []string{
			"-ss", fmt.Sprintf("%d", pre_seek),
			"-sn",
			"-i", url,
			"-max_muxing_queue_size", "50000",
			"-ss", fmt.Sprintf("%d", accurate_seek),
			"-t", fmt.Sprintf("%d", slice_length),
			"-y",
			"-f", "webm",
			"-acodec", "libvorbis",
			"-ac", "2",
			"-sws_flags", "lanczos",
			"-vcodec", "libvpx",
			"-cpu-used", "3",
			"-vb", "256k",
			"-vf", "\"scale=min(920\\,iw):trunc(ow/a/2)*2,scale=trunc(oh*a/2)*2:min(524\\,ih)\"",
			"-r", "24",
		}
		mp4_cmd_args = []string{
			"-ss", fmt.Sprintf("%d", pre_seek),
			"-sn",
			"-i", url,
			"-max_muxing_queue_size", "50000",
			"-ss", fmt.Sprintf("%d", accurate_seek),
			"-t", fmt.Sprintf("%d", slice_length),
			"-y",
			"-f", "mp4",
			"-acodec", "aac",
			"-ab", "96k",
			"-ac", "2",
			"-sws_flags", "lanczos",
			"-vcodec", "libx264",
			"-pix_fmt", "yuv420p",
			"-vb", "256k",
			"-vf", "\"scale=min(920\\,iw):trunc(ow/a/2)*2,scale=trunc(oh*a/2)*2:min(524\\,ih)\"",
			"-r", "24",
		}
	}

	webm_cmd := exec.Command(
		"ffmpeg",
		webm_cmd_args...,
	)

	mp4_cmd := exec.Command(
		"ffmpeg",
		mp4_cmd_args...,
	)

	webm_results, webm_error := ExecuteFFMpegCommand(webm_cmd, fmt.Sprintf("%v.webm", tmp_file_name))
	mp4_results, mp4_error := ExecuteFFMpegCommand(mp4_cmd, fmt.Sprintf("%v.mp4", tmp_file_name))

	results := &FFMpegResults{}
	results.FileNameWebm = webm_results.FileName
	results.ErrorWebm = webm_error
	results.FileNameMp4 = mp4_results.FileName
	results.ErrorMP4 = mp4_error

	return results
}

// Function used to run and return all ffprobe commands (genericized what can be generic)
func ExecuteFFProbeCommand(cmd *exec.Cmd)(*Media, *MediaError){
	errors := &MediaError{}

	var cmdOut, cmdErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &cmdOut, &cmdErr

	if err := cmd.Run(); err != nil {
		StdOut := cmdOut.Bytes()
		mapped := make(map[string]interface{})
		_ = json.Unmarshal(StdOut, &mapped)
		_ = mapstructure.Decode(mapped, errors)
		return nil, errors
	}

	info := &Media{}
	if err := json.Unmarshal([]byte(cmdOut.String()), &info); err != nil {
		return nil, nil
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

// Function used to run and return all ffmpeg commands (genericized what can be generic)
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

// Function used to run and return all media-info commands (genericized what can be generic)
func ExecuteMediaInfoCommand(cmd *exec.Cmd)(*MediaInfoJson, error){
	var cmdOut, cmdErr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &cmdOut, &cmdErr

	if err := cmd.Run(); err != nil{
		// TODO: Error Logic
		return nil, err
	}
	data := &MediaInfoJson{}
	if err := json.Unmarshal(cmdOut.Bytes(), &data); err != nil{
		// TODO: Error Logic
		return nil, err
	}
	return data, nil
}