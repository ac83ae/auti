package task

import (
	"fmt"
	"time"

	"go.dedis.ch/kyber/v3"

	"github.com/ac83ae/auti/benchmark/timecounter"
	closcaud "github.com/ac83ae/auti/internal/closc/auditor"
	"github.com/ac83ae/auti/internal/closc/committee"
	closccom "github.com/ac83ae/auti/internal/closc/committee"
	closcorg "github.com/ac83ae/auti/internal/closc/organization"
)

func generateEntities(numOrganizations int) (*closccom.Committee, []*closcaud.Auditor, []*closcorg.Organization) {
	organizations := make([]*closcorg.Organization, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		organizations[i] = closcorg.New("org" + string(rune(i)))
	}
	auditors := make([]*closcaud.Auditor, numOrganizations)
	for i := 0; i < numOrganizations; i++ {
		auditors[i] = closcaud.New("aud"+string(rune(i)), []*closcorg.Organization{organizations[i]})
	}
	com := closccom.New("com", auditors)
	return com, auditors, organizations
}

func INEpoch(numOrganizations, iterations int) error {
	fmt.Println("[CLOSC-IN] Default")
	fmt.Printf("Num Org: %d, Num iter: %d\n", numOrganizations, iterations)
	for i := 0; i < iterations; i++ {
		com, auditors, _ := generateEntities(numOrganizations)
		startTime := time.Now()
		err := com.InitializeEpoch(auditors)
		if err != nil {
			return err
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}

func INRandGen(num, iterations int) error {
	fmt.Println("[CLOSC-IN] Random generation")
	fmt.Printf("Num: %d, Num iter: %d\n", num, iterations)
	for i := 0; i < iterations; i++ {
		startTime := time.Now()
		epochIDs := make([]kyber.Point, num)
		for j := 0; j < num; j++ {
			epochIDs[j] = committee.GenerateAuditorEpochID()
		}
		elapsed := time.Since(startTime)
		timecounter.Print(elapsed)
	}
	fmt.Println()
	return nil
}
