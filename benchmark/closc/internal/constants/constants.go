package constants

const (
	SubmitTXBatchSize         = 5000
	SubmitTXMaxRetries        = 5
	SubmitTXRetryDelaySeconds = 2

	LocalChainTXIDLogPath     = "lc_tx_id.log"
	LocalChainCommitTXLogPath = "lc_commit_tx.log"
	OrgChainTXIDLogPath       = "oc_tx_id.log"
	AudChainTXIDLogPath       = "ac_tx_id.log"
)
