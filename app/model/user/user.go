package user

import (
	"database/sql"
	"github.com/georgysavva/scany/sqlscan"
	"log"
)

type User struct {
	Id        int    `db:"id"`
	Name      string `db:"name"`
	Education string `db:"education"`
	Position  string `db:"position"`
	Authority string `db:"authority"`
	Friendly  int8   `db:"friendly"`
	Fit       int8   `db:"fit"`
	Healthy   int8   `db:"healthy"`
	Pro       int8   `db:"pro"`
	Open      int8   `db:"open"`
	Eco       int8   `db:"eco"`
}

type Event struct {
	Id      int    `db:"id"`
	Rank    string `db:"rank"`
	Roles   string `db:"roles"`
	BeItmo  string `db:"be_itmo"`
	UserId  int    `db:"user_id"`
	EventId int    `db:"event_id"`
}

func (u User) Get(id int, db *sql.DB) (User, error) {
	row, err := db.Query("SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		log.Println(err)
		db.Close()
		return User{}, err
	}
	err = sqlscan.ScanOne(&u, row)
	if err != nil {
		log.Println(err)
		db.Close()
		return User{}, err
	}
	return u, nil
}

func (u User) GetBeMap() map[string]int8 {
	beMap := make(map[string]int8)
	beMap["friendly"] = u.Friendly
	beMap["fit"] = u.Fit
	beMap["healthy"] = u.Healthy
	beMap["pro"] = u.Pro
	beMap["open"] = u.Open
	beMap["eco"] = u.Eco

	return beMap
}

func (u User) GetLowestPoint() (string, int8) {

	beMap := u.GetBeMap()

	lowestValue := u.Friendly
	lowestName := "friendly"

	for index, value := range beMap {
		if value < lowestValue {
			lowestValue = value
			lowestName = index
		}
	}

	return lowestName, lowestValue
}

func (u User) GetAvgValue() int8 {
	beMap := u.GetBeMap()
	score := int8(0)
	i := int8(0)
	for _, value := range beMap {
		i++
		score += value
	}
	return score / i

}

func (u User) GetUserEventsByTag(tag string, db *sql.DB) ([]Event, error) {
	rows, err := db.Query("SELECT * FROM user_events WHERE be_itmo=$1 and user_id=$2", tag, u.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var events []Event
	err = sqlscan.ScanAll(&events, rows)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return events, nil
}

func (u User) GetCatImageSrc(lowestValue int8, db *sql.DB) (string, int) {
	id := 0
	switch {
	case lowestValue > 90:
		id = 1
	case lowestValue > 60:
		id = 2
	case lowestValue > 20:
		id = 5
	case lowestValue > 40:
		id = 6
	case lowestValue > 30:
		id = 7
	case lowestValue > 75:
		id = 8
	default:
		id = 9
	}

	type Cat struct {
		Id  int    `db:"id"`
		Src string `db:"src"`
	}
	var catDb Cat
	row := db.QueryRow("SELECT * FROM cat_image WHERE id=$1", id)
	row.Scan(&catDb)

	return catDb.Src, id
}
