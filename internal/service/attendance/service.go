package attendance

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/itsaFan/fleetify-be/internal/appErr"
	"github.com/itsaFan/fleetify-be/internal/helper"
	"github.com/itsaFan/fleetify-be/internal/model"
	atdrepo "github.com/itsaFan/fleetify-be/internal/repo/attendance"
	emprepo "github.com/itsaFan/fleetify-be/internal/repo/employee"
	"gorm.io/gorm"
)

type service struct {
	atdRepo atdrepo.Repository
	empRepo emprepo.Repository
}

type Service interface {
	CreateEmpAttendance(ctx context.Context, employeeID string) (*model.Attendance, error)
	CloseEmpAttendance(ctx context.Context, employeeID string) (*model.Attendance, error)

	ListEmployeeAtdHistories(ctx context.Context, p ListInputEmp) (*AttendanceHistoryOutput, error)
	ListDeparmentAtdHistories(ctx context.Context, p ListInputDept) (*AttendanceHistoryOutput, error)
}

func New(atdRepo atdrepo.Repository, empRepo emprepo.Repository) Service {
	return &service{atdRepo: atdRepo, empRepo: empRepo}
}

func (s *service) CreateEmpAttendance(ctx context.Context, employeeID string) (*model.Attendance, error) {
	if employeeID == "" {
		return nil, fmt.Errorf("%w: employee_id is required", appErr.ErrRequiredField)
	}

	normalizedEmpId := helper.NormalizeStringField(employeeID)

	emp, err := s.empRepo.GetByEmployeeIDJoinDept(ctx, normalizedEmpId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee %q", appErr.ErrNotFound, normalizedEmpId)
		}
		return nil, err
	}

	_ = emp

	now := time.Now().UTC()
	attID := uuid.New().String()

	att := &model.Attendance{
		EmployeeID:   normalizedEmpId,
		AttendanceID: attID,
		ClockIn:      &now,
	}

	hist := &model.AttendanceHistory{
		EmployeeID:     normalizedEmpId,
		AttendanceID:   attID,
		DateAttendance: now,
		AttendanceType: 1,
		Description:    "Clock in",
	}

	if err := s.atdRepo.WithTx(ctx, func(tx atdrepo.Repository) error {
		open, err := tx.FindEmpOpenAttendanceForUpdate(ctx, normalizedEmpId)
		if err != nil {
			return err
		}
		if open != nil {
			return fmt.Errorf("%w: already clocked in with attendance_id=%s", appErr.ErrAlreadyExists, open.AttendanceID)
		}
		if err := tx.CreateEmpAttendanceByEmpId(ctx, att); err != nil {
			return err
		}
		if err := tx.CreateAttendanceHistory(ctx, hist); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return att, nil
}

func (s *service) CloseEmpAttendance(ctx context.Context, employeeID string) (*model.Attendance, error) {
	if employeeID == "" {
		return nil, fmt.Errorf("%w: employee_id is required", appErr.ErrRequiredField)
	}

	normalizedEmpId := helper.NormalizeStringField(employeeID)

	emp, err := s.empRepo.GetByEmployeeIDJoinDept(ctx, normalizedEmpId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: employee %q", appErr.ErrNotFound, normalizedEmpId)
		}
		return nil, err
	}

	_ = emp

	now := time.Now().UTC()
	var updated *model.Attendance

	if err := s.atdRepo.WithTx(ctx, func(tx atdrepo.Repository) error {
		open, err := tx.FindEmpOpenAttendanceForUpdate(ctx, normalizedEmpId)

		if err != nil {
			return err
		}

		if open == nil {
			return fmt.Errorf("%w: no open attendance for employee %q", appErr.ErrNotFound, normalizedEmpId)
		}

		if err := tx.UpdateAttendanceOutByAttendanceID(ctx, open.AttendanceID, now); err != nil {
			return err
		}

		hist := &model.AttendanceHistory{
			EmployeeID:     normalizedEmpId,
			AttendanceID:   open.AttendanceID,
			DateAttendance: now,
			AttendanceType: 2,
			Description:    "Clock out",
		}

		if err := tx.CreateAttendanceHistory(ctx, hist); err != nil {
			return err
		}

		open.ClockOut = &now
		updated = open
		return nil

	}); err != nil {
		return nil, err
	}

	return updated, nil

}

func (s *service) ListEmployeeAtdHistories(ctx context.Context, p ListInputEmp) (*AttendanceHistoryOutput, error) {
	empId := helper.NormalizeStringField(p.EmployeeID)

	if empId == "" {
		return nil, fmt.Errorf("%w: employee_id is required", appErr.ErrRequiredField)
	}

	emp, err := s.empRepo.GetByEmployeeIDJoinDept(ctx, empId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	loc := helper.LoadLocationOrUTC(p.TZ)
	if p.FromLocal == "" || p.ToLocal == "" {
		return nil, fmt.Errorf("%w: from/to are required (YYYY-MM-DD)", appErr.ErrRequiredField)
	}

	y1, m1, d1, err := func() (int, time.Month, int, error) {
		y, m, d, e := helper.ParseYYYYMMDD(p.FromLocal)
		return y, m, d, e
	}()

	if err != nil {
		return nil, fmt.Errorf("%w: invalid 'from' date", appErr.ErrInvalidInput)
	}

	y2, m2, d2, err := func() (int, time.Month, int, error) {
		y, m, d, e := helper.ParseYYYYMMDD(p.ToLocal)
		return y, m, d, e
	}()

	if err != nil {
		return nil, fmt.Errorf("%w: invalid 'to' date", appErr.ErrInvalidInput)
	}

	fromUTC, _ := helper.DayBoundsLocalToUTC(loc, y1, m1, d1)
	_, toUTC := helper.DayBoundsLocalToUTC(loc, y2, m2, d2)

	rows, err := s.atdRepo.ListHistoryByEmpId(ctx, atdrepo.ListParamsEmp{
		EmployeeID: empId,
		FromUtc:    fromUTC,
		ToUtc:      toUTC,
	})

	if err != nil {
		return nil, err
	}

	items := groupAndCompute(rows, loc, emp.Department.MaxClockInTime, emp.Department.MaxClockOutTime, emp.Name)

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].DateLocal == items[j].DateLocal {
			return items[i].EmployeeID < items[j].EmployeeID
		}
		return items[i].DateLocal < items[j].DateLocal
	})

	limit := p.Limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	page := p.Page
	if page <= 0 {
		page = 1
	}

	total := int64(len(items))
	start := min((page-1)*limit, len(items))
	end := min(start+limit, len(items))
	pageItems := items[start:end]

	return &AttendanceHistoryOutput{
		Items:     pageItems,
		Total:     total,
		FromLocal: p.FromLocal,
		ToLocal:   p.ToLocal,
		TZUsed:    p.TZ,
	}, nil
}

func (s *service) ListDeparmentAtdHistories(ctx context.Context, p ListInputDept) (*AttendanceHistoryOutput, error) {
	loc := helper.LoadLocationOrUTC(p.TZ)
	if p.FromLocal == "" || p.ToLocal == "" {
		return nil, fmt.Errorf("%w: from/to are required (YYYY-MM-DD)", appErr.ErrRequiredField)
	}

	y1, m1, d1, err := func() (int, time.Month, int, error) {
		y, m, d, e := helper.ParseYYYYMMDD(p.FromLocal)
		return y, m, d, e
	}()

	if err != nil {
		return nil, fmt.Errorf("%w: invalid 'from' date", appErr.ErrInvalidInput)
	}

	y2, m2, d2, err := func() (int, time.Month, int, error) {
		y, m, d, e := helper.ParseYYYYMMDD(p.ToLocal)
		return y, m, d, e
	}()

	if err != nil {
		return nil, fmt.Errorf("%w: invalid 'to' date", appErr.ErrInvalidInput)
	}

	fromUTC, _ := helper.DayBoundsLocalToUTC(loc, y1, m1, d1)
	_, toUTC := helper.DayBoundsLocalToUTC(loc, y2, m2, d2)

	rows, err := s.atdRepo.ListHistoryByDepartment(ctx, atdrepo.ListParamsDept{
		DepartmentID: p.DepartmentID,
		FromUtc:      fromUTC,
		ToUtc:        toUTC,
	})

	if err != nil {
		return nil, err
	}

	type deptTimes struct {
		in, out string
		name    string
	}
	cache := map[string]deptTimes{}

	getDeptTimes := func(empId string) (deptTimes, error) {
		if v, ok := cache[empId]; ok {
			return v, nil
		}
		emp, err := s.empRepo.GetByEmployeeIDJoinDept(ctx, empId)
		if err != nil {
			return deptTimes{}, err
		}
		v := deptTimes{in: emp.Department.MaxClockInTime, out: emp.Department.MaxClockOutTime, name: emp.Name}
		cache[empId] = v
		return v, nil
	}

	grouped := map[string][]model.AttendanceHistory{}
	for _, r := range rows {
		grouped[r.EmployeeID] = append(grouped[r.EmployeeID], r)
	}

	all := make([]AttendanceHistoryItem, 0, len(rows))
	for empId, hist := range grouped {
		sort.SliceStable(hist, func(i, j int) bool {
			if hist[i].DateAttendance.Equal(hist[j].DateAttendance) {
				return hist[i].ID < hist[j].ID
			}
			return hist[i].DateAttendance.Before(hist[j].DateAttendance)
		})

		type dayAgg struct {
			firstInUTC *time.Time
			lastOutUTC *time.Time
			attendance string
		}
		byDay := map[string]*dayAgg{}

		for _, h := range hist {
			local := h.DateAttendance.In(loc)
			key := fmt.Sprintf("%04d-%02d-%02d", local.Year(), local.Month(), local.Day())
			if _, ok := byDay[key]; !ok {
				byDay[key] = &dayAgg{}
			}
			agg := byDay[key]
			switch h.AttendanceType {
			case 1:
				if agg.firstInUTC == nil || h.DateAttendance.Before(*agg.firstInUTC) {
					t := h.DateAttendance
					agg.firstInUTC = &t
					agg.attendance = h.AttendanceID
				}

			case 2:
				if agg.lastOutUTC == nil || h.DateAttendance.After(*agg.lastOutUTC) {
					t := h.DateAttendance
					agg.lastOutUTC = &t
				}
			}
		}

		dt, err := getDeptTimes(empId)
		if err != nil {
			continue
		}

		hIn, mIn, sIn, _ := helper.ParseCutoffHHMMSS(dt.in)
		hOut, mOut, sOut, _ := helper.ParseCutoffHHMMSS(dt.out)

		for day, agg := range byDay {
			item := AttendanceHistoryItem{
				EmployeeID:   empId,
				EmployeeName: dt.name,
				DateLocal:    day,
				StatusIn:     "missing_in",
				StatusOut:    "no_out",
			}

			y, _m, _d, _ := helper.ParseYYYYMMDD(day)
			deadlineInLocal := time.Date(y, time.Month(_m), _d, hIn, mIn, sIn, 0, loc)
			deadlineOutLocal := time.Date(y, time.Month(_m), _d, hOut, mOut, sOut, 0, loc)

			// Attendance IN
			if agg.firstInUTC != nil {
				local := agg.firstInUTC.In(loc)
				str := local.Format("15:04:05")
				item.ClockInLocal = &str
				item.ClockInUTC = agg.firstInUTC

				diffMin := signedCeilMinutes(local.Sub(deadlineInLocal))
				if diffMin == 0 {
					item.StatusIn = "on_time"
					z := 0
					item.DeltaInMinutes = &z
				} else if diffMin > 0 {
					item.StatusIn = "late"
					item.DeltaInMinutes = &diffMin
				} else {
					item.StatusIn = "early"
					item.DeltaInMinutes = &diffMin
				}
				item.AttendanceID = agg.attendance
			}

			// Attendance OUT
			if agg.lastOutUTC != nil {
				local := agg.lastOutUTC.In(loc)
				str := local.Format("15:04:05")
				item.ClockOutLocal = &str
				item.ClockOutUTC = agg.lastOutUTC

				diffMin := signedCeilMinutes(local.Sub(deadlineOutLocal))

				switch {
				case diffMin == 0:
					item.StatusOut = "normal"
					z := 0
					item.DeltaOutMinutes = &z
				case diffMin > 0:
					item.StatusOut = "overtime"
					item.DeltaOutMinutes = &diffMin
				default:
					item.StatusOut = "early_leave"
					item.DeltaOutMinutes = &diffMin
				}
			}
			all = append(all, item)
		}

	}

	sort.SliceStable(all, func(i, j int) bool {
		if all[i].DateLocal == all[j].DateLocal {
			return all[i].EmployeeID < all[j].EmployeeID
		}
		return all[i].DateLocal < all[j].DateLocal
	})

	limit := p.Limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	page := p.Page
	if page <= 0 {
		page = 1
	}

	total := int64(len(all))
	start := min((page-1)*limit, len(all))
	end := min(start+limit, len(all))
	pageItems := all[start:end]

	return &AttendanceHistoryOutput{
		Items:     pageItems,
		Total:     total,
		FromLocal: p.FromLocal,
		ToLocal:   p.ToLocal,
		TZUsed:    p.TZ,
	}, nil
}

func groupAndCompute(
	rows []model.AttendanceHistory,
	loc *time.Location,
	maxIn string,
	maxOut string,
	employeeName string,
) []AttendanceHistoryItem {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].DateAttendance.Equal(rows[j].DateAttendance) {
			return rows[i].ID < rows[j].ID
		}
		return rows[i].DateAttendance.Before(rows[j].DateAttendance)
	})

	type dayAgg struct {
		firstInUTC *time.Time
		lastOutUTC *time.Time
		attID      string
	}
	byDay := map[string]*dayAgg{}
	var eid string
	if len(rows) > 0 {
		eid = rows[0].EmployeeID
	}

	for _, r := range rows {
		if eid == "" {
			eid = r.EmployeeID
		}
		local := r.DateAttendance.In(loc)
		key := fmt.Sprintf("%04d-%02d-%02d", local.Year(), local.Month(), local.Day())
		if _, ok := byDay[key]; !ok {
			byDay[key] = &dayAgg{}
		}
		agg := byDay[key]

		switch r.AttendanceType {
		case 1:
			if agg.firstInUTC == nil || r.DateAttendance.Before(*agg.firstInUTC) {
				t := r.DateAttendance
				agg.firstInUTC = &t
				agg.attID = r.AttendanceID
			}
		case 2:
			if agg.lastOutUTC == nil || r.DateAttendance.After(*agg.lastOutUTC) {
				t := r.DateAttendance
				agg.lastOutUTC = &t
			}
		}
	}

	hIn, mIn, sIn, _ := helper.ParseCutoffHHMMSS(maxIn)
	hOut, mOut, sOut, _ := helper.ParseCutoffHHMMSS(maxOut)

	items := make([]AttendanceHistoryItem, 0, len(byDay))
	for day, agg := range byDay {
		item := AttendanceHistoryItem{
			EmployeeID:   eid,
			EmployeeName: employeeName,
			DateLocal:    day,
			StatusIn:     "missing_in",
			StatusOut:    "no_out",
		}

		y, _m, _d, _ := helper.ParseYYYYMMDD(day)
		deadlineInLocal := time.Date(y, time.Month(_m), _d, hIn, mIn, sIn, 0, loc)
		deadlineOutLocal := time.Date(y, time.Month(_m), _d, hOut, mOut, sOut, 0, loc)

		if agg.firstInUTC != nil {
			local := agg.firstInUTC.In(loc)
			str := local.Format("15:04:05")
			item.ClockInLocal = &str
			item.ClockInUTC = agg.firstInUTC

			diffMin := signedCeilMinutes(local.Sub(deadlineInLocal))
			if diffMin == 0 {
				item.StatusIn = "on_time"
				z := 0
				item.DeltaInMinutes = &z
			} else if diffMin > 0 {
				item.StatusIn = "late"
				item.DeltaInMinutes = &diffMin
			} else {
				item.StatusIn = "early"
				item.DeltaInMinutes = &diffMin
			}
			item.AttendanceID = agg.attID
		}

		if agg.lastOutUTC != nil {
			local := agg.lastOutUTC.In(loc)
			str := local.Format("15:04:05")
			item.ClockOutLocal = &str
			item.ClockOutUTC = agg.lastOutUTC

			diffMin := signedCeilMinutes(local.Sub(deadlineOutLocal))

			switch {
			case diffMin == 0:
				item.StatusOut = "normal"
				z := 0
				item.DeltaOutMinutes = &z
			case diffMin > 0:
				item.StatusOut = "overtime"
				item.DeltaOutMinutes = &diffMin
			default:
				item.StatusOut = "early_leave"
				item.DeltaOutMinutes = &diffMin
			}
		}
		items = append(items, item)
	}

	sort.SliceStable(items, func(i, j int) bool { return items[i].DateLocal < items[j].DateLocal })
	return items
}

func signedCeilMinutes(d time.Duration) int {
	secs := d.Seconds()
	if secs >= 0 {
		return int(math.Ceil(secs / 60.0))
	}
	return -int(math.Ceil(math.Abs(secs) / 60.0))
}
