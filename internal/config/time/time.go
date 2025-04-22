package time

import "time"

func SetServerTimeZone() {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	time.Local = loc
}
