package xfs

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var latestPrjID uint32
var prjManager *PrjManager

/*
	init function will load /etc/projects file to record max project ID for compute next
	besides, it also create Manager object for operations
*/
func init() {
	// get current max Project ID
	var maxProjID uint32

	f, err := os.Open(FILE_PROJECTS)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		projID := strings.Split(line, ":")[0]

		id, err := strconv.Atoi(projID)
		if err != nil {
			log.Fatal(err)
		}

		if uint32(id) > maxProjID {
			maxProjID = uint32(id)
		}
	}

	latestPrjID = maxProjID

	println("Latest project ID is: ", latestPrjID)

	// init Project mananger
	prjManager = NewPrjManager()
}
