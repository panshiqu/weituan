package handler

import (
	"fmt"
	"net/http"

	"github.com/panshiqu/weituan/db"
	"github.com/panshiqu/weituan/define"
)

func stat(tn, fn string) (*define.StatInfo, error) {
	info := &define.StatInfo{}

	if err := db.MySQL.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tn)).Scan(&info.Total); err != nil {
		return nil, err
	}

	var n int

	for i := 0; i <= define.GC.StatEverydays; i++ {
		if err := db.MySQL.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s < SUBDATE(CURRENT_DATE,?) AND %s > SUBDATE(CURRENT_DATE,?)", tn, fn, fn), i-1, i).Scan(&n); err != nil {
			return nil, err
		}

		info.Everydays = append(info.Everydays, n)
	}

	for _, v := range define.GC.StatRecentdays {
		if err := db.MySQL.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s < CURRENT_DATE AND %s > SUBDATE(CURRENT_DATE,?)", tn, fn, fn), v).Scan(&n); err != nil {
			return nil, err
		}

		info.Recentdays = append(info.Recentdays, n)
	}

	return info, nil
}

func serveAdmin(w http.ResponseWriter, r *http.Request) error {
	return nil
}
