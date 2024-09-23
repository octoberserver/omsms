package enums

import "errors"

type BackupStrat string

const (
	BACKUP_FULL_SERVER BackupStrat = "FULL_SERVER"
	BACKUP_WORLD       BackupStrat = "WORLD"
	BACKUP_CUSTOM      BackupStrat = "CUSTOM"
	BACKUP_NONE        BackupStrat = "NONE"
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
		return errors.New(`must be one of "FULL_SERVER", "WORLD", "CUSTOM" or "NONE"`)
	}
}

func (e *BackupStrat) Type() string {
	return "BackupStrat"
}
