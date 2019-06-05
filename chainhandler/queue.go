package chainhandler

type blockQueueItem struct {
	block   Block
	actions []Action
}
