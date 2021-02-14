package xfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

/*
	Project means a user directory, it has unique id, name and path
*/
type Project struct {
	name string
	dir  string
	id   string
}

// mapping maintains the Project name and id relationship in
// /etc/Projects and /etc/projid files
type PrjManager struct {
	Items []*Project
	sync.Mutex
}

func NewPrjManager() *PrjManager {
	return &PrjManager{}
}

// Add insert or update prj info to projects and projid files
func (m *PrjManager) Add(prj *Project) error {
	m.Lock()
	defer m.Unlock()

	if err := m.readMappingInfo(); err != nil {
		return errors.Wrap(err, "read mapping info error")
	}

	var exist bool
	for _, item := range m.Items {
		// if the prj has exist, update dir not id
		if item.name == prj.name {
			item.dir = prj.dir
			exist = true
			break
		}
	}

	if !exist {
		m.Items = append(m.Items, prj)
	}

	if err := m.persistence(); err != nil {
		return errors.Wrap(err, "overwrite to disk")
	}

	return nil
}

// Delete remove the project from projects and projid files
func (m *PrjManager) Delete(name string) error {
	m.Lock()
	defer m.Unlock()

	if err := m.readMappingInfo(); err != nil {
		return errors.Wrap(err, "read mapping info error")
	}

	for index, item := range m.Items {
		if item.name == name {
			m.Items = append(m.Items[0:index], m.Items[index+1:]...)
			break
		}
	}

	if err := m.persistence(); err != nil {
		return errors.Wrap(err, "persistence to disk")
	}

	return nil
}

func (m *PrjManager) readMappingInfo() error {
	// read /etc/projid file
	projidFile, err := ioutil.ReadFile(FILE_PROJID)
	if err != nil {
		return errors.Wrap(err, "read projid file error")
	}
	mappingIDToName := make(map[string]string)
	lines := strings.Split(string(projidFile), "\n")

	// iterate all lines
	for i := 0; i < len(lines); i++ {
		info := strings.Split(lines[i], ":")
		name, id := info[0], info[1]
		mappingIDToName[id] = name
	}

	// read /etc/Projects file
	ProjectsFile, err := ioutil.ReadFile(FILE_PROJECTS)
	if err != nil {
		return errors.Wrap(err, "read Projects file error")
	}
	mappingIDToDir := make(map[string]string)
	lines = strings.Split(string(ProjectsFile), "\n")

	// iterate all lines but last empty line
	for i := 0; i < len(lines); i++ {
		info := strings.Split(lines[i], ":")
		id, dir := info[0], info[1]
		mappingIDToDir[id] = dir
	}

	// populate items
	// make items slice empty first
	m.Items = make([]*Project, 0, len(mappingIDToName))
	for id, name := range mappingIDToName {
		m.Items = append(m.Items, &Project{
			name: name,
			id:   id,
			dir:  mappingIDToDir[id],
		})
	}

	return nil
}

func (m *PrjManager) persistence() error {
	// prepare projid and Projects file content
	var projidContent = make([]string, 0, len(m.Items))
	var ProjectsContent = make([]string, 0, len(m.Items))

	for _, item := range m.Items {
		var projidLine = fmt.Sprintf("%s:%s", item.name, item.id)
		projidContent = append(projidContent, projidLine)

		var ProjectsLine = fmt.Sprintf("%s:%s", item.id, item.dir)
		ProjectsContent = append(ProjectsContent, ProjectsLine)
	}

	// sync content to disk
	if err := overwriteFile(FILE_PROJID, strings.Join(projidContent, "\n")); err != nil {
		return errors.Wrap(err, "overwrite /etc/projid file")
	}

	if err := overwriteFile(FILE_PROJECTS, strings.Join(ProjectsContent, "\n")); err != nil {
		return errors.Wrap(err, "overwrite /etc/Projects file")
	}

	return nil
}

// overwriteFile overwrite file with content
func overwriteFile(file, content string) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrap(err, "open file")
	}

	_, err = f.WriteString(content)
	if err != nil {
		return errors.Wrap(err, "rewrite file")
	}

	return nil
}
