package handler

import (
	"fmt"
	"net/http"
	"time"

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
	s := [][]string{
		[]string{"user", "RegisterTime", "用户"},
		[]string{"sku", "PublishTime", "商品"},
		[]string{"share", "ShareTime", "分享"},
		[]string{"bargain", "BargainTime", "砍价"},
	}

	now := time.Now()

	for _, v := range s {
		info, err := stat(v[0], v[1])
		if err != nil {
			return err
		}

		fmt.Fprintf(w, "【%s】\n", v[2])
		fmt.Fprintf(w, "总量：%d\n", info.Total)

		var sum int
		for k, vv := range info.Everydays {
			switch k {
			case 0:
				fmt.Fprintf(w, "今日新增：%d\n", vv)
			case 1:
				sum += vv
				fmt.Fprintf(w, "昨日新增：%d\n", vv)
			default:
				sum += vv
				fmt.Fprintf(w, "%s新增：%d，近%d天总新增：%d\n", now.AddDate(0, 0, -k).Format("1-02"), vv, k, sum)
			}
		}

		for i := 0; i < len(define.GC.StatRecentdays); i++ {
			fmt.Fprintf(w, "最近%d天总新增：%d\n", define.GC.StatRecentdays[i], info.Recentdays[i])
		}

		fmt.Fprintln(w)
	}

	return nil
}
