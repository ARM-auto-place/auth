package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ARM-auto-place/auth/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type ConnectManager interface {
	GetMasterConn() *pgxpool.Pool
	GetReplicaConn() *pgxpool.Pool
}

type manager struct {
	log     zerolog.Logger
	master  *pgxpool.Pool
	replica *pgxpool.Pool
}

func (m *manager) GetMasterConn() *pgxpool.Pool {
	return m.master
}

func (m *manager) GetReplicaConn() *pgxpool.Pool {
	return m.replica
}

func New(log zerolog.Logger, cfg config.PostgresConfig) (ConnectManager, error) {
	ctx := context.Background()

	masterAddr, masterPort, err := parseDbAddressAndPort(cfg.Addr)
	if err != nil {
		return nil, err
	}

	master, err := dbConnect(ctx, &cfg, masterAddr, masterPort, cfg.DB, cfg.User, cfg.Password)
	if err != nil {
		return nil, err
	}

	replicaAddr, replicaPort, err := parseDbAddressAndPort(cfg.ReplicaAddr)
	if err != nil {
		return nil, err
	}

	replica, err := dbConnect(ctx, &cfg, replicaAddr, replicaPort, cfg.DB, cfg.UserRO, cfg.PasswordRO)
	if err != nil {
		return nil, err
	}

	return &manager{
		log:     log,
		master:  master,
		replica: replica,
	}, nil
}

func parseDbAddressAndPort(conn string) (string, int, error) {
	splits := strings.Split(conn, ":")
	address := splits[0]
	port, err := strconv.Atoi(splits[1])
	if err != nil {
		return "", 0, fmt.Errorf("failed parse db address and port, connection string=%s", conn)
	}

	return address, port, nil
}

func dbConnect(ctx context.Context, pcfg *config.PostgresConfig, dbAddr string, dbPort int, db, user, password string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d dbname=%s sslmode=disable user=%s password=%s pool_max_conns=%d",
		dbAddr, dbPort, db, user, password, pcfg.MaxConn,
	))

	if err != nil {
		return nil, fmt.Errorf("failed parse postgres dsn: %s:%v: %w", dbAddr, dbPort, err)
	}

	mci, err := time.ParseDuration(pcfg.MaxIdleLifetime)
	if err != nil {
		return nil, fmt.Errorf("failed parse max idle conn lifetime to duration: %w", err)
	}
	mc, err := time.ParseDuration(pcfg.MaxIdleLifetime)
	if err != nil {
		return nil, fmt.Errorf("failed parse max conn lifetime to duration: %w", err)
	}

	cfg.MaxConnIdleTime = mci
	cfg.MaxConnLifetime = mc

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed connect to pg %s:%v: %w", dbAddr, dbPort, err)
	}

	return pool, nil
}
