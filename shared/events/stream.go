package events

type DataProvider string

const (
	DataProviderESPN DataProvider = "ESPN"
)

type StreamMode string

const (
	StreamModeLive   StreamMode = "LIVE"
	StreamModeReplay StreamMode = "REPLAY"
)
