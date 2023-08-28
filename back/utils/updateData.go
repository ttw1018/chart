package utils

func UpdateData() error {
	localDir := "./data"
	remoteDir := "/public/home/wlchen/twtang/DST/result"

	return sync(localDir, remoteDir, "192.168.126.81", "wlchen", "22")
}
