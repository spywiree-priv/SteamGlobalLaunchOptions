package main

import (
	"io"
	"os"

	"github.com/guregu/null/v5"
	"github.com/spywiree-priv/SteamGlobalLaunchOptions/vdf"
)

func Backup(r io.ReadSeeker, path string) error {
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, r)
	return err
}

func ApplyLaunchOptions(value, path string, overwrite bool) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := vdf.ParseText(f)
	if err != nil {
		return err
	}

	apps, err := data.GetChildByPath("Software", "Valve", "Steam", "apps")
	if err != nil {
		return err
	}

	for child := range apps.ChildrenIter() {
		if overwrite || !child.HasNonZeroChild("LaunchOptions") {
			child.SetChild(vdf.KeyValue{
				Key:   "LaunchOptions",
				Value: null.StringFrom(value),
			})
		}
	}

	if err = Backup(f, path+".bak"); err != nil {
		return err
	}

	if err = f.Truncate(0); err != nil {
		return err
	}

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	return data.Write(f)
}
