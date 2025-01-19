package service

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Анонимный импорт пакета драйвера
)

type PostgresDBParams struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

type PostgresTransactionLogger struct {
	events chan<- Event // Канал только для записи; для передачи событий
	errors <-chan error // Канал только для чтения; для приема ошибок
	db     *sql.DB
}

func (p *PostgresTransactionLogger) verifyTableExist() error {
	_, err := p.db.Exec(`
CREATE SEQUENCE IF NOT EXISTS sequence_default START 1 INCREMENT 1;
CREATE TABLE IF NOT EXISTS transactions(sequence INT NOT NULL default nextval('sequence_default'), event_type INT, key VARCHAR(25), value TEXT);`)
	if err != nil {
		return fmt.Errorf("not exist Table %w", err)
	}
	return nil
}

func (p *PostgresTransactionLogger) createTable() error {
	_, err := p.db.Exec("CREATE TABLE $1 (int sequence, int event_type, string key , string value)", "transaction")
	if err != nil {
		return fmt.Errorf("create table %w", err)
	}
	return nil
}

func NewPostgresTransactionLogger(params PostgresDBParams) (TransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", params.Host, params.Port, params.User, params.Password, params.DbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}
	logger := &PostgresTransactionLogger{db: db}

	err = logger.verifyTableExist()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}
	//if !exists {
	//	if err = logger.createTable(); err != nil {
	//		return nil, fmt.Errorf("failed to create table: %w", err)
	//	}
	//}
	return logger, nil
}

func (p *PostgresTransactionLogger) WritePut(key, value string) {
	p.events <- Event{EventType: EventPut, Key: key, Value: value}
}
func (p *PostgresTransactionLogger) WriteDelete(key string) {
	p.events <- Event{EventType: EventDelete, Key: key, Value: "deleted"}
}

func (p *PostgresTransactionLogger) Error() <-chan error {
	return p.errors
}

func (p *PostgresTransactionLogger) Run() {
	events := make(chan Event, 16)
	p.events = events
	errors := make(chan error, 1)
	p.errors = errors
	go func() {
		query := `INSERT INTO transactions (event_type, key, value) VALUES ($1, $2, $3)`
		for e := range events {
			fmt.Printf("INSERT")
			_, err := p.db.Exec(query, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}

	}()
}
func (p *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvents := make(chan Event) // Небуферизованный канал событий
	outErrors := make(chan error, 1)
	go func() {
		defer close(outEvents)
		defer close(outErrors)

		query := `SELECT sequence, event_type, key, value FROM transactions ORDER BY sequence`
		rows, err := p.db.Query(query)
		if err != nil {
			outErrors <- fmt.Errorf("sql query error: %w", err)
		}
		defer rows.Close()
		event := Event{}
		for rows.Next() {
			err := rows.Scan(&event.Sequence, &event.EventType, &event.Key, &event.Value)
			if err != nil {
				outErrors <- fmt.Errorf("error reading row: %w", err)
				return
			}
			outEvents <- event
		}
		err = rows.Err()
		if err != nil {
			outErrors <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()
	return outEvents, outErrors
}

func (p *PostgresTransactionLogger) Initialize(storage *Storage) error {
	var err error
	events, errorsLog := p.ReadEvents()
	e, ok := Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errorsLog:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = storage.Delete(e.Key)
			case EventPut:
				err = storage.Put(e.Key, string(e.Value))
			}
		}
	}
	p.Run()
	return err
}

func (p *PostgresTransactionLogger) Close() error {
	dbErr := p.db.Close()
	if dbErr != nil {
		return fmt.Errorf("close db: %w", dbErr)
	}
	return nil
}
