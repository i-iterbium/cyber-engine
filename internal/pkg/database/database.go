package database

import (
	"database/sql"
	"fmt"
	"sync"
)

// Command опсиывает настройки запроса
type Command struct {
	Query string        `json:"query,omitempty"`
	Args  []interface{} `json:"args,omitempty"`
}

// ConnectionSetting опсиывает настройки подключения к БД
type ConnectionSetting struct {
	Driver                 string                 `json:"driver,omitempty"`
	ConnectionStringParams map[string]interface{} `json:"connectionStringParams,omitempty"`
	AfterConnection        []Command              `json:"afterConnection,omitempty"`
}

// DriverSetting описыает настройки драйвер к БД
type DriverSetting struct {
	GetConnectionString func(s map[string]interface{}) (string, error)
	AfterConnection     func(db *sql.DB, cs ConnectionSetting) error
}

// Settings описывает настройки БД по умолчанию и пул соединений
type Settings struct {
	Pool map[string]ConnectionSetting `json:"pool"`
}

// Pool описывает потокобезопасный пул именованных соединений
type Pool struct {
	*sync.RWMutex
	m map[string]*sql.DB
}

// NewPool создаёт новый пул и возвращает указатель на него
func NewPool() *Pool {
	return &Pool{
		RWMutex: &sync.RWMutex{},
		m:       make(map[string]*sql.DB),
	}
}

var settings Settings
var pool *Pool
var drivers = make(map[string]DriverSetting)

func init() {
	pool = NewPool()
}

// GetSettings возвращает указатель на настройки БД
func GetSettings() *Settings {
	return &settings
}

// RegisterDriver регистрирует драйвер БД и задаёт функцию, возвращающую connection string для подключения
func RegisterDriver(name string, driverSetting DriverSetting) {
	drivers[name] = driverSetting
}

// Open устанавливает соединение с базой данных из пула, если соединение ещё не установлено, и возвращает ссылку на него
func Open(name string) (*sql.DB, error) {
	connectionSettings, ok := settings.Pool[name]
	if !ok {
		return nil, fmt.Errorf("%q отсутствует в пуле соединений", name)
	}

	conn, _ := pool.Load(name)
	conn, err := getConn(conn, connectionSettings)
	if err != nil {
		return nil, err
	}

	pool.Store(name, conn)

	return conn, err
}

// Load возвращает соединение из пула
func (p *Pool) Load(key string) (*sql.DB, bool) {
	p.RLock()
	value, ok := p.m[key]
	p.RUnlock()

	return value, ok
}

// Store сохраняет соединение в пул
func (p *Pool) Store(key string, value *sql.DB) {
	p.Lock()
	p.m[key] = value
	p.Unlock()
}

func getConn(conn *sql.DB, cs ConnectionSetting) (*sql.DB, error) {
	if conn == nil {
		return openConn(cs)
	}

	if err := conn.Ping(); err != nil {
		return openConn(cs)
	}

	return conn, nil
}

func openConn(s ConnectionSetting) (*sql.DB, error) {
	ds, ok := drivers[s.Driver]
	if !ok {
		return nil, fmt.Errorf("Драйвер для %q не зарегистрирован", s.Driver)
	}

	cs, err := ds.GetConnectionString(s.ConnectionStringParams)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open(s.Driver, cs)
	if err != nil {
		return nil, err
	}

	if ds.AfterConnection != nil {
		ds.AfterConnection(conn, s)
	}

	for _, q := range s.AfterConnection {
		if _, err := conn.Exec(q.Query, q.Args...); err != nil {
			return nil, err
		}
	}

	return conn, err
}
