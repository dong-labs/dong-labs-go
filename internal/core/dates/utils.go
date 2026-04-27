// Package dates 提供日期处理工具函数
package dates

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// WeekStart 一周开始日的定义
type WeekStart string

const (
	WeekStartMonday    WeekStart = "monday"
	WeekStartSunday    WeekStart = "sunday"
	WeekStartSaturday  WeekStart = "saturday"
)

// Today 获取今天的日期
func Today() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}

// Yesterday 获取昨天的日期
func Yesterday() time.Time {
	return Today().AddDate(0, 0, -1)
}

// Tomorrow 获取明天的日期
func Tomorrow() time.Time {
	return Today().AddDate(0, 0, 1)
}

// ThisWeek 获取本周的日期范围
// 返回 (周开始日期, 周结束日期)
func ThisWeek(weekStart WeekStart) (time.Time, time.Time) {
	today := Today()
	weekday := today.Weekday() // Sunday = 0, Monday = 1, ..., Saturday = 6

	var daysToSubtract int
	switch weekStart {
	case WeekStartMonday:
		// Monday = 1, 需要转换为 0-based
		daysToSubtract = int(weekday) - 1
		if daysToSubtract < 0 {
			daysToSubtract = 6 // Sunday
		}
	case WeekStartSunday:
		daysToSubtract = int(weekday)
	case WeekStartSaturday:
		daysToSubtract = (int(weekday) + 1) % 7
	default:
		daysToSubtract = int(weekday) - 1
		if daysToSubtract < 0 {
			daysToSubtract = 6
		}
	}

	start := today.AddDate(0, 0, -daysToSubtract)
	end := start.AddDate(0, 0, 6)
	return start, end
}

// LastWeek 获取上周的日期范围
func LastWeek(weekStart WeekStart) (time.Time, time.Time) {
	start, end := ThisWeek(weekStart)
	lastWeekStart := start.AddDate(0, 0, -7)
	lastWeekEnd := end.AddDate(0, 0, -7)
	return lastWeekStart, lastWeekEnd
}

// ThisMonth 获取本月的日期范围
// 返回 (月初日期, 月末日期)
func ThisMonth() (time.Time, time.Time) {
	today := Today()
	firstDay := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	// 获取下个月第一天，再减一天得到本月最后一天
	nextMonth := firstDay.AddDate(0, 1, 0)
	lastDay := nextMonth.AddDate(0, 0, -1)

	return firstDay, lastDay
}

// LastMonth 获取上月的日期范围
func LastMonth() (time.Time, time.Time) {
	firstDay, _ := ThisMonth()

	// 上月最后日 = 本月第一日 - 1
	lastMonthLast := firstDay.AddDate(0, 0, -1)

	// 上月第一日
	lastMonthFirst := time.Date(lastMonthLast.Year(), lastMonthLast.Month(), 1, 0, 0, 0, 0, lastMonthLast.Location())

	return lastMonthFirst, lastMonthLast
}

// ThisQuarter 获取本季度的日期范围
func ThisQuarter() (time.Time, time.Time) {
	today := Today()
	quarter := (int(today.Month()) - 1) / 3
	quarterStartMonth := time.Month(quarter*3 + 1)
	quarterEndMonth := time.Month(quarter*3 + 3)

	start := time.Date(today.Year(), quarterStartMonth, 1, 0, 0, 0, 0, today.Location())

	// 计算季度最后一天
	if quarterEndMonth == 12 {
		return start, time.Date(today.Year(), 12, 31, 0, 0, 0, 0, today.Location())
	}
	nextMonth := time.Date(today.Year(), quarterEndMonth+1, 1, 0, 0, 0, 0, today.Location())
	end := nextMonth.AddDate(0, 0, -1)

	return start, end
}

// ThisYear 获取本年的日期范围
func ThisYear() (time.Time, time.Time) {
	today := Today()
	start := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, today.Location())
	end := time.Date(today.Year(), 12, 31, 0, 0, 0, 0, today.Location())
	return start, end
}

// Parse 解析日期字符串
// 支持格式: YYYY-MM-DD, YYYY/MM/DD, YYYYMMDD
// 关键词: today, yesterday, tomorrow
func Parse(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(strings.ToLower(dateStr))

	// 关键词
	switch dateStr {
	case "today":
		return Today(), nil
	case "yesterday":
		return Yesterday(), nil
	case "tomorrow":
		return Tomorrow(), nil
	case "this week":
		start, _ := ThisWeek(WeekStartMonday)
		return start, nil
	case "this month":
		start, _ := ThisMonth()
		return start, nil
	case "this year":
		start, _ := ThisYear()
		return start, nil
	}

	// 标准格式
	layouts := []string{
		"2006-01-02",
		"2006/01/02",
		"20060102",
		"2006-1-2",
		"2006/1/2",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析日期: %s", dateStr)
}

// Format 格式化日期为字符串
func Format(t time.Time, fmtStr string) string {
	if fmtStr == "" {
		fmtStr = "2006-01-02"
	}
	return t.Format(fmtStr)
}

// DateRange 生成日期范围内的所有日期
// inclusive: 是否包含结束日期
func DateRange(start, end time.Time, inclusive bool) []time.Time {
	if inclusive {
		end = end.Add(24 * time.Hour)
	}

	days := int(end.Sub(start).Hours() / 24)
	if days < 0 {
		return []time.Time{}
	}

	result := make([]time.Time, days)
	for i := 0; i < days; i++ {
		result[i] = start.AddDate(0, 0, i)
	}
	return result
}

// RelativeDate 计算相对日期
func RelativeDate(base time.Time, years, months, weeks, days int) time.Time {
	if base.IsZero() {
		base = Today()
	}

	// 处理年月
	if years != 0 || months != 0 {
		totalMonths := base.Year()*12 + int(base.Month()) - 1
		totalMonths += years*12 + months

		newYear := totalMonths / 12
		newMonth := totalMonths%12 + 1

		// 处理日期溢出
		maxDay := daysInMonth(newYear, newMonth)
		newDay := base.Day()
		if newDay > maxDay {
			newDay = maxDay
		}

		base = time.Date(newYear, time.Month(newMonth), newDay, 0, 0, 0, 0, base.Location())
	}

	// 处理周和日
	return base.AddDate(0, 0, weeks*7+days)
}

// daysInMonth 获取指定月份的天数
func daysInMonth(year, month int) int {
	if month == 12 {
		t := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
		return t.AddDate(0, 0, -1).Day()
	}
	t := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	return t.AddDate(0, 0, -1).Day()
}

// ParseDateRange 解析日期范围字符串
// 格式: "2024-01-01:2024-01-31" 或 "this month"
func ParseDateRange(rangeStr string) (start, end time.Time, err error) {
	rangeStr = strings.TrimSpace(rangeStr)

	// 预定义范围
	switch rangeStr {
	case "today":
		t := Today()
		return t, t, nil
	case "yesterday":
		t := Yesterday()
		return t, t, nil
	case "this week":
		start, end := ThisWeek(WeekStartMonday)
		return start, end, nil
	case "last week":
		start, end := LastWeek(WeekStartMonday)
		return start, end, nil
	case "this month":
		start, end := ThisMonth()
		return start, end, nil
	case "last month":
		start, end := LastMonth()
		return start, end, nil
	case "this quarter":
		start, end := ThisQuarter()
		return start, end, nil
	case "this year":
		start, end := ThisYear()
		return start, end, nil
	}

	// 解析 "start:end" 格式
	if strings.Contains(rangeStr, ":") {
		parts := strings.SplitN(rangeStr, ":", 2)
		start, err = Parse(parts[0])
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		end, err = Parse(parts[1])
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		return start, end, nil
	}

	// 单个日期
	t, err := Parse(rangeStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return t, t, nil
}

// ToISODate 转换为 ISO 日期字符串 (YYYY-MM-DD)
func ToISODate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FromISODate 从 ISO 日期字符串解析
func FromISODate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// GetQuarter 获取日期所属季度 (1-4)
func GetQuarter(t time.Time) int {
	return (int(t.Month())-1)/3 + 1
}

// GetWeekOfYear 获取周数
func GetWeekOfYear(t time.Time) (year, week int) {
	return t.ISOWeek()
}

// Age 计算年龄（基于年数）
func Age(birthday time.Time) int {
	t := Today()
	years := t.Year() - birthday.Year()
	// 如果还没到生日，减一年
	if t.Month() < birthday.Month() || (t.Month() == birthday.Month() && t.Day() < birthday.Day()) {
		years--
	}
	return years
}

// DaysBetween 计算两个日期之间的天数
func DaysBetween(start, end time.Time) int {
	duration := end.Sub(start)
	return int(duration.Hours() / 24)
}

// AddWorkdays 添加工作日（跳过周末）
func AddWorkdays(start time.Time, days int) time.Time {
	result := start
	added := 0
	for added < days {
		result = result.AddDate(0, 0, 1)
		if result.Weekday() != time.Saturday && result.Weekday() != time.Sunday {
			added++
		}
	}
	return result
}

// IsWeekend 判断是否是周末
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// ParseYearMonth 解析年月字符串 (YYYY-MM)
func ParseYearMonth(s string) (time.Time, error) {
	return time.Parse("2006-01", s)
}

// YearMonthRange 返回年月的日期范围
func YearMonthRange(year int, month time.Month) (start, end time.Time) {
	start = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	nextMonth := start.AddDate(0, 1, 0)
	end = nextMonth.AddDate(0, 0, -1)
	return start, end
}

// NowString 返回当前时间的 ISO 格式字符串
func NowString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// UnixMillis 返回毫秒级时间戳
func UnixMillis(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond()/1000000)
}

// FromUnixMillis 从毫秒级时间戳解析
func FromUnixMillis(ms int64) time.Time {
	sec := ms / 1000
	nsec := (ms % 1000) * 1000000
	return time.Unix(sec, nsec)
}

// QuarterStartEnd 返回指定季度的起止日期
func QuarterStartEnd(year int, quarter int) (start, end time.Time) {
	startMonth := time.Month((quarter-1)*3 + 1)
	endMonth := time.Month(quarter * 3)

	start = time.Date(year, startMonth, 1, 0, 0, 0, 0, time.UTC)

	if endMonth == 12 {
		end = time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	} else {
		nextMonth := time.Date(year, endMonth+1, 1, 0, 0, 0, 0, time.UTC)
		end = nextMonth.Add(-time.Second)
	}

	return start, end
}

// FormatDuration 格式化时间间隔为人类可读格式
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f秒", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0f分钟", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.0f小时", d.Hours())
	}
	return fmt.Sprintf("%.0f天", d.Hours()/24)
}

// ParseNaturalLanguage 解析自然语言日期
// 支持: "3天前", "2周后", "1个月前", etc.
func ParseNaturalLanguage(s string) (time.Time, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	// 提取数字和单位
	re := strings.NewReader(s)
	var num int
	_, err := fmt.Fscanf(re, "%d", &num)
	if err != nil {
		return time.Time{}, fmt.Errorf("无法解析数字: %w", err)
	}

	// 剩余部分是单位
	unit := strings.TrimPrefix(s, strconv.Itoa(num))
	unit = strings.TrimSpace(unit)

	var years, months, weeks, days int
	switch unit {
	case "年", "year", "years", "y":
		years = num
	case "月", "month", "months", "M":
		months = num
	case "周", "week", "weeks", "w":
		weeks = num
	case "天", "day", "days", "d":
		days = num
	default:
		return time.Time{}, fmt.Errorf("未知的时间单位: %s", unit)
	}

	// 检查是"前"还是"后"
	if strings.Contains(s, "前") || strings.Contains(s, "ago") {
		years, months, weeks, days = -years, -months, -weeks, -days
	}

	return RelativeDate(time.Time{}, years, months, weeks, days), nil
}
