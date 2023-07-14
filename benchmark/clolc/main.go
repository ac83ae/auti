package main

import (
	"flag"
	"log"

	. "github.com/ac83ae/auti/benchmark/clolc/internal/flag"
	bf "github.com/ac83ae/auti/benchmark/clolc/internal/task"
)

func main() {
	benchPhasePtr := flag.String("phase", "", GetPhases())
	benchProcessPtr := flag.String("process", "", GetPhasesAndProcesses())
	numOrgPtr := flag.Int("numOrg", 2, "Number of organizations")
	numIterPtr := flag.Int("numIter", 10, "Number of iterations")
	numTXsPtr := flag.Int("numTXs", 100, "Number of transactions")
	flag.Parse()

	var err error
	switch *benchPhasePtr {
	case PhaseInitialization:
		switch *benchProcessPtr {
		case ProcessINDefault:
			err = bf.INDefault(*numOrgPtr, *numIterPtr)
		}
	case PhaseTransactionRecord:
		switch *benchProcessPtr {
		case ProcessTRLocalChainSubmit:
			err = bf.TRLocalSubmitTX(*numTXsPtr, *numIterPtr)
		case ProcessTRLocalChainPrepare:
			err = bf.TRLocalPrepareTX(*numTXsPtr)
		case ProcessTRLocalChainRead:
			err = bf.TRLocalReadTX(*numTXsPtr, *numIterPtr)
		case ProcessTRLocalChainReadAll:
			err = bf.TRLocalReadAllTXs(*numTXsPtr, *numIterPtr)
		case ProcessTROrgChainSubmit:
			err = bf.TROrgSubmitTX(*numTXsPtr, *numIterPtr)
		case ProcessTROrgChainPrepare:
			err = bf.TROrgPrepareTX(*numTXsPtr)
		case ProcessTROrgChainRead:
			err = bf.TROrgReadTX(*numTXsPtr, *numIterPtr)
		case ProcessTROrgChainReadAll:
			err = bf.TROrgReadAllTXs(*numTXsPtr, *numIterPtr)
		case ProcessTRCommitment:
			err = bf.TRCommitment(*numTXsPtr, *numIterPtr)
		case ProcessTRAccumulate:
			err = bf.TRAccumulate(*numTXsPtr, *numIterPtr)
		}
	case PhaseConsistencyExamination:
		switch *benchProcessPtr {
		case ProcessCEAccumulateCommitment:
			err = bf.CEAccumulateCommitment(*numOrgPtr, *numIterPtr)
		case ProcessCEComputeB:
			err = bf.CEComputeB(*numOrgPtr, *numIterPtr)
		case ProcessCEComputeC:
			err = bf.CEComputeC(*numOrgPtr, *numIterPtr)
		case ProcessCEComputeD:
			err = bf.CEComputeD(*numOrgPtr, *numIterPtr)
		case ProcessCEEncrypt:
			err = bf.CEEncrypt(*numOrgPtr, *numIterPtr)
		case ProcessCEDecrypt:
			err = bf.CEDecrypt(*numIterPtr)
		case ProcessCEAudChainSubmit:
			err = bf.CEAudSubmitTX(*numTXsPtr, *numIterPtr)
		case ProcessCEAudChainPrepare:
			err = bf.CEAudPrepareTX(*numTXsPtr)
		case ProcessCEAudChainRead:
			err = bf.CEAudReadTX(*numTXsPtr, *numIterPtr)
		case ProcessCEAudChainReadAll:
			err = bf.CEAudReadAllTXs(*numTXsPtr, *numIterPtr)
		case ProcessCECheck:
			err = bf.CECheck(*numIterPtr)
		}
	case PhaseResultVerification:
		switch *benchProcessPtr {
		case ProcessRVVerifyOrgAndAudResult:
			err = bf.RVVerifyOrgAndAudResult(*numOrgPtr, *numIterPtr)
		case ProcessRVVerifyAuditPairResult:
			err = bf.RVVerifyAuditPairResult(*numOrgPtr, *numIterPtr)
		}
	default:
		log.Fatalf("Error: %v", "Invalid phase")
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
