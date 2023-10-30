/**
 * @Time: 2023/10/30 18:48
 * @Author: jzechen
 * @File: sina_test.go
 * @Software: GoLand toresa
 */

package sina

import (
	"github.com/jzechen/toresa/pkg/manager/config"
	"github.com/jzechen/toresa/pkg/manager/contants"
	"testing"
)

const DrivePath = "/home/jze/go/src/toresa/browser"

func Test_getCookieStr(t *testing.T) {
	type args struct {
		dc       *config.DriveConfig
		userName string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// Add test cases.
		{
			name: "test one",
			args: args{
				dc: &config.DriveConfig{
					Type: contants.DefaultDriveType,
					Path: DrivePath, // input the path of your browse drive for the test
					Port: contants.DefaultDrivePort,
				},
				userName: "userName",
				password: "password",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCookieStr(tt.args.dc, tt.args.userName, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getCookieStr() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("cookie: %s", got)
		})
	}
}
