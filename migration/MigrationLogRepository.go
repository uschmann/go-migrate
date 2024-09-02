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
	config     *Config
}

func NewMigrationLogRepository(connection *sql.DB, config *Config) *MigrationLogRepository {
	return &MigrationLogRepository{
		connection: connection,
		config:     config,
	}
}

func (m *MigrationLogRepository) MigrationLogsTableExists() bool {

	_, err := m.connection.Query("Select COUNT(*) FROM " + m.config.MigrationLogTable)

	if err != nil {
		return false
	}

	return true
}

func (m *MigrationLogRepository) CreateMigrationLogsTable() (bool, error) {
	if m.MigrationLogsTableExists() {
		return false, nil
	}

	ddl := `CREATE TABLE ` + m.config.MigrationLogTable + ` (
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
	return m.connection.Exec("insert into "+m.config.MigrationLogTable+" (name, batch) values (:1, :2)", name, batch)
}

func (m *MigrationLogRepository) DeleteMigrationLogById(id int) (sql.Result, error) {
	return m.connection.Exec("delete from "+m.config.MigrationLogTable+" where id = :1", id)
}

func (m *MigrationLogRepository) DeleteMigrationLogByName(name string) (sql.Result, error) {
	return m.connection.Exec("delete from "+m.config.MigrationLogTable+" where name = :1", name)
}

func (m *MigrationLogRepository) GetHighestBatch() (int, error) {
	var batch int
	err := m.connection.QueryRow("select nvl(max(batch), 0) from " + m.config.MigrationLogTable).Scan(&batch)

	if err != nil {
		return -1, err
	}

	return batch, nil
}

func (m *MigrationLogRepository) GetAllMigrationLogs() ([]MigrationLog, error) {
	rows, err := m.connection.Query("SELECT ID, NAME, BATCH FROM " + m.config.MigrationLogTable + " ORDER BY NAME ASC, BATCH ASC")

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

func (m *MigrationLogRepository) IsMigrationExecuted(name string) (bool, error) {
	var count int
	err := m.connection.QueryRow("select count(*) from "+m.config.MigrationLogTable+" where name = :1", name).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *MigrationLogRepository) GetMigrationLogsToRollback() ([]string, error) {
	batch, err := m.GetHighestBatch()
	check(err)

	rows, err := m.connection.Query("SELECT NAME FROM "+m.config.MigrationLogTable+" WHERE BATCH = :1 ORDER BY NAME DESC", batch)

	if err != nil {
		return nil, err
	}

	var migrations []string

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return migrations, err
		}
		migrations = append(migrations, name)
	}

	if err = rows.Err(); err != nil {
		return migrations, err
	}

	return migrations, nil
}

// TODO: Query all MigrationLogs
// See: https://go.dev/doc/database/querying
