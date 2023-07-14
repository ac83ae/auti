package flag

const (
	PhaseInitialization         = "in"
	PhaseTransactionRecord      = "tr"
	PhaseConsistencyExamination = "ce"
	PhaseResultVerification     = "rv"

	ProcessINDefault = "default"
	ProcessINRandGen = "rand_gen"

	ProcessTRCommitment                  = "commitment"
	ProcessTRMerkleProofGen              = "merkle_proof_gen"
	ProcessTRLocalChainSubmit            = "local_submit"
	ProcessTRLocalChainPrepare           = "local_prepare"
	ProcessTRLocalChainRead              = "local_read"
	ProcessTRLocalChainReadAll           = "local_read_all"
	ProcessTRLocalChainCommitmentSubmit  = "local_commit_submit"
	ProcessTRLocalChainCommitmentPrepare = "local_commit_prepare"
	ProcessTRLocalCHainCommitmentRead    = "local_commit_read"
	ProcessTRLocalChainCommitmentReadAll = "local_commit_read_all"
	ProcessTROrgChainSubmit              = "org_submit"
	ProcessTROrgChainPrepare             = "org_prepare"
	ProcessTROrgChainRead                = "org_read"
	ProcessTROrgChainReadAll             = "org_read_all"
	ProcessTRLocalChainSC                = "local_sc"

	ProcessCEMerkleProofVerify                       = "merkle_proof_verify"
	ProcessCEMerkleProofMerge                        = "merkle_proof_merge"
	ProcessCESummarizeMerkleProofVerificationResults = "summarize_proof_result"
	ProcessCEVerifyCommitments                       = "verify_commit"
	ProcessCEAccumulateCommitments                   = "accumulate_commit"
	ProcessCEAudChainSubmit                          = "aud_submit"
	ProcessCEAudChainPrepare                         = "aud_prepare"
	ProcessCEAudChainRead                            = "aud_read"
	ProcessCEAudChainReadAll                         = "aud_read_all"

	ProcessRVVerifyMerkleBatchProof                       = "verify_batch_proof"
	ProcessRVSummarizeMerkleBatchProofVerificationResults = "summarize_batch_proof"
	ProcessRVVerifyCommitments                            = "verify_commits"
)

var PhaseProcessMap = map[string][]string{
	PhaseInitialization: {
		ProcessINDefault,
		ProcessINRandGen,
	},
	PhaseTransactionRecord: {
		ProcessTRCommitment,
		ProcessTRMerkleProofGen,
		ProcessTRLocalChainSubmit,
		ProcessTRLocalChainPrepare,
		ProcessTRLocalChainRead,
		ProcessTRLocalChainReadAll,
		ProcessTRLocalChainCommitmentSubmit,
		ProcessTRLocalChainCommitmentPrepare,
		ProcessTRLocalCHainCommitmentRead,
		ProcessTRLocalChainCommitmentReadAll,
		ProcessTROrgChainSubmit,
		ProcessTROrgChainPrepare,
		ProcessTROrgChainRead,
		ProcessTROrgChainReadAll,
		ProcessTRLocalChainSC,
	},
	PhaseConsistencyExamination: {
		ProcessCEMerkleProofVerify,
		ProcessCEMerkleProofMerge,
		ProcessCESummarizeMerkleProofVerificationResults,
		ProcessCEVerifyCommitments,
		ProcessCEAccumulateCommitments,
		ProcessCEAudChainSubmit,
		ProcessCEAudChainPrepare,
		ProcessCEAudChainRead,
		ProcessCEAudChainReadAll,
	},
}

func GetPhases() string {
	var phases string
	for key := range PhaseProcessMap {
		phases += key + ", "
	}
	return phases
}

func GetPhasesAndProcesses() string {
	var result string
	for key, val := range PhaseProcessMap {
		result += key + ": "
		for _, process := range val {
			result += process + ", "
		}
		result += "\n"
	}
	// remove the last "\n"
	result = result[:len(result)-1]
	return result
}
