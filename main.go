package main

import (
	"fmt"
	"strconv"

	"github.com/jllucas/minitree/history"
)

func main() {
	var rootHash []byte

	// Create tree
	hTree := history.NewHistoryTree()
	fmt.Println(hTree)

	// Insert events.
	for i := 0; i <= 7; i++ {
		fmt.Println("-> Adding event", i)
		rootHash = hTree.Add([]byte(strconv.Itoa((i))))
		fmt.Printf("Root hash --> %x\n", rootHash)
	}

	// Generate membership proof.
	eventIndex := 2
	hash := []byte{0xfd, 0x9a, 0xf3, 0xd2, 0x89, 0xe1, 0xc7, 0xb2, 0x2c, 0x64, 0x8a, 0x4f, 0x31, 0x63, 0x81, 0xe7, 0xc, 0x63, 0xf8, 0x29, 0x4d, 0x3d, 0xc5, 0x7d, 0xbe, 0xc8, 0x2d, 0x32, 0x92, 0x40, 0x91, 0x5}
	version := 2

	membershipProof := hTree.GenerateMembershipProof(eventIndex, hash, version)
	fmt.Printf("\nMembership Proof for event %d in version-%d tree. \n", eventIndex, version)
	membershipProof.Prettyfy()

	// Generate incremental proof.
	eventI := 2
	eventJ := 4

	incrementalProof := hTree.GenerateIncrementalProof(eventI, eventJ)
	fmt.Printf("\nIncremental Proof between events %d and %d. \n", eventI, eventJ)
	incrementalProof.Prettyfy()

	// Verify membership proof.
	/* 	isMember := membershipProof.VerifyMembershipProof(eventIndex, hash, hash)
	   	fmt.Printf("\nIs event %d member of version-%d tree? %v\n", eventIndex, version, isMember)
	*/
	// Verify incremental proof.
	/* 	isCorrect := incrementalProof.VerifyIncrementalProof(hash, hash)
	   	fmt.Printf("Is Incremental proof correct? %v", isCorrect) */
}
