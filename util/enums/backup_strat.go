package enums

import "errors"

type BackupStrat string

const (
	BACKUP_FULL_SERVER BackupStrat = "FULL_SERVER"
	BACKUP_WORLD       BackupStrat = "WORLD"
	BACKUP_CUSTOM      BackupStrat = "CUSTOM"
	BACKUP_NONE        BackupStrat = "NONE"
	BACKUP_NULL        BackupStrat = ""
)

func (e *BackupStrat) String() string {
	return string(*e)
}

func (e *BackupStrat) Set(v string) error {
	switch v {
	case "FULL_SERVER", "WORLD", "CUSTOM", "NONE":
		*e = BackupStrat(v)
		return nil
	default:
		return errors.New(`必須是 "FULL_SERVER", "WORLD", "CUSTOM" 或 "NONE"`)
	}
}

func (e *BackupStrat) Type() string {
	return "BackupStrat"
}
