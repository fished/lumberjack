package backend

import "time"

const (
	YearField = 1 << iota
	MonthField
	DayField
	TzField
)

// TimestampParser provides a tool for parsing a timestamp, hopefull in any format.
type TimestampParser struct {
	// Layout represents the layout string as used for time.ParseTimestamp
	Layout string
	// DefaultedFields represents the fields that are NOT returned by the Layout,
	// and should be set to match Time.now()
	DefaultedFields int
	// DefaultLocation is the timezone that should be used by default for this rule.
	DefaultLocation *time.Location
}

func (parser *TimestampParser) Parse(s string) (time.Time, error) {
	return parser.ParseAtTime(s, time.Now())
}

func (parser *TimestampParser) ParseAtTime(s string, t time.Time) (time.Time, error) {
	t, err := time.Parse(parser.Layout, s)
	if err != nil {
		return time.Now(), err
	}

	now := time.Now()

	year := t.Year()
	month := t.Month()
	day := t.Day()
	hour := t.Hour()
	min := t.Minute()
	sec := t.Second()
	ns := t.Nanosecond()
	loc := t.Location()

	if parser.DefaultedFields&TzField != 0 {
		loc = time.Local
	}

	if parser.DefaultedFields&DayField != 0 {
		if hour <= now.Hour() {
			day = now.Day()
		} else {
			day = now.AddDate(0, 0, -1).Day()
		}
	}

	if parser.DefaultedFields&MonthField != 0 {
		if day <= now.Day() {
			month = now.Month()
		} else {
			month = now.AddDate(0, -1, 0).Month()
		}
	}

	if parser.DefaultedFields&YearField != 0 {
		if month <= now.Month() {
			year = now.Year()
		} else {
			year = now.Year() - 1
		}
	}

	return time.Date(year, month, day, hour, min, sec, ns, loc), nil
}

var TimestampParserRegistry = map[string]TimestampParser{
	"ANSIC": TimestampParser{
		Layout:          "Mon Jan _2 15:04:05 2006",
		DefaultedFields: TzField,
		DefaultLocation: time.Local,
	},
	"UnixDate": TimestampParser{
		Layout:          "Mon Jan _2 15:04:05 MST 2006",
		DefaultedFields: 0,
		DefaultLocation: time.UTC,
	},
	"RubyDate": TimestampParser{
		Layout:          "Mon Jan 02 15:04:05 -0700 2006",
		DefaultedFields: 0,
		DefaultLocation: time.UTC,
	},
	"RFC822": TimestampParser{
		Layout:          "02 Jan 06 15:04 MST",
		DefaultedFields: 0,
		DefaultLocation: time.UTC,
	},
	"RFC822Z": TimestampParser{
		Layout:          "02 Jan 06 15:04 -0700", // RFC822 with numeric zone,
		DefaultedFields: 0,
		DefaultLocation: time.UTC,
	},
	"RFC850": TimestampParser{
		Layout:          "Monday, 02-Jan-06 15:04:05 MST",
		DefaultedFields: 0,
		DefaultLocation: time.UTC,
	},
	"RFC1123": TimestampParser{
		Layout:          "Mon, 02 Jan 2006 15:04:05 MST",
		DefaultedFields: 0,
		DefaultLocation: time.UTC,
	},
	"RFC1123Z": TimestampParser{
		Layout:          "Mon, 02 Jan 2006 15:04:05 -0700", // RFC1123 with numeric zone,
		DefaultedFields: 0,
		DefaultLocation: time.UTC,
	},
	"RFC3339": TimestampParser{
		Layout:          "2006-01-02T15:04:05Z07:00",
		DefaultedFields: 0,
		DefaultLocation: time.UTC,
	},
	"RFC3339Nano": TimestampParser{
		Layout:          "2006-01-02T15:04:05.999999999Z07:00",
		DefaultedFields: 0,
		DefaultLocation: time.UTC,
	},
	"Kitchen": TimestampParser{
		Layout:          "3:04PM",
		DefaultedFields: TzField | DayField | MonthField | YearField,
		DefaultLocation: time.Local,
	},
	"Stamp": TimestampParser{
		Layout:          "Jan _2 15:04:05",
		DefaultedFields: TzField | YearField,
		DefaultLocation: time.Local,
	},
	"StampMilli": TimestampParser{
		Layout:          "Jan _2 15:04:05.000",
		DefaultedFields: TzField | YearField,
		DefaultLocation: time.Local,
	},
	"StampMicro": TimestampParser{
		Layout:          "Jan _2 15:04:05.000000",
		DefaultedFields: TzField | YearField,
		DefaultLocation: time.Local,
	},
	"StampNano": TimestampParser{
		Layout:          "Jan _2 15:04:05.000000000",
		DefaultedFields: TzField | YearField,
		DefaultLocation: time.Local,
	},
}
