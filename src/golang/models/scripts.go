package models

import (
	"database/sql"
)

//AnsibleScript - sctruct to provide sql tables mapping
type AnsibleScript struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

//AnsibleScriptCollection -- struct provide collection of AnsibleScripts
type AnsibleScriptCollection struct {
	Scripts []AnsibleScript `json:"items"`
}

//GetAnsibleScripts -- get collection of scrits from db
func GetAnsibleScripts(db *sql.DB) (AnsibleScriptCollection, error) {
	sql := "SELECT * FROM ansible_scripts"
	rows, err := db.Query(sql)
	if err != nil {
		return AnsibleScriptCollection{}, err
	}
	defer rows.Close()

	result := AnsibleScriptCollection{}
	for rows.Next() {
		ansibleScript := AnsibleScript{}
		err = rows.Scan(&ansibleScript.ID, &ansibleScript.Name, &ansibleScript.Path, &ansibleScript.Description)
		if err != nil {
			return AnsibleScriptCollection{}, err
		}
		result.Scripts = append(result.Scripts, ansibleScript)
	}
	return result, nil
}

//GetAnsibleScripts -- get collection of scrits from db
func GetAnsibleScriptsFiltered(db *sql.DB, filter string) (AnsibleScriptCollection, error) {
	sql := "SELECT * FROM ansible_scripts where path like '%" + filter + "%'"
	rows, err := db.Query(sql)
	if err != nil {
		return AnsibleScriptCollection{}, err
	}
	defer rows.Close()

	result := AnsibleScriptCollection{}
	for rows.Next() {
		ansibleScript := AnsibleScript{}
		err = rows.Scan(&ansibleScript.ID, &ansibleScript.Name, &ansibleScript.Path, &ansibleScript.Description)
		if err != nil {
			return AnsibleScriptCollection{}, err
		}
		result.Scripts = append(result.Scripts, ansibleScript)
	}
	return result, nil
}

//PutAnsibleScripts -- put script to database
func PutAnsibleScripts(db *sql.DB, script AnsibleScript) error {
	//Пока нет мыслей зачем это понадобится
	return nil
}

//DeleteAnsibleScripts -- delete scripts from db
func DeleteAnsibleScripts(db *sql.DB, script AnsibleScript) error {
	//Пока нет мыслей зачем это понадобится
	return nil
}
