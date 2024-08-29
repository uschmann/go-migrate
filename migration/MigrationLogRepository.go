package migration

import (
	"database/sql"
)

type MigrationLog struct {
	Id    int
	Name  string
	Batch int
}

type MigrationLogRepository struct {
	connection *sql.DB
}

func NewMigrationLogRepository(connection *sql.DB) *MigrationLogRepository {
	return &MigrationLogRepository{
		connection: connection,
	}
}

func (m *MigrationLogRepository) MigrationLogsTableExists() bool {

	_, err := m.connection.Query("Select COUNT(*) FROM MIGRATION_LOGS")

	if err != nil {
		return false
	}

	return true
}

func (m *MigrationLogRepository) CreateMigrationLogsTable() (bool, error) {
	if m.MigrationLogsTableExists() {
		return false, nil
	}

	ddl := `CREATE TABLE MIGRATION_LOGS (
				ID NUMBER(6,0) GENERATED ALWAYS AS IDENTITY,
				NAME VARCHAR2(256) NOT NULL,
				BATCH NUMBER(6,0) NOT NULL,

				CONSTRAINT migration_logs_pk PRIMARY KEY (ID)
			)`

	_, err := m.connection.Exec(ddl)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *MigrationLogRepository) AddMigrationLog(name string, batch int) (sql.Result, error) {
	return m.connection.Exec("insert into migration_logs (name, batch) values (:1, :2)", name, batch)
}

func (m *MigrationLogRepository) DeleteMigrationLogById(id int) (sql.Result, error) {
	return m.connection.Exec("delete from migration_logs where id = :1", id)
}

func (m *MigrationLogRepository) GetHighestBatch() (int, error) {
	var batch int
	err := m.connection.QueryRow("select max(batch) from migration_logs").Scan(&batch)

	if err != nil {
		return -1, err
	}

	return batch, nil
}

func (m *MigrationLogRepository) GetAllMigrationLogs() ([]MigrationLog, error) {
	rows, err := m.connection.Query("SELECT ID, NAME, BATCH FROM migration_logs")

	if err != nil {
		return nil, err
	}

	var migrationLogs []MigrationLog

	for rows.Next() {
		var log MigrationLog

		if err := rows.Scan(&log.Id, &log.Name, &log.Batch); err != nil {
			return migrationLogs, err
		}
		migrationLogs = append(migrationLogs, log)
	}

	if err = rows.Err(); err != nil {
		return migrationLogs, err
	}

	return migrationLogs, nil
}

// TODO: Query all MigrationLogs
// See: https://go.dev/doc/database/querying
