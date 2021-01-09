package xfs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/Wang-Kai/quotar/pkg/conf"
)

const (
	FILE_PROJECTS = "/etc/projects"
	FILE_PROJID   = "/etc/projid"
)

// CreatePrj
func CreatePrj(name, quota string) error {
	// create project directory
	prjPath := fmt.Sprintf("%s/%s", conf.WORKSPACE, name)
	if _, err := os.Stat(prjPath); os.IsNotExist(err) {
		if err := os.Mkdir(prjPath, os.ModeDir|0755); err != nil {
			return err
		}
	}

	var prjID = genPrjID()

	// insert the mapping of project ID and directory
	var mappingIDAndDir = fmt.Sprintf("%s:%s\n", prjID, prjPath)

	projectsFilePointer, err := os.OpenFile(FILE_PROJECTS, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer projectsFilePointer.Close()

	_, err = projectsFilePointer.WriteString(mappingIDAndDir)
	if err != nil {
		return err
	}

	// insert the mapping of project name and project ID
	var mappingNameAndID = fmt.Sprintf("%s:%s\n", name, prjID)

	projidFilePointer, err := os.OpenFile(FILE_PROJID, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer projidFilePointer.Close()

	_, err = projidFilePointer.WriteString(mappingNameAndID)
	if err != nil {
		return err
	}

	// init xfs quota project
	initQuotaCmd := fmt.Sprintf("xfs_quota -x -c 'project -s %s' %s", name, conf.WORKSPACE)
	initQuotaExecCmd := exec.Command("bash", "-c", initQuotaCmd)
	if err := initQuotaExecCmd.Run(); err != nil {
		return err
	}

	// limit project quota
	limitQuotaCmd := fmt.Sprintf("xfs_quota -x -c 'limit -p bsoft=%s bhard=%s %s' %s", quota, quota, name, conf.WORKSPACE)
	limitQuotaExecCmd := exec.Command("sh", "-c", limitQuotaCmd)
	if err := limitQuotaExecCmd.Run(); err != nil {
		return err
	}

	return nil
}

var currentPrjID uint32

// genPrjID generate project ID while add xfs project
func genPrjID() string {
	nextPrjID := atomic.AddUint32(&currentPrjID, 1)
	return fmt.Sprintf("%d", nextPrjID)
}

func init() {
	// get current max project ID
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

	currentPrjID = maxProjID

	println("Current max project ID is", currentPrjID)
}
