package rtc

type KeyFrameRequestManager struct {
	keyFrameRequestDelay uint32
}

func NewKeyFrameRequestManager(keyFrameRequestDelay uint32) *KeyFrameRequestManager {
	return &KeyFrameRequestManager{
		keyFrameRequestDelay: keyFrameRequestDelay,
	}
}
