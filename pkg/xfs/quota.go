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
	"github.com/pkg/errors"
)

const (
	FILE_PROJECTS = "/etc/projects"
	FILE_PROJID   = "/etc/projid"
)

var prjManager *PrjManager

func DeletePrj(name string) error {
	// limit project quota to zero
	if err := limitPrjQuota(name, "0"); err != nil {
		return err
	}

	// remove project directory
	if err := prjManager.Delete(name); err != nil {
		return errors.Wrap(err, "delete project")
	}

	return nil
}

// CreatePrj
func CreatePrj(name, quota string) error {
	// create project directory
	dir := fmt.Sprintf("%s/%s", conf.WORKSPACE, name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModeDir|0755); err != nil {
			return err
		}
	}

	// generate project ID
	var prjID = genPrjID()

	prj := &project{
		name: name,
		id:   prjID,
		dir:  dir,
	}

	if err := prjManager.Add(prj); err != nil {
		return errors.Wrap(err, "add project")
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

func limitPrjQuota(prjName, quota string) error {
	limitQuotaCmd := fmt.Sprintf("xfs_quota -x -c 'limit -p bsoft=%s bhard=%s %s' %s", quota, quota, prjName, conf.WORKSPACE)

	limitQuotaExecCmd := exec.Command("sh", "-c", limitQuotaCmd)
	if err := limitQuotaExecCmd.Run(); err != nil {
		return errors.Wrap(err, "limit project quota")
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

	// init project mananger
	prjManager = newPrjManager()
}
