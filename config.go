package godueros

const (
	Sample_Rate 		int   = 16000
	// 填满一个buffer调用一次回调
	Frame_Per_Buffer 	int   = 256
	Channels_Num		int32 = 1
	// 10ms per chunk
	Chunk_Per_Sec       int   = 100
	// record voice 5m timeout
	Timeout_Chunk_Num   int   = 500
)

var (
	ACCESS_TOKEN = "26.687423ae6cd02c090ef92af1b64be322.2592000.1514813398.2300547068-10218111"
	DEVICE_ID 	 = "wh9foSD4KnZAa8kyG1l62SdwkNHlVfuH"
	REPLY_FILE 	 = "dureply.mp3"
)

