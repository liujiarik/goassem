package run

import "testing"

func TestCheckMatch(t *testing.T) {

	t.Log(CheckMatch("sdsd.sh", "*"))
	t.Log(CheckMatch("sdsd.sh", "sdd"))
	t.Log(CheckMatch("sdsd.sh", "*.sh"))
	t.Log(CheckMatch("sdsd.sh ", "sdsd.sh"))
	t.Log(CheckMatch("sdsd.sh ", "sd*"))
	t.Log(CheckMatch("sdsd.sh", "sd*DF"))

}
