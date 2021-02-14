package xfs

import (
	"fmt"
	"os"
	"os/exec"
	"sync/atomic"

	"github.com/Wang-Kai/quotar/pkg/conf"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	FILE_PROJECTS = "/etc/projects"
	FILE_PROJID   = "/etc/projid"
)

// DeletePrj delete the project
// 1. limit project quota to zero
// 2. remove project directory if needed
// 3. remove mapping info
func DeletePrj(name string) error {
	// limit Project quota to zero
	if err := limitPrjQuota(name, "0"); err != nil {
		return err
	}

	// TODO: remove Project directory

	// remove project info from mapping files
	if err := prjManager.Delete(name); err != nil {
		return errors.Wrap(err, "delete Project")
	}

	return nil
}

// CreatePrj create a new project
// 1. create directory
// 2. write project info into mapping files
// 3. init xfs project
// 4. limit project quota
func CreatePrj(name, quota string) error {
	// create Project directory
	dir := fmt.Sprintf("%s/%s", conf.WORKSPACE, name)
	if err := createPrjDir(dir); err != nil {
		return err
	}

	// generate Project ID
	var prjID = genPrjID()
	prj := &Project{
		name: name,
		id:   prjID,
		dir:  dir,
	}

	if err := prjManager.Add(prj); err != nil {
		return errors.Wrap(err, "add Project")
	}

	// init xfs quota Project
	if err := initPrjQuota(name, quota); err != nil {
		return err
	}

	// limit Project quota
	if err := limitPrjQuota(name, quota); err != nil {
		return err
	}

	return nil
}

// initPrjQuota init the quota setting for prj project
func initPrjQuota(prj, quota string) error {
	initQuotaCmd := fmt.Sprintf("xfs_quota -x -c 'project -s %s' %s", prj, conf.WORKSPACE)
	log.WithFields(log.Fields{
		"project": prj,
		"quota":   quota,
		"command": initQuotaCmd,
	}).Info("Init project quota")

	initQuotaExecCmd := exec.Command("bash", "-c", initQuotaCmd)
	if err := initQuotaExecCmd.Run(); err != nil {
		return errors.Wrap(err, "init project quota")
	}

	return nil
}

// limitPrjQuota execute xfs_quota command to limit prj to the size
func limitPrjQuota(prj, quota string) error {
	limitQuotaCmd := fmt.Sprintf("xfs_quota -x -c 'limit -p bsoft=%s bhard=%s %s' %s", quota, quota, prj, conf.WORKSPACE)
	log.WithFields(log.Fields{
		"project": prj,
		"quota":   quota,
		"command": limitQuotaCmd,
	}).Info("Limit project quota")

	limitQuotaExecCmd := exec.Command("sh", "-c", limitQuotaCmd)
	if err := limitQuotaExecCmd.Run(); err != nil {
		return errors.Wrap(err, "limit project quota")
	}

	return nil
}

// genPrjID generate Project ID while add xfs Project
func genPrjID() string {
	nextPrjID := atomic.AddUint32(&latestPrjID, 1)
	return fmt.Sprintf("%d", nextPrjID)
}

// createPrjDir create directory with 0755 mode
func createPrjDir(name string) error {
	log.WithField("directory", name).Info("Create project directoy")

	if err := os.Mkdir(name, os.ModeDir|0755); err != nil {
		return errors.Wrapf(err, "create %s directory", name)
	}
	return nil
}
