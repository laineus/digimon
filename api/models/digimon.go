package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

// DigimonRepository デジモンリポジトリ
type DigimonRepository struct{}

// NewDigimonRepository リポジトリ取得
func NewDigimonRepository() DigimonRepository {
	return DigimonRepository{}
}

// GetAll 全デジモンを取得
func (m DigimonRepository) GetAll() *[]Digimon {
	db := connect()
	digimons := []Digimon{}
	db.Find(&digimons)
	for i := range digimons {
		json.Unmarshal([]byte(digimons[i].Data), &digimons[i].DataObj)
	}
	db.Close()
	return &digimons
}

// GetByID idでデジモンを取得
func (m DigimonRepository) GetByID(id int) *Digimon {
	db := connect()
	var ret *Digimon
	digimon := Digimon{}
	db.Find(&digimon, "id = ?", id)
	db.Close()
	if digimon.ID != 0 {
		json.Unmarshal([]byte(digimon.Data), &digimon.DataObj)
		ret = &digimon
	}
	return ret
}

// ErrPostRateLimited は全ユーザー共通の投稿制限を表す。
var ErrPostRateLimited = errors.New("posting is limited to once per minute")

// PostIntervalSeconds は全ユーザー共通の投稿間隔。
const PostIntervalSeconds = 60

// SetNewDigimon はDB上の投稿を確認し、全体で60秒に1件だけ作成する。
func (m DigimonRepository) SetNewDigimon(name string, data *[][]int) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	db := connect()
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := db.DB().Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// MySQLの名前付きロックは同じ接続で取得・解放する。
	// DBごとにロックを分け、複数プロセスからの同時投稿も直列化する。
	var locked sql.NullInt64
	if err := conn.QueryRowContext(ctx,
		"SELECT GET_LOCK(CONCAT('digimon-post:', SHA1(DATABASE())), 0)").Scan(&locked); err != nil {
		return err
	}
	if !locked.Valid {
		return errors.New("failed to acquire posting lock")
	}
	if locked.Int64 != 1 {
		return ErrPostRateLimited
	}
	defer func() {
		// リクエスト処理のタイムアウト後もロックを解放する。
		unlockCtx, unlockCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer unlockCancel()
		conn.ExecContext(unlockCtx, "SELECT RELEASE_LOCK(CONCAT('digimon-post:', SHA1(DATABASE())))")
	}()

	var recent bool
	if err := conn.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM digimons WHERE created_at > DATE_SUB(NOW(), INTERVAL ? SECOND))",
		PostIntervalSeconds).Scan(&recent); err != nil {
		return err
	}
	if recent {
		return ErrPostRateLimited
	}

	_, err = conn.ExecContext(ctx,
		"INSERT INTO digimons (name, data, created_at, updated_at) VALUES (?, ?, NOW(), NOW())",
		name, string(bytes))
	return err
}
