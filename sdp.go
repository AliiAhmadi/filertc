package main

import "strings"

func stripSDP(sdp string) string {
	sdp = strings.Replace(sdp, "a=group:BUNDLE audio video data", "a=group:BUNDLE data", -1)
	tmp := strings.Split(sdp, "m=audio")

	bsdp := tmp[0]
	var esdp string

	if len(tmp) > 1 {
		tmp = strings.Split(tmp[1], "a=end-of-candidates")
		esdp = strings.Join(tmp[2:], "a=end-of-candidates")
	} else {
		esdp = strings.Join(tmp[1:], "a=end-of-candidates")
	}

	sdp = bsdp + esdp
	sdp = strings.Replace(sdp, "\r\n\r\n", "\r\n", -1)
	sdp = strings.Replace(sdp, "\n\n", "\n", -1)
	return sdp
}
