package db

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestDb(t *testing.T) {
	q, err := DB()
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()

	sample := []*Dongle{
		{
			IMEI: "000000",
			IMSI: "000001",
			Path: "/dev/ttyUSB5",
			TTY:  5,
		},
		{
			IMEI: "000001",
			IMSI: "000002",
			Path: "/dev/tty6",
			TTY:  6,
		},
		{
			IMEI: "000002",
			IMSI: "000003",
			Path: "/dev/tty7",
			TTY:  7,
		},
	}

	for _, v := range sample {
		err = CreateDongle(q, v)
		if err != nil {
			t.Error(err)
		}
	}

	for _, v := range sample {
		d, err := GetDongle(q, v.Path)
		if err != nil {
			t.Error(err)
		}
		if d != nil {
			if d.IMEI != v.IMEI {
				t.Errorf("expected %s got %s", v.IMEI, d.IMEI)
			}
		}
	}
	a, err := GetAllDongles(q)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != len(sample) {
		t.Errorf("expected %d got %d", len(sample), len(a))
	}
	d, err := GetDistinct(q)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(d))

	err = RemoveDongle(q, a[0])
	if err != nil {
		t.Error(err)
	}
	_, err = GetDongle(q, a[0].Path)
	if err != sql.ErrNoRows {
		t.Error("expected %v got %v", sql.ErrNoRows, err)
	}

	sample[0], sample[1] = sample[1], sample[0]
	imei := sample[0].IMEI
	for i := range sample {
		sample[i].IMEI = imei
	}
	q.Close()
	qq, err := DB()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range sample {
		err = CreateDongle(qq, v)
		if err != nil {
			t.Error(err)
		}
	}
	low, err := GetSymlinkCandidate(qq, imei)
	if err != nil {
		t.Fatal(err)
	}
	expect := 5
	if low.TTY != 5 {
		t.Errorf("expected %s got %s", expect, low.TTY)
	}
}

func TestRemoveDongle(t *testing.T) {
	q, err := dbWIthName("test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer q.Close()
	imei := "123456"

	for v := range make([]struct{}, 10) {
		err = CreateDongle(q, &Dongle{IMEI: imei, Path: fmt.Sprint(v)})
		if err != nil {
			t.Fatal(err)
		}
	}
	dongels, err := GetAllDongles(q)
	if err != nil {
		t.Fatal(err)
	}
	if len(dongels) != 10 {
		t.Errorf("expected 10 got %d", len(dongels))
	}
	err = RemoveDongle(q, &Dongle{IMEI: imei})
	if err != nil {
		t.Fatal(err)
	}
	dongels, err = GetAllDongles(q)
	if err != nil {
		t.Fatal(err)
	}
	if len(dongels) != 0 {
		t.Errorf("expected 10 got %d", len(dongels))
	}

}
