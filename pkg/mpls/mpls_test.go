package mpls

import (
	"reflect"
	"testing"
)

func TestParseMPLS(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name              string
		args              args
		wantHeader        *MPLSHeader
		wantAppinfo       *AppInfo
		wantPlaylist      *PlayList
		wantChapterMarks  *PlaylistMarks
		wantExtensiondata *Extensions
		wantErr           bool
	}{
		{
			name:       "valid MPLS file",
			args:       args{filePath: "testdata/valid.mpls"}, // Place a test .mpls file in pkg/mpls/testdata/
			wantHeader: &MPLSHeader{
				// Fill with expected values, e.g., TypeIndicator: "MPLS0200"
			},
			wantAppinfo: &AppInfo{
				// Fill with expected values
			},
			wantPlaylist: &PlayList{
				// Fill with expected values
			},
			wantChapterMarks: &PlaylistMarks{
				// Fill with expected values
			},
			wantExtensiondata: &Extensions{
				// Fill with expected values
			},
			wantErr: false,
		},
		{
			name:    "invalid MPLS file",
			args:    args{filePath: "testdata/invalid.mpls"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHeader, gotAppinfo, gotPlaylist, gotChapterMarks, gotExtensiondata, err := ParseMPLS(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMPLS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHeader, tt.wantHeader) {
				t.Errorf("ParseMPLS() gotHeader = %v, want %v", gotHeader, tt.wantHeader)
			}
			if !reflect.DeepEqual(gotAppinfo, tt.wantAppinfo) {
				t.Errorf("ParseMPLS() gotAppinfo = %v, want %v", gotAppinfo, tt.wantAppinfo)
			}
			if !reflect.DeepEqual(gotPlaylist, tt.wantPlaylist) {
				t.Errorf("ParseMPLS() gotPlaylist = %v, want %v", gotPlaylist, tt.wantPlaylist)
			}
			if !reflect.DeepEqual(gotChapterMarks, tt.wantChapterMarks) {
				t.Errorf("ParseMPLS() gotChapterMarks = %v, want %v", gotChapterMarks, tt.wantChapterMarks)
			}
			if !reflect.DeepEqual(gotExtensiondata, tt.wantExtensiondata) {
				t.Errorf("ParseMPLS() gotExtensiondata = %v, want %v", gotExtensiondata, tt.wantExtensiondata)
			}
		})
	}
}
