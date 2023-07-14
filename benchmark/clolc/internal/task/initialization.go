package task

import (
	"fmt"
	"time"

	"github.com/ac83ae/auti/benchmark/timecounter"
	clolcaud "github.com/ac83ae/auti/internal/clolc/auditor"
	clolccom "github.com/ac83ae/auti/internal/clolc/committee"
	clolcorg "github.com/ac83ae/auti/internal/clolc/organization"
)

func generateEntities(numOrganizations int) (*clolccom.Committee, []*clolcaud.Auditor, []*clolcorg.Organization) {
	organizations := make([]*clolcorg.Organization, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		organizations[i] = clolcorg.New("org" + string(rune(i)))
	}
	auditors := make([]*clolcaud.Auditor, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		auditors[i] = clolcaud.New("aud"+string(rune(i)), []*clolcorg.Organization{organizations[i]})
	}
	com := clolccom.New("com", auditors)
	return com, auditors, organizations
}

func INDefault(numOrganizations, iterations int) error {
	fmt.Println("[CLOLC-IN] Default")
	fmt.Printf("Num Org: %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, organizations := generateEntities(numOrganizations)
		startTime := time.Now()
		_, err := com.InitializeEpoch(auditors, organizations)
		if err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}
