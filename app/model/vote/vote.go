package vote

import (
	"database/sql"
	"fmt"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/lib/pq"
	"log"
)

type Vote struct {
	Id          int    `db:"id" json:"id,omitempty"`
	BeItmo      string `db:"be_itmo" json:"beItmo,omitempty"`
	Title       string `db:"title" json:"title,omitempty"`
	Vote1       string `db:"vote_1" json:"vote1,omitempty"`
	Vote1Value1 string `db:"vote_1_value_1" json:"vote1Value1,omitempty"`
	Vote1Value2 string `db:"vote_1_value_2" json:"vote1Value2,omitempty"`
	Vote2       string `db:"vote_2" json:"vote2,omitempty"`
	Vote2Value1 string `db:"vote_2_value_1" json:"vote2Value1,omitempty"`
	Vote2Value2 string `db:"vote_2_value_2" json:"vote2Value2,omitempty"`
	Vote2Value3 string `db:"vote_2_value_3" json:"vote2Value3,omitempty"`
	Status      bool   `json:"status,omitempty"`
}

func GetVoteByTag(tag string, db *sql.DB) Vote {
	row, err := db.Query("SELECT * FROM vote WHERE be_itmo=$1", tag)
	if err != nil {
		log.Println(err)
		db.Close()
		return Vote{}
	}
	var v Vote
	err = sqlscan.ScanOne(&v, row)
	if err != nil {
		log.Println(err)
		db.Close()
		return Vote{}
	}

	return v
}

func (v Vote) CheckUserVote(userId int, db *sql.DB) bool {

	row := db.QueryRow("SELECT * FROM vote_log WHERE user_id=$1 and vote_id=$2 and vote_date = current_date", userId, v.Id)
	var r interface{}
	err := row.Scan(&r)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {

	}

	return true
}

func UpdateBeItmo(tag string, userId int, value int, db *sql.DB) error {
	quoted := pq.QuoteIdentifier(tag)

	_, err := db.Exec(fmt.Sprintf("UPDATE users SET %s=%s+$1 WHERE id=$2", quoted, quoted), value, userId)
	if err != nil {
		return err
	}
	return nil
}
