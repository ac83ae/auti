package main

import (
	"flag"
	"log"

	. "github.com/ac83ae/auti/benchmark/closc/internal/flag"
	"github.com/ac83ae/auti/benchmark/closc/internal/task"
)

func main() {
	benchPhasePtr := flag.String("phase", "", GetPhases())
	benchProcessPtr := flag.String("process", "", GetPhasesAndProcesses())
	numOrgPtr := flag.Int("numOrg", 2, "Number of organizations")
	numIterPtr := flag.Int("numIter", 10, "Number of iterations")
	numPtr := flag.Int("num", 100, "Number/Quantity/Depth/Number of SC")
	flag.Parse()

	var err error
	switch *benchPhasePtr {
	case PhaseInitialization:
		switch *benchProcessPtr {
		case ProcessINDefault:
			err = task.INEpoch(*numOrgPtr, *numIterPtr)
		case ProcessINRandGen:
			err = task.INRandGen(*numPtr, *numIterPtr)
		}
	case PhaseTransactionRecord:
		switch *benchProcessPtr {
		case ProcessTRCommitment:
			err = task.TRCommitment(*numPtr, *numIterPtr)
		case ProcessTRMerkleProofGen:
			err = task.TRMerkleProofGen(*numPtr, *numIterPtr)
		case ProcessTRLocalChainSubmit:
			err = task.TRLocalSubmitTX(*numPtr, *numIterPtr)
		case ProcessTRLocalChainPrepare:
			err = task.TRLocalPrepareTX(*numPtr)
		case ProcessTRLocalChainRead:
			err = task.TRLocalReadTX(*numPtr, *numIterPtr)
		case ProcessTRLocalChainReadAll:
			err = task.TRLocalReadAllTXs(*numPtr, *numIterPtr)
		case ProcessTRLocalChainCommitmentSubmit:
			err = task.TRLocalCommitmentSubmitTX(*numPtr, *numIterPtr)
		case ProcessTRLocalChainCommitmentPrepare:
			err = task.TRLocalCommitmentPrepareTX(*numPtr)
		case ProcessTRLocalCHainCommitmentRead:
			err = task.TRLocalCommitmentReadTX(*numPtr, *numIterPtr)
		case ProcessTRLocalChainCommitmentReadAll:
			err = task.TRLocalCommitmentReadAllTXs(*numPtr, *numIterPtr)
		case ProcessTROrgChainSubmit:
			err = task.TROrgSubmitTX(*numPtr, *numIterPtr)
		case ProcessTROrgChainPrepare:
			err = task.TROrgPrepareTX(*numPtr)
		case ProcessTROrgChainRead:
			err = task.TROrgReadTX(*numPtr, *numIterPtr)
		case ProcessTROrgChainReadAll:
			err = task.TROrgReadAllTXs(*numPtr, *numIterPtr)
		case ProcessTRLocalChainSC:
			err = task.TRLocalSC(*numPtr, *numIterPtr)
		}
	case PhaseConsistencyExamination:
		switch *benchProcessPtr {
		case ProcessCEMerkleProofVerify:
			err = task.CEMerkleProofVerify(*numPtr, *numIterPtr)
		case ProcessCEMerkleProofMerge:
			err = task.CEMerkleProofMerge(*numPtr, *numIterPtr)
		case ProcessCESummarizeMerkleProofVerificationResults:
			err = task.CESummarizeMerkleProofVerificationResults(*numPtr, *numIterPtr)
		case ProcessCEVerifyCommitments:
			err = task.CEVerifyCommitments(*numPtr, *numIterPtr)
		case ProcessCEAccumulateCommitments:
			err = task.CEAccumulateCommitments(*numPtr, *numIterPtr)
		case ProcessCEAudChainSubmit:
			err = task.CEAudSubmitTX(*numPtr, *numIterPtr)
		case ProcessCEAudChainPrepare:
			err = task.CEAudPrepareTX(*numPtr)
		case ProcessCEAudChainRead:
			err = task.CEAudReadTX(*numPtr, *numIterPtr)
		case ProcessCEAudChainReadAll:
			err = task.CEAudReadAllTXs(*numPtr, *numIterPtr)
		}
	case PhaseResultVerification:
		switch *benchProcessPtr {
		case ProcessRVVerifyMerkleBatchProof:
			err = task.RVVerifyMerkleBatchProof(*numPtr, *numIterPtr)
		case ProcessRVSummarizeMerkleBatchProofVerificationResults:
			err = task.RVSummarizeMerkleBatchProofVerificationResults(*numPtr, *numIterPtr)
		case ProcessRVVerifyCommitments:
			err = task.RVVerifyCommitments(*numPtr, *numIterPtr)
		}
	default:
		log.Fatalf("Error: %v", "Invalid phase")
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
