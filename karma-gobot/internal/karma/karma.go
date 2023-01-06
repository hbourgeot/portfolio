package karma

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Karma struct {
	User        string
	Count       int
	LastUpdated time.Time
}

type KarmaModel struct {
	DB *sql.DB
}

// CREATE & INSERT METHODS

func (m *KarmaModel) CreateTable(channel string) error {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(username VARCHAR(200) NOT NULL PRIMARY KEY, karma INT NOT NULL, last_updated DATETIME NOT NULL)", channel)
	if _, err := m.DB.Exec(query); err != nil {
		return err
	}

	infoLog.Println("Created table for: ", channel, "group")

	return nil
}

func (m KarmaModel) InsertUsers(username, channel string) error {
	query := fmt.Sprintf("INSERT INTO `%s`(username, karma, last_updated) VALUES (?, ?, ?)", channel)
	if _, err := m.DB.Exec(query, username, 0, time.Now()); err != nil {
		return err
	}
	return nil
}

// GET methods

func (m *KarmaModel) GetKarmas(channel string, top bool) ([]*Karma, error) {
	var query string
	if top {
		query = fmt.Sprintf("SELECT username, karma FROM `%s` ORDER BY ASC TOP 10", channel)
	} else {
		query = fmt.Sprintf("SELECT username, karma FROM `%s` ORDER BY DESC TOP 10", channel)
	}

	rows, err := m.DB.Query(query)
	defer rows.Close()

	var karmas []*Karma

	for rows.Next() {
		var karma *Karma
		if err = rows.Scan(&karma.User, &karma.Count); err != nil {
			if err == sql.ErrNoRows {
				return nil, err
			}

			return nil, err
		}

		karmas = append(karmas, karma)
	}
	return karmas, nil
}

func (m *KarmaModel) GetActualKarma(username, channel string) (int, bool, error) {
	var user Karma
	query := fmt.Sprintf("SELECT karma FROM `%s` WHERE username = ?", channel)

	if err := m.DB.QueryRow(query, username).Scan(&user.Count); err != nil {
		if err == sql.ErrNoRows {
			return 0, true, err
		}
		return 0, true, err // return true if there is no rows
	}

	return user.Count, false, nil
}

func (m *KarmaModel) GetLastUpdated(username, channel string) (time.Time, bool) {
	var user Karma
	if channel == "" {
		return user.LastUpdated, false
	}

	query := fmt.Sprintf("SELECT last_updated FROM `%s` WHERE username = ?", channel)

	if err := m.DB.QueryRow(query, username).Scan(&user.LastUpdated); err != nil {
		if err == sql.ErrNoRows {
			return user.LastUpdated, true
		}
		return user.LastUpdated, true
	}

	return user.LastUpdated, false
}

// UPDATE METHODS

func (m *KarmaModel) AddKarma(karmaTransmitter, karmaReceiver, channel string) error {
	query := fmt.Sprintf("UPDATE `%s` SET karma = ? WHERE username = ?", channel)

	karma, noRows, err := m.GetActualKarma(karmaReceiver, channel)
	if noRows {
		err := m.InsertUsers(karmaReceiver, channel)
		if err != nil {
			return err
		}
		karma = 0
	} else {
		return err
	}

	karma++
	_, err = m.DB.Exec(query, karma, karmaReceiver)
	if err != nil {
		return err
	}

	err = m.updateLastKarma(time.Now(), channel, karmaTransmitter)
	if err != nil {
		return err
	}

	return nil
}

func (m *KarmaModel) SubstractKarma(karmaTransmitter, karmaReceiver, channel string) error {
	query := fmt.Sprintf("UPDATE `%s` SET Karma = ? WHERE username = ?", channel)

	karma, noRows, err := m.GetActualKarma(karmaReceiver, channel)
	if noRows {
		err := m.InsertUsers(karmaReceiver, channel)
		if err != nil {
			return err
		}
		karma = 0
	} else {
		return err
	}

	karma--
	_, err = m.DB.Exec(query, karma, karmaReceiver)
	if err != nil {
		return err
	}

	err = m.updateLastKarma(time.Now(), channel, karmaTransmitter)
	if err != nil {
		return err
	}

	return nil
}

func (m *KarmaModel) updateLastKarma(date time.Time, channel, username string) error {
	query := fmt.Sprintf("UPDATE `%s` SET last_updated = ? WHERE username = ?", channel)
	_, err := m.DB.Exec(query, date, username)
	if err != nil {
		return err
	}

	return nil
}
